package main

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"
	"fmt"
)

// Block代表区块链中的一个区块
type Block struct {
	Index		int 	// 区块的高度
	Timestamp	string	// 时间戳
	Data		string	// 区块数据
	PrevHash	string	// 前一个区块的哈希
	Hash		string	// 当前区块的哈希

	Nonce		int		// 用于挖矿的随机数
}

// 计算区块的哈希值
// 我们将区块的所有关键信息拼接在一起，然后进行SHA256运算
func (b *Block) CalculateHash() string {
	// 简单起见，我们将所有字段拼接成字符串
	// 在实际生产中，通常会将结构体序列化为JSON后再哈希，以保证一致性
	// Integer to ASCII
	record := string(b.Index) + b.Timestamp + b.Data + b.PrevHash + strconv.Itoa(b.Nonce)
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)

	return hex.EncodeToString(hashed)
}

// MineBlock 执行工作量证明(挖矿)
// difficulty 代表哈希值前缀需要有多少个0
func (b *Block) MineBlock(difficulty int) {
	// 目标前缀，例如 difficulty=2，目标就是 "00"
	target := ""
	for i := 0; i < difficulty; i++ {
		target += "0"
	}

	// 不断尝试不同的 Nonce, 直到哈希值满足难度要求
	for {
		hash := b.CalculateHash()
		// 检查哈希值是否以目标前缀开头
		if hash[:difficulty] == target {
			fmt.Printf("区块挖矿成功！Hash: %s, Nonce :%d\n", hash, b.Nonce)
			b.Hash = hash
			break
		}
		b.Nonce++
	}
 }

// NewBlock 创建一个新区块
// 创建后不再直接计算哈希，而是需要挖矿
func NewBlock(index int, data string, prevHash string) *Block {
	block := &Block {
		Index: index,
		Timestamp: time.Now().String(),
		Data: data,
		PrevHash: prevHash,
		Nonce: 0,
	}

	// 创建后立即计算哈希
	// block.Hash = block.CalculateHash()

	return block
}