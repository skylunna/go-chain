package main

// Blockchain 代表整个区块链
type Blockchain struct {
	Blocks []*Block
}

// NewBlockchain 创建一个新的区块链，包含创世区块
func NewBlockchain() *Blockchain {
	return &Blockchain {
		Blocks: []*Block{NewBlock(0, "Genesis Block", "")},
	}
}

// 添加新区块到链上
func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(prevBlock.Index+1, data, prevBlock.Hash)
	bc.Blocks = append(bc.Blocks, newBlock)
}

// IsChainValid 验证整个区块链是否有效
// 返回true表示链条完整未被篡改
func (bc *Blockchain) IsChainValid() bool {
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