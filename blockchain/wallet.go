package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/crypto/ripemd160"
)

const (
	version            = byte(0x00)
	walletFile         = "wallet.dat"
	addressChecksumLen = 4
)

//Wallet 交易钱包
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

//Wallets 钱包集合
type Wallets struct {
	Wallets map[string]*Wallet
}

//NewWallet 构造钱包
func NewWallet() *Wallet {
	private, public := newKeyPair()
	wallet := Wallet{private, public}
	return &wallet
}

func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, pubKey
}

//GenAddress 产生钱包地址
func (w Wallet) GenAddress() []byte {
	pubKeyHash := HashPubKey(w.PublicKey)
	versionPayload := append([]byte{version}, pubKeyHash...)
	checksum := checksum(versionPayload)
	fullPayload := append(versionPayload, checksum...)
	address := Base58Encode(fullPayload)
	return address
}

//HashPubKey 公钥hash
func HashPubKey(pubKey []byte) []byte {
	publicSha256 := sha256.Sum256(pubKey)
	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSha256[:])
	if err != nil {
		log.Panic(err)
	}
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)
	return publicRIPEMD160
}

//ValidateAddress 检查地址是否有效
func ValidateAddress(address string) bool {
	pubKeyHash := Base58Decode([]byte(address))
	actualCheckSum := pubKeyHash[len(pubKeyHash)-addressChecksumLen:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-addressChecksumLen]
	targetChecksum := checksum(append([]byte{version}, pubKeyHash...))
	return bytes.Compare(actualCheckSum, targetChecksum) == 0
}

func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])
	return secondSHA[:addressChecksumLen]
}

//NewWallets 构造钱包集合
func NewWallets() (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)
	err := wallets.LoadFromFile()
	return &wallets, err
}

//CreateWallet 向钱包集合中添加一个钱包
func (ws *Wallets) CreateWallet() string {
	wallet := NewWallet()
	address := fmt.Sprintf("%s", wallet.GenAddress())
	ws.Wallets[address] = wallet
	return address
}

//GetAddresses 返回钱包地址集
func (ws *Wallets) GetAddresses() []string {
	var addresses []string
	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}
	return addresses
}

//GetWallet 通过地址获取钱包
func (ws *Wallets) GetWallet(address string) Wallet {
	return *ws.Wallets[address]
}

//LoadFromFile 从文件中加载钱包集合
func (ws *Wallets) LoadFromFile() error {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}
	file, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}
	var wallets Wallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(file))
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panic(err)
	}
	ws.Wallets = wallets.Wallets
	return nil
}

//SaveToFile 保存到文件
func (ws Wallets) SaveToFile() {
	var content bytes.Buffer
	//接口可能产生的类型需要注册
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}
	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}
