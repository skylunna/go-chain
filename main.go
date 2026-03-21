package main

import (
	"fmt"
)

func main() {

	// // 1. 创建创世区块
	// // 创世区块的特点是索引为0，且没有前一个区块哈希(通常设置为 "0" 或空字符串)
	// genesisBlock := NewBlock(0, "Genesis Block", "")

	// // 2. 创建第二个区块
	// // 它的前一个哈希应该是创世区块的哈希
	// secondBlock := NewBlock(1, "Hello BlockChain", genesisBlock.Hash)

	// // 3. 打印信息
	// fmt.Println("=== Genesis Block ===")
	// fmt.Printf("Index: %d\n", genesisBlock.Index)
	// fmt.Printf("Data: %s\n", genesisBlock.Data)
	// fmt.Printf("Hash: %s\n", genesisBlock.Hash)
	// fmt.Printf("PrevHash: %s\n", genesisBlock.PrevHash)

	// fmt.Println("\n=== Second Block ===")
	// fmt.Printf("Index: %d\n", secondBlock.Index)
	// fmt.Printf("Data: %s\n", secondBlock.Data)
	// fmt.Printf("Hash: %s\n", secondBlock.Hash)
	// fmt.Printf("PrevHash: %s\n", secondBlock.PrevHash)


	// // 4. 验证链接
	// if secondBlock.PrevHash == genesisBlock.Hash {
	// 	fmt.Println("\n ✅ 验证成功: 第二个区块正确连接到了创世区块!")
	// } else {
	// 	fmt.Println("\n ❌ 验证失败，链条断裂")
	// }

	// 1. 创建区块链
	chain := NewBlockchain()

	// 2. 添加几个区块
	chain.AddBlock("Transaction 1: A send 10 BTC to B")
	chain.AddBlock("Transaction 2: B send 5 BTC to C")
	chain.AddBlock("Transaction 3: C send 3 BTC to D")

	// 3. 验证初始化状态
	fmt.Println("==== 初始状态验证 ====")
	if chain.IsChainValid() {
		fmt.Println("✅ 区块链验证通过，数据完整! ")
	} else {
		fmt.Println("❌ 区块链验证失败")
	}

	// 4. 演示篡改攻击
	fmt.Println("\n==== 模拟黑客篡改数据 ====")
	// 黑客偷偷修改第二个区块的数据 (索引1)
	chain.Blocks[1].Data = "HACKED! Hacker stole all coins!"
	// 注意：黑客无法重新计算哈希，因为需要大量算力(后续会实现PoW)
	hk := chain.Blocks[1].CalculateHash()
	fmt.Printf("👋 黑客: %s\n", hk)



	// 5. 再次验证
	fmt.Println("==== 篡改后验证 ====")
	if chain.IsChainValid() {
		fmt.Println("✅ 区块链验证通过，数据完整！")
	} else {
		fmt.Println("❌ 区块链验证失败！数据已经被篡改! ")
	}

	// 6. 打印被篡改的区块信息
	fmt.Println("\n====被篡改的区块详情====")
	tamperedBlock := chain.Blocks[1]
	fmt.Printf("存储的 Hash: %s\n", tamperedBlock.Hash)
	fmt.Printf("实际计算的 Hash: %s\n", tamperedBlock.CalculateHash())
	fmt.Printf("两个哈希是否匹配：%v\n", tamperedBlock.Hash == tamperedBlock.CalculateHash())
}