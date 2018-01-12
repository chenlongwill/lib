package lib

import (
	"net"
	"net/http"
	"time"
)

// 带超时设置的http请求
var c *http.Client = &http.Client{
	Transport: &http.Transport{
		Dial: func(netw, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(netw, addr, time.Second*3)
			if err != nil {
				return nil, err
			}
			return c, nil

		},
		MaxIdleConnsPerHost:   10,
		ResponseHeaderTimeout: time.Second * 2,
	},
}

// var HttpPost = c.Post
// var HttpGet = c.Get

var HttpPost = http.Post
var HttpGet = http.Get

// 案例 get
// func httpGet() {
// 	resp, err := http.Get("http://www.01happy.com/demo/accept.php?id=1")
// 	if err != nil {
// 		// handle error
// 	}
// 	defer resp.Body.Close()
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		// handle error
// 	}
// 	fmt.Println(string(body))
// }

// 案例 post
// 	url := ps("http://%s/push/room?op=%d&rid=%d", beego.AppConfig.DefaultString("push_addr", "localhost:57172"), op, rid)
// 	js, err := json.Marshal(v)
// 	if err != nil {
// 		logs.Error("err[%v]json解析失败", err)
// 		return false
// 	}
// 	resp, err := HttpPost(url, "application/json;charset=utf-8", bytes.NewBuffer(js))
// 	if err != nil {
// 		logs.Error("err[%v]HttpPost失败", err)
// 		return false
// 	}
//  defer resp.Body.Close()
// 	body, _ := ioutil.ReadAll(resp.Body)
// 	var rep SystemMessage
// 	err = json.Unmarshal(body, &rep)
// 	if err != nil {
// 		logs.Error("err[%v]json解析失败", err)
// 		return false
// 	}
// 	if rep.Code != 0 {
// 		logs.Error("err[%v]HttpPost失败", rep)
// 		return false
