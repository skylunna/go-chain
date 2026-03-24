package main

import (
	"fmt"
	"sync"
)

// Blockchain 代表整个区块链
type Blockchain struct {
	Blocks []*Block
	mu     sync.Mutex // 新增：互斥锁，保护并发安全
}

// 定义挖矿难度，数字越大越慢
const Difficulty = 2

// 创建固定的创世区块
// 所有节点调用此函数都会生成完全相同的区块
func CreateGenesisBlock() *Block {
	block := &Block{
		Index:     0,
		Timestamp: "2026-03-04 00:00:00",
		Data:      "Genesis Block - Go-Chain v1.0",
		PrevHash:  "",
		Nonce:     0,
	}

	// 挖矿得到固定哈希
	// 难度固定，否则不同电脑挖出的Nonce不同，但哈希结果相同
	block.MineBlock(Difficulty)
	return block
}

// NewBlockchain 创建一个新的区块链，包含创世区块
func NewBlockchain() *Blockchain {
	// 创世区块也需要挖矿
	genesis := CreateGenesisBlock()
	// genesis.MineBlock(Difficulty)
	fmt.Printf("✅ 创世区块已加载: Hash=%s\n", genesis.Hash[:16]+"...")

	return &Blockchain{
		Blocks: []*Block{genesis},
	}
}

// 添加新区块到链上 - 线程安全
func (bc *Blockchain) AddBlock(data string) *Block {

	// 上锁，确保同一时间只有一个协程能修改区块链
	bc.mu.Lock()
	defer bc.mu.Unlock()

	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(prevBlock.Index+1, data, prevBlock.Hash)

	// 挖矿过程
	fmt.Printf("正在挖掘区块 %d ...\n", newBlock.Index)
	newBlock.MineBlock(Difficulty)

	bc.Blocks = append(bc.Blocks, newBlock)
	fmt.Printf("区块 %d 挖掘成功! Hash: %s\n", newBlock.Index, newBlock.Hash[:10]+"...")

	return newBlock
}

// IsChainValid 验证整个区块链是否有效
// 返回true表示链条完整未被篡改
func (bc *Blockchain) IsChainValid() bool {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	for i := 1; i < len(bc.Blocks); i++ {
		currentBlock := bc.Blocks[i]
		prevBlock := bc.Blocks[i-1]

		// 1. 检查当前区块的哈希值是否正确（数据是否被篡改）
		if currentBlock.Hash != currentBlock.CalculateHash() {
			return false
		}

		// 2. 检查当前区块的前哈希是否等于上一个区块的哈希（链条是否断裂）
		if currentBlock.PrevHash != prevBlock.Hash {
			return false
		}
	}

	return true
}

// 获取所有区块的副本（防止外部直接修改）
func (bc *Blockchain) GetBlocks() []*Block {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	// 返回副本
	blocksCopy := make([]*Block, len(bc.Blocks))
	copy(blocksCopy, bc.Blocks)
	return blocksCopy
}

// 获取单个区块
func (bc *Blockchain) GetBlock(index int) *Block {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	if index < 0 || index >= len(bc.Blocks) {
		return nil
	}

	return bc.Blocks[index]
}
