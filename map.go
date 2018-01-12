package lib

import (
	"sync"
)

type Maps struct {
	Lock   sync.RWMutex
	Bucket map[string]interface{}
}

func (this *Maps) NewMaps(n int) {
	this.Lock.Lock()
	if n == 0 {
		this.Bucket = make(map[string]interface{})
	} else {
		this.Bucket = make(map[string]interface{}, n)
	}
	this.Lock.Unlock()
	return
}

func (this *Maps) Get(key string) (value interface{}, ok bool) {
	value, ok = this.Bucket[key]
	return
}

func (this *Maps) Put(key string, value interface{}) {
	this.Lock.Lock()
	this.Bucket[key] = value
	this.Lock.Unlock()
	return
}

func (this *Maps) PutNoExit(key string, value interface{}) (val interface{}, ok bool) {
	if val, ok = this.Bucket[key]; ok {
		return
	}
	this.Lock.Lock()
	this.Bucket[key] = value
	this.Lock.Unlock()
	return
}

func (this *Maps) Delete(key string) {
	this.Lock.Lock()
	delete(this.Bucket, key)
	this.Lock.Unlock()
	return
}
