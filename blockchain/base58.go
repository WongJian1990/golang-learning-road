package main

import (
	"bytes"
	"math/big"
)

var b58Alphabet = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

//ReverseBytes 字节反序
func ReverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

//Base58Encode 加密
func Base58Encode(in []byte) []byte {
	var res []byte
	x := big.NewInt(0).SetBytes(in)
	base := big.NewInt(int64(len(b58Alphabet)))
	zero := big.NewInt(0)
	mod := &big.Int{}
	for x.Cmp(zero) != 0 {
		x.DivMod(x, base, mod)
		res = append(res, b58Alphabet[mod.Int64()])
	}
	ReverseBytes(res)
	for b := range in {
		if b == 0x00 {
			res = append([]byte{b58Alphabet[0]}, res...)
		} else {
			break
		}
	}
	return res
}

//Base58Decode 解密
func Base58Decode(in []byte) []byte {
	res := big.NewInt(0)
	zeroBytes := 0
	for b := range in {
		if b == 0x00 {
			zeroBytes++
		}
	}
	payload := in[zeroBytes:]
	for _, b := range payload {
		charIndex := bytes.IndexByte(b58Alphabet, b)
		res.Mul(res, big.NewInt(58))
		res.Add(res, big.NewInt(int64(charIndex)))
	}
	decoded := res.Bytes()
	decoded = append(bytes.Repeat([]byte{byte(0x00)}, zeroBytes), decoded...)
	return decoded
}
