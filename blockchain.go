package main

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
)

// 挖矿难度
const Difficulty = 2

// Blockchain 现在包含数据库引用
type Blockchain struct {
	db *leveldb.DB
	mu sync.Mutex
}

// InitBlockchain 初始化区块链（打开数据库或创建创世块）
func InitBlockchain(dbPath string) (*Blockchain, error) {
	// 1. 打开数据库
	db, err := leveldb.OpenFile(dbPath, nil) // nil 传默认配置
	if err != nil {
		return nil, err
	}

	// 创建 Blockchain 结构体实例，返回其指针
	// 创建一条区块链实例，并且把已打开的数据库连接交给它管理 bc
	bc := &Blockchain{db: db}

	// 2. 检查数据库是否为空
	hasData, _ := db.Has([]byte("chain_tip"), nil)
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

// saveBlock 将区块保存到数据库
func (bc *Blockchain) saveBlock(block *Block) error {
	// 创建一个空的 LevelDB 批量操作对象（batch)
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
	batch.Put([]byte("chain_tip"), []byte(fmt.Sprintf("%d", block.Index)))

	// 4. 写入数据库
	return bc.db.Write(batch, nil)
}

// GetBlock 根据高度获取区块
func (bc *Blockchain) GetBlock(height int) (*Block, error) {
	key := fmt.Sprintf("block:%d", height)
	data, err := bc.db.Get([]byte(key), nil)
	if err != nil {
		return nil, err
	}

	var block Block
	err = json.Unmarshal(data, &block)
	return &block, err
}

// GetTip 获取最新区块高度
func (bc *Blockchain) GetTip() (int, error) {
	data, err := bc.db.Get([]byte("chain_tip"), nil)
	if err != nil {
		return -1, err
	}
	// 简单解析字符串为 int
	var tip int
	fmt.Sscanf(string(data), "%d", &tip)
	return tip, nil
}

// AddBlock 添加新区块（持久化版本）
func (bc *Blockchain) AddBlock(data string) (*Block, error) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	// 1. 获取最新高度
	tip, err := bc.GetTip()
	if err != nil {
		return nil, err
	}

	// 2. 获取上一个区块（为了获取 PrevHash）
	prevBlock, err := bc.GetBlock(tip)
	if err != nil {
		return nil, err
	}

	// 3. 创建新区块
	newBlock := NewBlock(tip+1, data, prevBlock.Hash)
	fmt.Printf("正在挖掘区块 %d ...\n", newBlock.Index)
	newBlock.MineBlock(Difficulty)

	// 4. 保存到数据库
	err = bc.saveBlock(newBlock)
	if err != nil {
		return nil, err
	}

	fmt.Printf("区块 %d 挖掘成功并已持久化！\n", newBlock.Index)
	return newBlock, nil
}

// GetBlocks 获取所有区块（用于 API 展示，注意性能）
// 实际生产中不会一次性全取，这里为了兼容原有 API
func (bc *Blockchain) GetBlocks() ([]*Block, error) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	tip, err := bc.GetTip()
	if err != nil {
		return nil, err
	}

	blocks := make([]*Block, 0, tip+1)
	for i := 0; i <= tip; i++ {
		block, err := bc.GetBlock(i)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, block)
	}
	return blocks, nil
}

// IsChainValid 验证链条（逻辑不变，只是数据来源变了）
func (bc *Blockchain) IsChainValid() bool {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	tip, err := bc.GetTip()
	if err != nil {
		return false
	}

	for i := 1; i <= tip; i++ {
		currentBlock, err := bc.GetBlock(i)
		if err != nil {
			return false
		}
		prevBlock, err := bc.GetBlock(i - 1)
		if err != nil {
			return false
		}

		if currentBlock.Hash != currentBlock.CalculateHash() {
			return false
		}
		if currentBlock.PrevHash != prevBlock.Hash {
			return false
		}
	}
	return true
}

// SyncChain 同步更长的区块链到本地数据库
// 返回是否成功同步
func (bc *Blockchain) SyncChain(newChain []*Block) bool {
	if len(newChain) == 0 {
		return false
	}

	bc.mu.Lock()
	defer bc.mu.Unlock()

	// 1. 获取本地当前高度
	localTip, err := bc.GetTip()
	if err != nil {
		fmt.Printf("❌ 获取本地链高度失败: %v\n", err)
		return false
	}

	// 2. 如果新链不比本地长，不同步
	if len(newChain) <= localTip+1 {
		return false
	}

	// 3. 验证新链的完整性（从创世块开始验证）
	// 注意：实际生产中可以优化，只验证差异部分
	for i := 1; i < len(newChain); i++ {
		current := newChain[i]
		prev := newChain[i-1]

		// 验证哈希
		if current.Hash != current.CalculateHash() {
			fmt.Printf("❌ 同步失败: 区块 %d 哈希验证失败\n", current.Index)
			return false
		}

		// 验证链条连续性
		if current.PrevHash != prev.Hash {
			fmt.Printf("❌ 同步失败: 区块 %d 的 PrevHash 不匹配\n", current.Index)
			return false
		}

		// 验证难度（可选，简化处理）
		// 如果需要严格验证，可以检查哈希前缀是否符合难度要求
	}

	// 4. 验证通过，开始写入数据库
	// 策略：清空本地链，重新写入新链（简单但有效）
	// 生产环境可以只写入差异部分

	// 4.1 先备份当前数据库路径（可选，用于回滚）
	// 这里简化处理，直接覆盖

	// 4.2 逐个保存新区块
	// 注意：LevelDB 支持批量写入，性能更好
	batch := new(leveldb.Batch)

	for _, block := range newChain {
		data, err := json.Marshal(block)
		if err != nil {
			fmt.Printf("❌ 序列化区块 %d 失败: %v\n", block.Index, err)
			return false
		}
		blockKey := fmt.Sprintf("block:%d", block.Index)
		batch.Put([]byte(blockKey), data)
	}

	// 更新链顶
	batch.Put([]byte("chain_tip"), []byte(fmt.Sprintf("%d", len(newChain)-1)))

	// 4.3 一次性写入
	if err := bc.db.Write(batch, nil); err != nil {
		fmt.Printf("❌ 批量写入数据库失败: %v\n", err)
		return false
	}

	fmt.Printf("✅ 成功同步 %d 个区块，当前链高度: %d\n", len(newChain), len(newChain)-1)
	return true
}

// GetChainLength 获取当前区块链长度（辅助方法）
func (bc *Blockchain) GetChainLength() int {
	tip, err := bc.GetTip()
	if err != nil {
		return 0
	}
	return tip + 1 // 高度从 0 开始，长度 = 高度 + 1
}

// Close 关闭数据库连接
func (bc *Blockchain) Close() {
	if bc.db != nil {
		bc.db.Close()
	}
}
