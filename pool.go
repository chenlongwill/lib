package lib

import (
	"fmt"
)

// 消息缓存池，主要用来过滤峰谷大流量请求
// 消息将存储在内存之中，重启丢失
type Worker struct {
	Name       string
	ID         int
	JobChannel chan interface{}
	WorkerPool chan chan interface{}
	QuitChan   chan bool
}

func NewWorker(id int, workerPool chan chan interface{}) Worker {
	worker := Worker{
		ID:         id,
		JobChannel: make(chan interface{}),
		WorkerPool: workerPool,
		QuitChan:   make(chan bool)}
	return worker
}

func (w Worker) Start(function func(interface{})) {
	go func() {
		for {
			w.WorkerPool <- w.JobChannel
			select {
			case value := <-w.JobChannel:
				function(value)
			case <-w.QuitChan:
				return
			}
		}
	}()
}

func (w Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}

type pool struct {
	Work       map[string]*Worker
	JobQueue   chan interface{}
	WorkerPool chan chan interface{}
}

func (p *pool) Stop(args ...string) {
	go func() {
		if len(args) >= 1 {
			for _, v := range args {
				if v, ok := p.Work[v]; ok {
					v.Stop()
				}
			}
		} else {
			for _, v := range p.Work {
				v.Stop()
			}
		}
	}()
}

func NewPool(function func(interface{}), size int) *pool {
	if size <= 0 || size > 20 {
		panic("缓存队列最大开启20个，最小1个")
	}
	pl := &pool{Work: make(map[string]*Worker, size), JobQueue: make(chan interface{}), WorkerPool: make(chan chan interface{}, size)}
	for i := 0; i < size; i++ {
		worker := NewWorker(i+1, pl.WorkerPool)
		worker.Start(function)
		pl.Work[fmt.Sprintf("%d", i)] = &worker
	}
	go func() {
		for {
			select {
			case job := <-pl.JobQueue:
				go func(jobChannel interface{}) {
					worker := <-pl.WorkerPool
					worker <- jobChannel
				}(job)
			}
		}
	}()
	return pl
}
