package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var P2P *P2PManager

func main() {
	// 1. 加载配置
	config := LoadConfig()

	// 2. 初始化数据库路径 (默认当前目录下的 chain_data)
	dbPath := "./chain_data"
	// 可以通过环境变量自定义
	if os.Getenv("DB_PATH") != "" {
		dbPath = os.Getenv("DB_PATH")
	}

	// 3. 初始化区块链 (持久化)
	var err error
	BlockChain, err = InitBlockchain(dbPath)
	if err != nil {
		log.Fatalf("❌ 初始化区块链失败：%v", err)
	}

	// 4. 优雅退出处理 (确保关闭数据库)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\n🛑 正在关闭数据库...")
		BlockChain.Close()
		os.Exit(0)
	}()

	fmt.Println("✅ 区块链持久化加载完成")

	// 5. 初始化 P2P 管理器
	P2P = NewP2PManager(config)

	// 6. 同步其他节点的链 (略，逻辑不变)
	if len(config.Peers) > 0 {
		// 注意：P2P 同步逻辑也需要适配数据库版本，暂时简化处理
		fmt.Println("🔄 跳过初始同步以简化演示，实际需实现 DB 同步逻辑")
	}

	// 7. 注册路由 (略，保持不变)
	http.HandleFunc("/blockchain", handleGetBlockchain)
	http.HandleFunc("/mine", handleMineBlock)
	http.HandleFunc("/valid", handleIsValid)
	http.HandleFunc("/tamper", handleTamper)
	http.HandleFunc("/block/receive", handleReceiveBlock)

	// 8. 启动服务器
	fmt.Println("===========================================")
	fmt.Println("🚀 持久化区块链节点启动成功！")
	fmt.Printf("📂 数据目录：%s\n", dbPath)
	fmt.Println("===========================================")

	log.Fatal(http.ListenAndServe(config.Port, nil))
}
