package main

import (
	"fmt"
	"sync"
)

// Blockchain 代表整个区块链
type Blockchain struct {
	Blocks []*Block
	mu		sync.Mutex	// 新增：互斥锁，保护并发安全
}

// 定义挖矿难度，数字越大越慢
const Difficulty = 2

// NewBlockchain 创建一个新的区块链，包含创世区块
func NewBlockchain() *Blockchain {
	// 创世区块也需要挖矿
	genesis := NewBlock(0, "Genesis Block", "")
	genesis.MineBlock(Difficulty)

	return &Blockchain {
		Blocks: []*Block{genesis},
	}
}

// 添加新区块到链上 - 线程安全
func (bc *Blockchain) AddBlock(data string) {

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