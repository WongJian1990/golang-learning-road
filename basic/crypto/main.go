package main

import (
	"crypto"
	"crypto/aes"
	"crypto/md5"
	"fmt"
	"io"
)

func main() {

	//md5
	res := crypto.MD5.Available()
	fmt.Println("res: ", res)

	sun := md5.Sum([]byte("Hello, World"))
	fmt.Printf("%x\n", sun)

	md := md5.New()
	fmt.Println(md.Size())

	io.WriteString(md, "加入redis里面有1亿个key,其中有10w个key是以固定的前缀开头的，如何全部找出 使用keys指令可以扫出指定模式的key列表")
	io.WriteString(md, "edis单节点存在单点故障，为了解决单点问题，一般需要对redis配置从节点， 然后使用哨兵来监听主节点的存活状态， 如果主节点挂掉，"+
		"从节点能继续提供缓存功能")

	fmt.Printf("%x\n", md.Sum(nil))

	//aes
	ciph, err := aes.NewCipher([]byte("3381681290492102"))
	if err != nil {
		fmt.Println("aes: ", err)
	}
	fmt.Println("Here .....")
	dst := make([]byte, 16)
	src := []byte("Hello Hello World")
	for len(src)%aes.BlockSize != 0 {
		src = append(src, ' ')
	}
	count := len(src) / aes.BlockSize
	for i := 0; i < count; i++ {
		data := src[i*aes.BlockSize : (i+1)*aes.BlockSize]
		ciph.Encrypt(dst, data)
		fmt.Printf("%x\n", dst)

		ciph.Decrypt(dst, dst)
		fmt.Printf("%s\n", string(dst))
	}

}
