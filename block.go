package main

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

// Block代表区块链中的一个区块
type Block struct {
	Index		int 	// 区块的高度
	Timestamp	string	// 时间戳
	Data		string	// 区块数据
	PrevHash	string	// 前一个区块的哈希
	Hash		string	// 当前区块的哈希
}

// 计算区块的哈希值
// 我们将区块的所有关键信息拼接在一起，然后进行SHA256运算
func (b *Block) CalculateHash() string {
	// 简单起见，我们将所有字段拼接成字符串
	// 在实际生产中，通常会将结构体序列化为JSON后再哈希，以保证一致性
	record := string(b.Index) + b.Timestamp + b.Data + b.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)

	return hex.EncodeToString(hashed)
}

// NewBlock 创建一个新区块
func NewBlock(index int, data string, prevHash string) *Block {
	block := &Block {
		Index: index,
		Timestamp: time.Now().String(),
		Data: data,
		PrevHash: prevHash,
	}

	// 创建后立即计算哈希
	block.Hash = block.CalculateHash()

	return block
}