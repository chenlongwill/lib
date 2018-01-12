package lib

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

const (
	// 客户端请求建立握手
	OP_HANDSHAKE       = int32(0)
	OP_HANDSHAKE_REPLY = int32(1)
	// 客户端发起心跳包
	OP_HEARTBEAT       = int32(2)
	OP_HEARTBEAT_REPLY = int32(3)
	// 断开连接
	OP_DISCONNECT       = int32(4)
	OP_DISCONNECT_REPLY = int32(5)
	// 客户端发起权限认证
	OP_AUTH       = int32(6)
	OP_AUTH_REPLY = int32(7)

	// 系统警告通知
	OP_SYS_WARNING_REPLY = int32(9)
	// 系统信息通知
	OP_SYS_INFO_REPLY = int32(11)
	// 系统资金变化通知
	OP_SYS_AMT_REPLY = int32(13)
	// 系统ODS行情推送
	OP_SYS_ODS_REPLY = int32(15)

	// 进入和退出房间
	OP_ROOM_IN        = int32(16)
	OP_ROOM_IN_REPLY  = int32(17)
	OP_ROOM_OUT       = int32(18)
	OP_ROOM_OUT_REPLY = int32(19)

	// 客户端给用户或者房间发消息
	OP_SEND_USER_MSG       = int32(20)
	OP_SEND_USER_MSG_REPLY = int32(21)
	OP_SEND_ROOM_MSG       = int32(22)
	OP_SEND_ROOM_MSG_REPLY = int32(23)

	// 用户A收到用户B的消息
	OP_RECV_USER_MSG_REPLY = int32(25)
	// 用户A收到房间001的消息
	OP_RECV_ROOM_MSG_REPLY = int32(27)

	// 系统订单结算通知
	OP_SYS_SETTLE_REPLY = int32(29)
	// 系统用户信息变化通知
	OP_SYS_USER_INFO_REPLY = int32(31)

	// 消息系统内部专用
	OP_RAW          = int32(255)
	OP_PROTO_READY  = int32(254)
	OP_PROTO_FINISH = int32(253)
)

// 消息推送更新用户账户余额
type SystemMessage struct {
	Code int32
	Msg  string
	Data interface{}
}

func PushUserMessage(op int32, uid int, v interface{}) bool {
	url := ps("http://%s/push/user?op=%d&uid=%d", beego.AppConfig.DefaultString("push_addr", "localhost:57172"), op, uid)
	js, err := json.Marshal(v)
	if err != nil {
		logs.Error("err[%v]json解析失败", err)
		return false
	}
	resp, err := HttpPost(url, "application/json;charset=utf-8", bytes.NewBuffer(js))
	if err != nil {
		logs.Error("err[%v]HttpPost失败", err)
		return false
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var rep SystemMessage
	err = json.Unmarshal(body, &rep)
	if err != nil {
		logs.Error("err[%v]json解析失败", err)
		return false
	}
	if rep.Code != 0 {
		logs.Error("err[%v]HttpPost失败", rep)
		return false
	}
	logs.Info("给用户[%d]推送消息成功,message[%v]", uid, v)
	return true
}

func PushRoomMessage(op int32, rid int, v interface{}) bool {
	url := ps("http://%s/push/room?op=%d&rid=%d", beego.AppConfig.DefaultString("push_addr", "localhost:57172"), op, rid)
	js, err := json.Marshal(v)
	if err != nil {
		logs.Error("err[%v]json解析失败", err)
		return false
	}
	resp, err := HttpPost(url, "application/json;charset=utf-8", bytes.NewBuffer(js))
	if err != nil {
		// logs.Error("err[%v]HttpPost失败", err)
		return false
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var rep SystemMessage
	err = json.Unmarshal(body, &rep)
	if err != nil {
		logs.Error("err[%v]json解析失败", err)
		return false
	}
	if rep.Code != 0 {
		logs.Error("err[%v]HttpPost失败", rep)
		return false
	}
	if op == OP_SYS_ODS_REPLY {
		// logs.Debug("给房间[%d]推送消息成功,message[%v]", rid, v)
	}
	return true
}
