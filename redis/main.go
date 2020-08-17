package main

import (
	"fmt"
	"hash/crc32"
	"time"

	"github.com/gomodule/redigo/redis"
)

//redis分片连接池
var (
	pools []*redis.Pool
)

//构建redis连接池
func newPool(addr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 360 * time.Second,
		// Dial or DialContext must be set. When both are set, DialContext takes precedence over Dial.
		Dial: func() (redis.Conn, error) { return redis.Dial("tcp", addr) },
	}
}

//构建redis分片
func newPools(addrs []string) []*redis.Pool {

	temp := make([]*redis.Pool, 0)

	for _, addr := range addrs {
		pool := newPool(addr)
		if pool != nil {
			temp = append(temp, pool)
		}
	}
	return temp
}

//Get 通过key返回redis分片连接
func Get(pools []*redis.Pool, key string) (redis.Conn, error) {
	if len(pools) == 0 {
		return nil, fmt.Errorf("Invalid Pools")
	}
	cs := crc32.ChecksumIEEE([]byte(key))
	index := cs % uint32(len(pools))
	conn := pools[index].Get()
	return conn, nil
}

func main() {

	pools := newPools([]string{"localhost:6379", "localhost:6380", "localhost:6381"})
	conn, err := Get(pools, "gemini/vdts")
	if err != nil {
		fmt.Println("Get redis connection failed: ", err)
		return
	}
	//设置string
	_, err = conn.Do("set", "gemini/vdts", "hello world")
	if err != nil {
		fmt.Println("ERR: ", err)
	}
	conn.Close()

	conn, err = Get(pools, "gemini/zdts")
	defer conn.Close()
	if err != nil {
		fmt.Println("Get redis connection failed: ", err)
		return
	}
	_, err = conn.Do("set", "gemini/zdts", "hello world")
	if err != nil {
		fmt.Println("ERR[SET]: ", err)
	}
	//获取string
	reply, err := redis.String(conn.Do("get", "gemini/zdts"))
	if err != nil {
		fmt.Println("ERR[GET]: ", err)
	}
	fmt.Println("reply: ", reply)

	//设置list
	_, err = conn.Do("lpush", "gemini/zdts/list", "hello", "world", "from", "wang", "jian")

	//获取list
	replies, err := redis.Values(conn.Do("lrange", "gemini/zdts/list", 0, 20))
	if err != nil {
		fmt.Println("ERR[LRANGE]: ", err)
	}
	for _, v := range replies {
		// fmt.Println("kind: ", reflect.TypeOf(v).Kind().String())
		fmt.Println("value: ", string(v.([]byte)))
	}

	//删除list
	_, err = conn.Do("del", "gemini/zdts/list")
	if err != nil {

		fmt.Println("ERR[DEL]: ", err)
	}

	_, err = conn.Do("lpush", "gemini/zdts/list", 10, 20, 30, 40, 50)
	replies, err = redis.Values(conn.Do("lrange", "gemini/zdts/list", 0, 20))
	if err != nil {
		fmt.Println("ERR[LRANGE]: ", err)
	}
	for _, v := range replies {
		// fmt.Println("kind: ", reflect.TypeOf(v).Kind().String())
		fmt.Println("value: ", string(v.([]byte)))
	}

	//删除list
	_, err = conn.Do("del", "gemini/zdts/list")
	if err != nil {

		fmt.Println("ERR[DEL]: ", err)
	}

	//设置 hash字典
	_, err = conn.Do("hmset", "gemini/zdts/hash", "value1", 10, "value2", 20, "value3", 30, "value3", 40)
	if err != nil {
		fmt.Println("ERR[HMSET]: ", err)
	}
	//获取hash字典
	replies, err = redis.Values(conn.Do("HGETALL", "gemini/zdts/hash"))
	if err != nil {
		fmt.Println("ERR[HGETALL]: ", err)
	}
	for _, v := range replies {
		fmt.Printf("v:%v\n", string(v.([]byte)))
	}

	//删除
	_, err = conn.Do("del", "gemini/zdts/hash")
	if err != nil {
		fmt.Println("ERR[DEL]: ", err)
	}

	//设置set集合
	_, err = conn.Do("sadd", "gemini/zdts/set", "value1", "value2", "value2", "value3", "value3", "value4", "value4")
	if err != nil {
		fmt.Println("ERR[SADD]: ", err)
	}
	//获取所有成员
	replies, err = redis.Values(conn.Do("smembers", "gemini/zdts/set"))
	if err != nil {
		fmt.Println("ERR[SMEMBERS]: ", err)
	}
	for _, v := range replies {
		fmt.Println("value: ", string(v.([]byte)))
	}

	//删除set key
	_, err = conn.Do("del", "gemini/zdts/set")
	if err != nil {
		fmt.Println("ERR[DEL]: ", err)
	}

	t1 := time.Now().UnixNano()
	time.Sleep(1 * time.Millisecond)
	t2 := time.Now().UnixNano()
	time.Sleep(1 * time.Millisecond)
	t3 := time.Now().UnixNano()

	//设置有序集合
	_, err = conn.Do("zadd", "gemini/zdts/zset", t3, "value1", t2, "value2",
		t1, "value3")
	if err != nil {
		fmt.Println("ERR[ZADD]: ", err)
	}
	//获取range范围内的成员
	replies, err = redis.Values(conn.Do("zrange", "gemini/zdts/zset", 0, 20, "WITHSCORES"))
	if err != nil {
		fmt.Println("ERR[ZRANGE]: ", err)
	}
	for _, v := range replies {
		fmt.Println("value: ", string(v.([]byte)))
	}

	//删除 zset key
	_, err = conn.Do("del", "gemini/zdts/zset")
	if err != nil {
		fmt.Println("ERR[DEL]: ", err)
	}

}
