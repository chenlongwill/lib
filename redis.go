package lib

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	"time"
)

var (
	// redis连接池
	redispool    *redis.Pool
	defaultTable string
)

func RedisConnect() {
	defaultTable = beego.AppConfig.String("dbname")
	if defaultTable == "" {
		panic("配置文件未配置appname变量，redis数据库名")
	}
	connectInit()
	c := redispool.Get()
	defer c.Close()
	err := c.Err()
	if err != nil {
		logs.GetLogger("init").Println("redis数据库连接失败[", err, "]")
		panic("redis数据库连接失败")
	} else {
		logs.GetLogger("init").Println("redis数据库连接成功")
	}
}

// 初始化链接
func connectInit() {
	dialFunc := func() (c redis.Conn, err error) {
		c, err = redis.Dial("tcp", beego.AppConfig.DefaultString("redis", ":6379"))
		if err != nil {
			return nil, err
		}

		if beego.AppConfig.String("redispwd") != "" {
			if _, err := c.Do("AUTH", beego.AppConfig.String("redispwd")); err != nil {
				c.Close()
				return nil, err
			}
		}

		_, selecterr := c.Do("SELECT", beego.AppConfig.DefaultInt("redisdefaultdb", 0))
		if selecterr != nil {
			c.Close()
			return nil, selecterr
		}
		return
	}
	redispool = &redis.Pool{
		MaxIdle:     beego.AppConfig.DefaultInt("redismaxidle", 10),
		IdleTimeout: 180 * time.Second,
		Dial:        dialFunc,
	}
}

type proredis string
type RedisType proredis

// 初始化一个集合
func NewRedis(table string) (rs proredis) {
	if table == "" {
		logs.Warning("redis Using未设置参数，将使用默认表")
		rs = proredis(defaultTable)
	} else {
		rs = proredis(table)
	}
	return
}

// redis基础命令
func (rs *proredis) do(commandName string, args ...interface{}) (reply interface{}, err error) {
	c := redispool.Get()
	defer c.Close()
	reply, err = c.Do(commandName, args...)
	if err != nil {
		logs.Error("redis error:[", err, "]")
	}
	return
}

// redis基础命令
func (rs *proredis) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	c := redispool.Get()
	defer c.Close()
	reply, err = c.Do(commandName, args...)
	if err != nil {
		logs.Error("redis error:[", err, "]")
	}
	return
}

// 根据key获取value
func (rs *proredis) Get(key string) interface{} {
	if v, err := rs.do("GET", fmt.Sprintf("%v_%s", *rs, key)); err == nil {
		return v
	}
	return nil
}

func (rs *proredis) GetInt(key string) (int, error) {
	return redis.Int(rs.do("GET", fmt.Sprintf("%v_%s", *rs, key)))
}

func (rs *proredis) GetString(key string) (string, error) {
	return redis.String(rs.do("GET", fmt.Sprintf("%v_%s", *rs, key)))
}

func (rs *proredis) GetDefaultString(key string) string {
	str, _ := redis.String(rs.do("GET", fmt.Sprintf("%v_%s", *rs, key)))
	return str
}

func (rs *proredis) GetBool(key string) (bool, error) {
	return redis.Bool(rs.do("GET", fmt.Sprintf("%v_%s", *rs, key)))
}

func (rs *proredis) GetIntMap(key string) (map[string]int, error) {
	return redis.IntMap(rs.do("GET", fmt.Sprintf("%v_%s", *rs, key)))
}

func (rs *proredis) GetInt64Map(key string) (map[string]int64, error) {
	return redis.Int64Map(rs.do("GET", fmt.Sprintf("%v_%s", *rs, key)))
}

func (rs *proredis) GetFloat64(key string) (float64, error) {
	return redis.Float64(rs.do("GET", fmt.Sprintf("%v_%s", *rs, key)))
}

func (rs *proredis) GetStringMap(key string) (map[string]string, error) {
	return redis.StringMap(rs.do("GET", fmt.Sprintf("%v_%s", *rs, key)))
}

func (rs *proredis) GetInt64(key string) (int64, error) {
	return redis.Int64(rs.do("GET", fmt.Sprintf("%v_%s", *rs, key)))
}

func (rs *proredis) GetInts(key string) ([]int, error) {
	return redis.Ints(rs.do("GET", fmt.Sprintf("%v_%s", *rs, key)))
}

func (rs *proredis) GetValues(key string) ([]interface{}, error) {
	return redis.Values(rs.do("GET", fmt.Sprintf("%v_%s", *rs, key)))
}

func (rs *proredis) GetStrings(key string) ([]string, error) {
	return redis.Strings(rs.do("GET", fmt.Sprintf("%v_%s", *rs, key)))
}

func (rs *proredis) GetBytes(key string) ([]byte, error) {
	return redis.Bytes(rs.do("GET", fmt.Sprintf("%v_%s", *rs, key)))
}

func (rs *proredis) GetNoBytes(key string) ([]byte, error) {
	return redis.Bytes(rs.do("GET", key))
}

func (rs *proredis) GetUint64(key string) (uint64, error) {
	return redis.Uint64(rs.do("GET", fmt.Sprintf("%v_%s", *rs, key)))
}

// 获取一个集合中所有值
func (rs *proredis) GetAll() []interface{} {
	if v, err := rs.do("HKEYS", *rs); err == nil {
		var str []string
		for _, val := range v.([]interface{}) {
			str = append(str, fmt.Sprintf("%s", val))
		}
		logs.Debug("redis GetAll() keys:", str)
		return rs.GetMulti(str)
	}
	return nil
}

// 获取多个key值
func (rs *proredis) GetMulti(keys []string) []interface{} {
	size := len(keys)
	var rv []interface{}
	c := redispool.Get()
	defer c.Close()
	var err error
	for _, key := range keys {
		err = c.Send("GET", fmt.Sprintf("%v_%s", *rs, key))
		if err != nil {
			goto ERROR
		}
	}
	if err = c.Flush(); err != nil {
		goto ERROR
	}
	for i := 0; i < size; i++ {
		if v, err := c.Receive(); err == nil {
			if v != nil {
				rv = append(rv, v.([]byte))
			}
		} else {
			rv = append(rv, err)
		}
	}
	return rv
ERROR:
	rv = rv[0:0]
	for i := 0; i < size; i++ {
		rv = append(rv, nil)
	}

	return rv
}

// 添加key-value，存储string、int、[]byte(把对象转成json字节组)
// var buf bytes.Buffer
// fmt.Fprint(&buf, arg)
// err = c.writeBytes(buf.Bytes())
func (rs *proredis) Put(key string, val interface{}) error {
	var err error
	if _, err = rs.do("SET", fmt.Sprintf("%v_%s", *rs, key), val); err != nil {
		return err
	}

	if _, err = rs.do("HSET", *rs, fmt.Sprintf("%v_%s", *rs, key), true); err != nil {
		return err
	}
	return err
}

// 有序队列对于对象存储，使用json字节组转换
// key集合名字,sort排序依据
func (rs *proredis) PutSortStructList(key string, sort int64, val interface{}) error {
	err := rs.PutStruct(key, val)
	if err != nil {
		return fmt.Errorf("Redis PutSortStructList PutJson %v", err)
	}
	err = rs.PutSortStringList(sort, fmt.Sprintf("%v_%s", *rs, key))
	if err != nil {
		return fmt.Errorf("PutSortStringList %v", err)
	}
	return nil
}

// 有序队列对于字符串key存储
// key集合名字,sort排序依据
func (rs *proredis) PutSortStringList(sort int64, val interface{}) error {
	if _, err := rs.do("ZADD", fmt.Sprintf("sort_%v", *rs), sort, val); err != nil {
		return err
	}
	return nil
}

// 有序队列对于字符串key存储
// key集合名字,sort排序依据
func (rs *proredis) DeleteSortStringList(val interface{}) error {
	logs.Debug(val, fmt.Sprintf("sort_%v", *rs))
	if _, err := rs.do("ZREM", fmt.Sprintf("sort_%v", *rs), val); err != nil {
		return err
	}
	return nil
}

// 有序队列对于字符串key存储
// sort排序依据,max大于标识
func (rs *proredis) GetSortStringList(sort int64, max bool) (str []string, err error) {
	var strs []string
	if max {
		strs, err = redis.Strings(rs.do("ZRANGEBYSCORE", fmt.Sprintf("sort_%v", *rs), "+inf", sort, "WITHSCORES"))
	} else {
		strs, err = redis.Strings(rs.do("ZRANGEBYSCORE", fmt.Sprintf("sort_%v", *rs), "-inf", sort, "WITHSCORES"))
	}
	for i := 0; i < len(strs); i += 2 {
		str = append(str, strs[i])
	}
	return
}

// 有序队列对于对象获取，使用json字节组转换
// key集合名字,sort排序依据
func (rs *proredis) GetSortStructList(sort int64, max bool) []interface{} {
	strs, _ := rs.GetSortStringList(sort, max)
	return rs.GetMulti(strs)
}

// for i := 0; i < len(jss); i += 2 {
// 	err = json.Unmarshal(jss[i], &tt)
// 	if err != nil {
// 		panic(err)
// 	}
// 	logs.Debug("%v", tt)
// }

// 有序队列对于对象获取，使用json字节组转换
// key集合名字,sort排序依据
func (rs *proredis) DeleteSortList(key string) error {
	err := rs.Delete(key)
	if err != nil {
		return err
	}
	if _, err = rs.do("ZREM", fmt.Sprintf("sort_%v", *rs), key); err != nil {
		return err
	}
	return nil
}

// 获取对象
func (rs *proredis) GetStructSort(key string, val interface{}) error {
	js, err := redis.Bytes(rs.do("GET", key))
	if err != nil {
		return err
	}
	err = json.Unmarshal(js, val)
	if err != nil {
		return err
	}
	return nil
}

// 获取对象
func (rs *proredis) GetStructSortS(key string, val interface{}) error {
	js, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("Redis GETSET PutJson %v", err)
	}
	jss, err := redis.Bytes(rs.do("GETSET", key, js))
	if err != nil {
		return err
	}
	err = json.Unmarshal(jss, val)
	if err != nil {
		return err
	}
	return nil
}

// 对于对象存储，使用json字节组转换
func (rs *proredis) PutStruct(key string, val interface{}) error {
	js, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("Redis PutJson %v", err)
	}
	if _, err = rs.do("SET", fmt.Sprintf("%v_%s", *rs, key), js); err != nil {
		return err
	}

	if _, err = rs.do("HSET", *rs, fmt.Sprintf("%v_%s", *rs, key), true); err != nil {
		return err
	}
	return err
}

// 获取对象
func (rs *proredis) GetStruct(key string, val interface{}) error {
	js, err := rs.GetNoBytes(fmt.Sprintf("%v_%s", *rs, key))
	if err != nil {
		return err
	}
	err = json.Unmarshal(js, val)
	if err != nil {
		return err
	}
	return nil
}

// 对于对象存储，使用json字节组转换
func (rs *proredis) PutStructEx(key string, val interface{}, timeout time.Duration) error {
	js, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("Redis PutJson %v", err)
	}
	if _, err = rs.do("SETEX", fmt.Sprintf("%v_%s", *rs, key), int64(timeout/time.Second), js); err != nil {
		return err
	}

	if _, err = rs.do("HSET", *rs, fmt.Sprintf("%v_%s", *rs, key), true); err != nil {
		return err
	}
	return err
}

// 添加带定时器的key
func (rs *proredis) PutEX(key string, val interface{}, timeout time.Duration) error {
	var err error
	if _, err = rs.do("SETEX", fmt.Sprintf("%v_%s", *rs, key), int64(timeout/time.Second), val); err != nil {
		return err
	}

	if _, err = rs.do("HSET", *rs, fmt.Sprintf("%v_%s", *rs, key), true); err != nil {
		return err
	}
	return err
}

// 删除一个key
func (rs *proredis) Delete(key string) error {
	logs.Warning("redis delete key [", key, "]")
	var err error
	if _, err = rs.do("DEL", fmt.Sprintf("%v_%s", *rs, key)); err != nil {
		return err
	}
	_, err = rs.do("HDEL", *rs, fmt.Sprintf("%v_%s", *rs, key))
	return err
}

// 判断key是否存在
func (rs *proredis) IsExist(key string) bool {
	v, err := redis.Bool(rs.do("EXISTS", fmt.Sprintf("%v_%s", *rs, key)))
	if err != nil {
		return false
	}
	if v == false {
		if _, err = rs.do("HDEL", *rs, fmt.Sprintf("%v_%s", *rs, key)); err != nil {
			return false
		}
	}
	return v
}

// 对key值累加.
func (rs *proredis) Incr(key string) error {
	_, err := redis.Bool(rs.do("INCRBY", fmt.Sprintf("%v_%s", *rs, key), 1))
	return err
}

// 对key值累减.
func (rs *proredis) Decr(key string) error {
	_, err := redis.Bool(rs.do("INCRBY", fmt.Sprintf("%v_%s", *rs, key), -1))
	return err
}

// 删除redis集合.
func (rs *proredis) ClearAll() error {
	cachedKeys, err := redis.Strings(rs.do("HKEYS", *rs))
	if err != nil {
		return err
	}
	for _, str := range cachedKeys {
		if _, err = rs.do("DEL", str); err != nil {
			return err
		}
	}
	_, err = rs.do("DEL", *rs)
	return err
}
