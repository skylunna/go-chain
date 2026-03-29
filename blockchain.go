package main

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// Blockchain 代表整个区块链
type Blockchain struct {
	db 		*leveldb.DB	// 数据库引用
	mu     sync.Mutex // 新增：互斥锁，保护并发安全
}

// InitBlockchain 初始化区块链 (打开数据库或创建创世区块)
func InitBlockchain(dbPath string) (*Blockchain, error) {	// 多返回值
	// 1. 打开数据库
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return nil, err
	}

	bc := &Blockchain{
		db: db
	}

	// 2. 检查数据库是否为空
	hasData, _ := db.Has([]byte("chain_tip", nil))
	if !hasData {
		// 数据库为空，创建创世区块
		fmt.Println("📦 数据库为空，初始化创世区块...")
		genesis := NewBlock(0, "Genesis Block", "")
		genesis.MineBlock(Difficulty)
		err = bc.saveBlock(genesis)

		if err != nil {
			return nil, err
		}
	} else {
		fmt.Println("📦 检测到已有数据，加载区块链...")
	}

	return bc, nil
}

// 将区块保存到数据库
func (bc *Blockchain) saveBlock(block *Block) error {

	batch := new(leveldb.Batch)

	// 1. 序列化区块
	data, err := json.Marshal(block)
	if err != nil {
		return err
	}

	// 2. 保存区块数据 key: "block:1"
	blockKey := fmt.Sprintf("block:%d", block.Index)
	batch.Put([]byte(blockKey), data)

	// 3. 更新链顶高度 key: "chain_tip"
	batch.Put([]byte("chain_tip", []byte(fmt.Sprintf("%d", block.Index))))

	// 4. 写入数据库
	return bc.db.Write(batch, nil)
}

// 根据高度获取区块
func (bc *Blockchain) GetBlock(height int) (*Block, error) {
	key := fmt.Sprintf("block:%d", height)
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
