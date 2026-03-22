package main

import (
	"fmt"
	// "time"
	// "log"
	"net/http"
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

	// // 1. 创建区块链
	// chain := NewBlockchain()

	// // 2. 添加几个区块
	// chain.AddBlock("Transaction 1: A send 10 BTC to B")
	// chain.AddBlock("Transaction 2: B send 5 BTC to C")
	// chain.AddBlock("Transaction 3: C send 3 BTC to D")

	// // 3. 验证初始化状态
	// fmt.Println("==== 初始状态验证 ====")
	// if chain.IsChainValid() {
	// 	fmt.Println("✅ 区块链验证通过，数据完整! ")
	// } else {
	// 	fmt.Println("❌ 区块链验证失败")
	// }

	// // 4. 演示篡改攻击
	// fmt.Println("\n==== 模拟黑客篡改数据 ====")
	// // 黑客偷偷修改第二个区块的数据 (索引1)
	// chain.Blocks[1].Data = "HACKED! Hacker stole all coins!"
	// // 注意：黑客无法重新计算哈希，因为需要大量算力(后续会实现PoW)
	// hk := chain.Blocks[1].CalculateHash()
	// fmt.Printf("👋 黑客: %s\n", hk)



	// // 5. 再次验证
	// fmt.Println("==== 篡改后验证 ====")
	// if chain.IsChainValid() {
	// 	fmt.Println("✅ 区块链验证通过，数据完整！")
	// } else {
	// 	fmt.Println("❌ 区块链验证失败！数据已经被篡改! ")
	// }

	// // 6. 打印被篡改的区块信息
	// fmt.Println("\n====被篡改的区块详情====")
	// tamperedBlock := chain.Blocks[1]
	// fmt.Printf("存储的 Hash: %s\n", tamperedBlock.Hash)
	// fmt.Printf("实际计算的 Hash: %s\n", tamperedBlock.CalculateHash())
	// fmt.Printf("两个哈希是否匹配：%v\n", tamperedBlock.Hash == tamperedBlock.CalculateHash())

	// fmt.Println("🚀 启动 Go-Chain 挖矿演示...")

	// // 1. 创建区块链 (会自动挖创世区块)
	// chain := NewBlockchain()

	// // 2. 添加区块并观察耗时
	// transactions := []string{
	// 	"A 转账 10 BTC 给 B",
	// 	"B 转账 5 BTC 给 C",
	// 	"C 转账 3 BTC 给 D",
	// }

	// for _, tx := range transactions {
	// 	start := time.Now()
	// 	chain.AddBlock(tx)
	// 	elapsed := time.Since(start)
	// 	fmt.Printf("⏰ 耗时: %v\n\n", elapsed)
	// }

	// // 最终验证
	// fmt.Println("==== 最终验证 ====")
	// if chain.IsChainValid() {
	// 	fmt.Println("✅ 区块链验证通过，所有区块均有效且满足难度要求！")
	// } else {
	// 	fmt.Println("❌ 区块链验证失败！")
	// }

	NetDemo()
}

func NetDemo() {
	// 1. 初始化区块链
	BlockChain = NewBlockchain()
	fmt.Println("✅ 区块链初始化完成，创世区块已挖掘")

	// 2. 注册路由
	// 查看整个链
	http.HandleFunc("/blockchain", handleGetBlockchain)

	// 挖矿/添加区块
	http.HandleFunc("/mine", handleMineBlock)

	// 验证链
	http.HandleFunc("/valid", handleIsValid)

	// 3. 启动服务器
	port := ":8080"
	fmt.Printf("🚀 服务器启动中... 访问 http://localhost%s/blockchain\n", port)

	// 启动HTTP服务器
	// 注意：log.Fatal会在服务器出错时终止程序
	// 修改后：添加更清晰的错误提示
	fmt.Printf("🚀 尝试在 %s 启动服务器...\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("❌ 服务器启动失败: %v\n", err)
		fmt.Println("💡 可能的原因：")
		fmt.Println("   1. 端口 8080 被其他程序占用")
		fmt.Println("   2. 没有权限绑定该端口")
		fmt.Println("   3. 防火墙阻止")
		return
	}
}