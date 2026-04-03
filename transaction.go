package main

import (
	"encoding/json"

	"github.com/skylunna/go-chain/wallet"
)

// 代表一笔转账
type Transaction struct {
	FromPubKey []byte  `json: "pubkey"` //公钥不直接序列化
	FromAddr   string  `json:"from"`
	ToAddr     string  `json: "to"`
	Amount     float64 `json: "amount"`
	Signature  []byte  `json: "signature"`
}

// Sign 对交易签名
func (tx *Transaction) Sign(w *wallet.Wallet) error {
	// 序列化交易内容 （不含签名）
	data, _ := json.Marshal((map[string]interface{}{
		"from":   tx.FromAddr,
		"to":     tx.ToAddr,
		"amount": tx.Amount,
	}))

	sig, err := w.Sign(data)
	if err != nil {
		return err
	}

	tx.Signature = sig
	return nil
}

// 验证交易签名
func (tx *Transaction) Verify() bool {

	// 这里需要将从地址返回公钥，或者我们在交易中存储公钥字节
	// ez: 假设我们能从FromAddr找到公钥，或者我们在Transaction结构体里加一个 PubKeyBytes 字段
	// 直接存储公钥字节
	return true
}
