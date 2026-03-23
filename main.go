package main

import (
	"fmt"
	"log"
	"net/http"
	// "time"
)

func main() {

	P2PDemo()
}

// func NetDemo() {
// 	// 1. 初始化区块链
// 	BlockChain = NewBlockchain()
// 	fmt.Println("✅ 区块链初始化完成，创世区块已挖掘")

// 	// 2. 注册路由
// 	// 查看整个链
// 	http.HandleFunc("/blockchain", handleGetBlockchain)

// 	// 挖矿/添加区块
// 	http.HandleFunc("/mine", handleMineBlock)

// 	// 验证链
// 	http.HandleFunc("/valid", handleIsValid)

// 	// 黑客攻击
// 	http.HandleFunc("/tamper", handleTamper)

// 	// 3. 启动服务器
// 	// 3. 启动服务器
// 	port := ":8080"
// 	fmt.Println("===========================================")
// 	fmt.Println("🚀 区块链服务器启动成功！")
// 	fmt.Println("===========================================")
// 	fmt.Println("📡 API 接口列表：")
// 	fmt.Println("   GET  /blockchain  - 查看完整区块链")
// 	fmt.Println("   POST /mine        - 挖掘新区块")
// 	fmt.Println("   GET  /valid       - 验证区块链完整性")
// 	fmt.Println("   POST /tamper      - 🚨 模拟黑客篡改数据")
// 	fmt.Println("===========================================")
// 	fmt.Printf("🌐 访问 http://localhost%s/blockchain\n", port)
// 	fmt.Println("===========================================")

// 	log.Fatal(http.ListenAndServe(port, nil))
// }

var P2P *P2PManager

func P2PDemo() {
	// 1. 加载配置
	config := LoadConfig()

	// 2. 初始化区块链
	BlockChain = NewBlockchain()
	fmt.Println("✅ 区块链初始化完成")

	// 3. 初始化P2P管理器
	P2P = NewP2PManager(config)

	// 4. 同步其他节点的链
	if len(config.Peers) > 0 {
		P2P.SyncWithPeers()
	}

	// 5. 注册路由
	http.HandleFunc("/blockchain", handleGetBlockchain)
	http.HandleFunc("/mine", handleMineBlock)
	http.HandleFunc("/valid", handleIsValid)
	http.HandleFunc("/tamper", handleTamper)
	http.HandleFunc("/block/receive", handleReceiveBlock) // 🆕 接收区块

	// 6. 启动服务器
	fmt.Println("===========================================")
	fmt.Println("🚀 P2P 区块链节点启动成功！")
	fmt.Println("===========================================")
	fmt.Printf("📍 当前端口：%s\n", config.Port)
	fmt.Printf("🔗 已知节点：%v\n", config.Peers)
	fmt.Println("===========================================")

	log.Fatal(http.ListenAndServe(config.Port, nil))
}
