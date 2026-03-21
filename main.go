package main

import (
	"fmt"
)

func main() {

	// 1. 创建创世区块
	// 创世区块的特点是索引为0，且没有前一个区块哈希(通常设置为 "0" 或空字符串)
	genesisBlock := NewBlock(0, "Genesis Block", "")

	// 2. 创建第二个区块
	// 它的前一个哈希应该是创世区块的哈希
	secondBlock := NewBlock(1, "Hello BlockChain", genesisBlock.Hash)

	// 3. 打印信息
	fmt.Println("=== Genesis Block ===")
	fmt.Printf("Index: %d\n", genesisBlock.Index)
	fmt.Printf("Data: %s\n", genesisBlock.Data)
	fmt.Printf("Hash: %s\n", genesisBlock.Hash)
	fmt.Printf("PrevHash: %s\n", genesisBlock.PrevHash)

	fmt.Println("\n=== Second Block ===")
	fmt.Printf("Index: %d\n", secondBlock.Index)
	fmt.Printf("Data: %s\n", secondBlock.Data)
	fmt.Printf("Hash: %s\n", secondBlock.Hash)
	fmt.Printf("PrevHash: %s\n", secondBlock.PrevHash)


	// 4. 验证链接
	if secondBlock.PrevHash == genesisBlock.Hash {
		fmt.Println("\n ✅ 验证成功: 第二个区块正确连接到了创世区块!")
	} else {
		fmt.Println("\n ❌ 验证失败，链条断裂")
	}
}