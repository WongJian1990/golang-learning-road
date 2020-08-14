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
	_, err = conn.Do("set", "gemini/vdts", "hello world")
	if err != nil {
		fmt.Println("ERR: ", err)
	}
	conn.Close()

	conn, err = Get(pools, "gemini/zdts")
	if err != nil {
		fmt.Println("Get redis connection failed: ", err)
		return
	}
	_, err = conn.Do("set", "gemini/zdts", "hello world")
	if err != nil {
		fmt.Println("ERR: ", err)
	}

	conn.Close()

}
