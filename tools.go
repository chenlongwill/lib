package lib

import (
	"crypto/md5"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/mahonia"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"math/rand"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"time"
)

var ps = fmt.Sprintf

// 定时器
func StartTimer(f func()) {
	go func() {
		for {
			f()
			now := time.Now()
			// 计算下一个零点
			next := now.Add(time.Hour * 24)
			next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
			t := time.NewTimer(next.Sub(now))
			<-t.C
		}
	}()
}

// 函数 generator2，返回 通道(Channel)
func generator2() chan int {
	// 创建通道
	out := make(chan int)
	// 创建协程
	go func() {
		for {
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			//向通道内写入数据，如果无人读取会等待
			out <- r.Intn(999999)
		}
	}()
	return out
}

// 函数 Generator ，返回通道(Channel)，多路复用技术(双倍生成随机数)
func Generator() chan int {
	// 创建两个随机数生成器服务
	rand_generator_1 := generator2()
	rand_generator_2 := generator2()

	//创建通道
	out := make(chan int)

	//创建协程
	go func() {
		for {
			//读取生成器1中的数据，整合
			out <- <-rand_generator_1
		}
	}()
	go func() {
		for {
			//读取生成器2中的数据，整合
			out <- <-rand_generator_2
		}
	}()
	return out
}

// 对字符串进行md5加密
func StrToMD5(str string) string {
	md5Ctx1 := md5.New()
	md5Ctx1.Write([]byte(str))
	return fmt.Sprintf("%x", md5Ctx1.Sum(nil))
}

// 获取一个唯一键值
func GetSid() string {
	return fmt.Sprintf("%x", string(bson.NewObjectId()))
}

// 字符串防sql注入过滤,可以为空
func CheckArgString(a ...string) bool {
	for _, arg := range a {
		if arg == "" {
			return true
		} else {
			reg := `('|and|exec|insert|select|delete|update|count|%|chr|mid|master|truncate|char|declare|;|or|<)`
			rgx := regexp.MustCompile(reg)
			if rgx.MatchString(arg) {
				return false
			}
		}
	}
	return true
}

// 支付密码校验
func CheckArgZfpwd(zfpwd string) bool {
	if zfpwd == "" {
		return false
	} else {
		reg := `^\d{6}$`
		rgx := regexp.MustCompile(reg)
		if !rgx.MatchString(zfpwd) {
			return false
		}
	}
	return true
}

// 字符串手机号校验
func CheckArgPhone(phone string) bool {
	if phone == "" {
		return false
	} else {
		reg := `^1(3|4|5|7|8|9)\d{9}$`
		rgx := regexp.MustCompile(reg)
		if !rgx.MatchString(phone) {
			return false
		}
	}
	return true
}

// 身份证验证
func CheckArgIdCard(idCard string) bool {
	if idCard == "" {
		return false
	} else {
		reg := `^\d{17}(\d|X)$`
		rgx := regexp.MustCompile(reg)
		if !rgx.MatchString(idCard) {
			return false
		}
	}
	return true
}

// 参数非空校验
func CheckArgNotNull(a ...interface{}) bool {
	for _, arg := range a {
		switch reflect.TypeOf(arg).Kind() {
		case reflect.String:
			if arg.(string) == "" {
				return false
			} else {
				reg := `('|and|exec|insert|select|delete|update|count|%|chr|mid|master|truncate|char|declare|;|or|<)`
				rgx := regexp.MustCompile(reg)
				if rgx.MatchString(arg.(string)) {
					return false
				}
			}
		case reflect.Int64:
			if arg.(int64) == 0 {
				return false
			}
		case reflect.Int32:
			if arg.(int32) == 0 {
				return false
			}
		case reflect.Int:
			if arg.(int) == 0 {
				return false
			}
		case reflect.Float32:
			if arg.(float32) == 0 {
				return false
			}
		case reflect.Float64:
			if arg.(float64) == 0 {
				return false
			}
		default:
			return false
		}
	}
	return true
}

// 结构体参数非空校验
func CheckStructArgNotNull(stru interface{}, fields ...string) string {
	// t := reflect.TypeOf(stru).Elem()
	doc := make(map[string]string)
	t := reflect.ValueOf(stru)
	for _, fname := range fields {
		t.FieldByName(fname).IsValid()
		m[fname] = val.FieldByName(fname).Interface()
	}
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Tag.Get("description") == "" {
			doc[t.Field(i).Name] = t.Field(i).Name
		} else {
			doc[t.Field(i).Name] = t.Field(i).Tag.Get("description")
		}
	}
	for _, arg := range a {
		if v, k := doc[reflect.TypeOf(arg).Name()]; k {
			switch reflect.TypeOf(arg).Kind() {
			case reflect.String:
				if arg.(string) == "" {
					return v
				} else {
					reg := `('|and|exec|insert|select|delete|update|count|%|chr|mid|master|truncate|char|declare|;|or|<)`
					rgx := regexp.MustCompile(reg)
					if rgx.MatchString(arg.(string)) {
						return "合理值:" + v
					}
				}
			case reflect.Int64:
				if arg.(int64) == 0 {
					return v
				}
			case reflect.Int32:
				if arg.(int32) == 0 {
					return v
				}
			case reflect.Int:
				if arg.(int) == 0 {
					return v
				}
			case reflect.Float32:
				if arg.(float32) == 0 {
					return v
				}
			case reflect.Float64:
				if arg.(float64) == 0 {
					return v
				}
			default:
				return v
			}
		}
	}
	return ""
}

// 结构体参数非空校验
func CheckStructStringSql(stru interface{}, a ...string) string {
	t := reflect.TypeOf(stru).Elem()
	doc := make(map[string]string)
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Tag.Get("description") == "" {
			doc[t.Field(i).Name] = t.Field(i).Name
		} else {
			doc[t.Field(i).Name] = t.Field(i).Tag.Get("description")
		}
	}
	for _, arg := range a {
		if v, k := doc[arg]; k {
			if arg == "" {
				return ""
			} else {
				reg := `('|and|exec|insert|select|delete|update|count|%|chr|mid|master|truncate|char|declare|;|or|<)`
				rgx := regexp.MustCompile(reg)
				if rgx.MatchString(arg) {
					return "合理值:" + v
				}
			}
		}
	}
	return ""
}

// 发送验证码
func SendVcode(telnum string) bool {
	// 将验证码存储在redis里
	redis := NewRedis("vcode")
	// 生成6位数字验证码
	rad := Generator()
	vcode := fmt.Sprintf("%06v", <-rad)
	err := redis.PutEX(telnum, vcode, 300*time.Second)
	if err != nil {
		logs.Error("验证码存储redis失败[%v]", err)
		return false
	}
	return SendMsgToPhone(telnum, fmt.Sprintf("尊敬的客户，您的手机验证码为：%s，本验证码5分钟之内有效。请保证是本人使用，否则请忽略此短信【XX商城】", vcode))
}
func SendMsgToPhone(telnum string, label string) bool {
	enc := mahonia.NewEncoder("GBK")
	content := enc.ConvertString(label)
	tmp := fmt.Sprintf("http://baidu.com/get/url?msg=%s", content)
	vl, _ := url.Parse(tmp)
	msg := vl.Query().Encode()
	req := fmt.Sprintf("cmd=send&uid=%s&psw=%s&mobiles=%s&msgid=%0404d%s&%s", "username", "pwd", telnum, time.Now().Unix(), telnum, msg)
	address := "http://kltx.sms10000.com.cn/sdk/SMS"
	resp, err := HttpPost(address, "application/x-www-form-urlencoded;charset=GB2312", strings.NewReader(req))
	if err != nil {
		logs.Error("发送信息到手机失败error[%v]resp[%v]", err, resp)
		return false
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error("发送信息到手机失败error[%v]", err)
		return false
	}
	if string(body) != "100" {
		logs.Error("发送信息到手机失败error[%s]返回码不等于100", body)
		return false
	}
	logs.Info("给用户[%s]推送短信消息成功,message[%s]", telnum, label)
	return true
}

//手机验证码验证
func CheckVcode(telnum, vcode string) bool {
	redis := NewRedis("vcode")

	vcodephone, err := redis.GetString(telnum)
	if err != nil {
		logs.Error("手机验证码查询失败:", err)
		return false
	}
	if vcodephone != vcode {
		logs.Error("手机验证码:[%s != %s]", vcodephone, vcode)

		return false

	}
	err = redis.Delete(telnum)
	if err != nil {
		logs.Error("[%s]手机验证码redis删除失败:[%v]", telnum, err)
	}
	return true

}

//正则验证必须是数字
func OnlyNumber(str string, maxLimit string) (bool, string) {
	if str == "" {
		return false, "请输入有效参数"
	}
	var flag bool = false
	var msg string = "验证通过"

	r, _ := regexp.Compile("^[0-9]{" + maxLimit + "}$")
	flag = r.MatchString(str)
	if !flag {
		msg = "参数格式错误，必须是" + maxLimit + "位数字"
	}
	return flag, msg
}
