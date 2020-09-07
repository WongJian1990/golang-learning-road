package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"
)

var subsidy = 50

//TxInput 交易输入
type TxInput struct {
	Txid      []byte
	Vout      int
	ScriptSig []byte
	PubKey    []byte
}

//TxOutput 交易输出
type TxOutput struct {
	Value        int
	ScriptPubKey []byte
}

//TxOutputs 输出集合
type TxOutputs struct {
	Outputs []TxOutput
}

//Transaction 交易单元
type Transaction struct {
	ID   []byte
	Vin  []TxInput
	Vout []TxOutput
}

//IsCoinBase 判断是否是coinbase交易
func (tx Transaction) IsCoinBase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && (tx.Vin[0].Vout == -1)
}

//Serialize 序列化
func (outs TxOutputs) Serialize() []byte {
	var encoded bytes.Buffer
	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(outs)
	if err != nil {
		log.Panic(err)
	}
	return encoded.Bytes()

}

//DeserializeOutputs 反序列化TxOutputs
func DeserializeOutputs(data []byte) TxOutputs {
	var outs TxOutputs
	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&outs)
	if err != nil {
		log.Panic(err)
	}
	return outs
}

//Serialize Transaction Serialize
func (tx Transaction) Serialize() []byte {
	var encoded bytes.Buffer
	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	return encoded.Bytes()
}

//Hash 返回交易的哈希
func (tx *Transaction) Hash() []byte {
	var hash [32]byte
	txCopy := *tx
	txCopy.ID = []byte{}
	hash = sha256.Sum256(txCopy.Serialize())
	return hash[:]
}

//Sign 签名
func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTxs map[string]Transaction) {
	if tx.IsCoinBase() {
		return
	}
	for _, vin := range tx.Vin {
		if prevTxs[hex.EncodeToString(vin.Txid)].ID == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}
	txCopy := tx.TrimmedCopy()
	for inID, vin := range txCopy.Vin {
		prevTx := prevTxs[hex.EncodeToString(vin.Txid)]
		txCopy.Vin[inID].ScriptSig = nil
		txCopy.Vin[inID].PubKey = prevTx.Vout[vin.Vout].ScriptPubKey
		txCopy.ID = txCopy.Hash()
		txCopy.Vin[inID].PubKey = nil
		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.ID)
		if err != nil {
			log.Panic(err)
		}
		sign := append(r.Bytes(), s.Bytes()...)
		tx.Vin[inID].ScriptSig = sign
	}
}

//TrimmedCopy 裁剪可以用作签名数据的交易数据
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	for _, vin := range tx.Vin {
		inputs = append(inputs, TxInput{vin.Txid, vin.Vout, nil, nil})
	}

	for _, vout := range tx.Vout {
		outputs = append(outputs, TxOutput{vout.Value, vout.ScriptPubKey})
	}
	txCopy := Transaction{tx.ID, inputs, outputs}
	return txCopy
}

//Verify 验证
func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if tx.IsCoinBase() {
		return true
	}
	for _, vin := range tx.Vin {
		if prevTXs[hex.EncodeToString(vin.Txid)].ID == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}
	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()
	for inID, vin := range tx.Vin {
		prevTx := prevTXs[hex.EncodeToString(vin.Txid)]
		txCopy.Vin[inID].ScriptSig = nil
		txCopy.Vin[inID].PubKey = prevTx.Vout[vin.Vout].ScriptPubKey
		txCopy.ID = txCopy.Hash()
		txCopy.Vin[inID].PubKey = nil
		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.ScriptSig)
		r.SetBytes(vin.ScriptSig[:(sigLen / 2)])
		s.SetBytes(vin.ScriptSig[(sigLen / 2):])
		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PubKey)
		x.SetBytes(vin.PubKey[:(keyLen / 2)])
		y.SetBytes(vin.PubKey[(keyLen / 2):])
		rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
		if ecdsa.Verify(&rawPubKey, txCopy.ID, &r, &s) == false {
			return false
		}
	}
	return true
}

//NewCoinBaseTx 创建一个coinbase交易
func NewCoinBaseTx(to, data string) *Transaction {
	if data == "" {
		randData := make([]byte, 20)
		_, err := rand.Read(randData)
		if err != nil {
			log.Panic(err)
		}
		data = fmt.Sprintf("%x", randData)
	}

	txin := TxInput{Txid: []byte{}, Vout: -1, ScriptSig: []byte(data)}
	txout := NewTxOutput(subsidy, to)
	tx := Transaction{nil, []TxInput{txin}, []TxOutput{*txout}}
	tx.ID = tx.Hash()
	return &tx
}

func (tx Transaction) String() string {
	var lines []string
	lines = append(lines, fmt.Sprintf("---Transaction %x", tx.ID))
	for i, input := range tx.Vin {
		lines = append(lines, fmt.Sprintf("	Input %d:", i))
		lines = append(lines, fmt.Sprintf("	TXID:	%x", input.Txid))
		lines = append(lines, fmt.Sprintf("	Out: 	%d", input.Vout))
		lines = append(lines, fmt.Sprintf("	Sign:	%x", input.ScriptSig))
		lines = append(lines, fmt.Sprintf("	PubKey:	%x", input.PubKey))
	}

	for i, output := range tx.Vout {
		lines = append(lines, fmt.Sprintf("	Output %d:", i))
		lines = append(lines, fmt.Sprintf("	Value:	%d", output.Value))
		lines = append(lines, fmt.Sprintf("	PubKey:	%x", output.ScriptPubKey))
	}
	return strings.Join(lines, "\n")
}

//NewUTXOTransaction 创建一笔新的交易
func NewUTXOTransaction(from, to string, amount int, UTXOSet *UTXOSet) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput
	//找到足够的未花费输出
	// acc, validOuptputs := chain.FindSpendableOutputs(from, amount)
	wallets, err := NewWallets()
	if err != nil {
		log.Panic(err)
	}
	wallet := wallets.GetWallet(from)
	pubKeyHash := HashPubKey(wallet.PublicKey)
	acc, validOuptputs := UTXOSet.FindSpendableOutputs(pubKeyHash, amount)
	if acc < amount {
		log.Panic("Error:Not enought founds")
	}
	for txid, outs := range validOuptputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}
		for _, out := range outs {
			input := TxInput{txID, out, nil, wallet.PublicKey}
			inputs = append(inputs, input)
		}
	}
	outputs = append(outputs, *NewTxOutput(amount, to))
	//如果UTXO总数超过所需，则找零
	if acc > amount {
		outputs = append(outputs, *NewTxOutput(acc-amount, from))
	}
	tx := Transaction{nil, inputs, outputs}
	tx.ID = tx.Hash()
	UTXOSet.BlockChain.SignTransaction(&tx, wallet.PrivateKey)
	return &tx
}

//Lock 锁定公钥
func (out *TxOutput) Lock(address []byte) {
	pubKeyHash := Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-addressChecksumLen]
	out.ScriptPubKey = pubKeyHash
}

//IsLockedWithKey 检查是否被公钥锁定
func (out *TxOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.ScriptPubKey, pubKeyHash) == 0
}

//NewTxOutput 构造交易输出
func NewTxOutput(value int, address string) *TxOutput {
	txo := &TxOutput{value, nil}
	txo.Lock([]byte(address))
	return txo
}

//UsesKey 检查交易地址是否可用
func (in *TxInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := HashPubKey(in.PubKey)
	return bytes.Compare(pubKeyHash, lockingHash) == 0
}
