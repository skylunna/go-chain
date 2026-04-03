package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/skylunna/go-chain/wallet"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("用法：wallet-cli <command>")
		fmt.Println("命令：")
		fmt.Println("  generate - 生成新钱包")
		fmt.Println("  import <private_key> - 导入私钥")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "generate":
		w, err := wallet.GenerateWallet()
		if err != nil {
			fmt.Printf("❌ 生成失败：%v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ 钱包生成成功！")
		fmt.Println("===========================================")
		fmt.Printf("📍 地址：%s\n", w.Address())
		fmt.Printf("🔑 私钥：%s\n", w.PrivateKeyToHex())
		fmt.Println("===========================================")
		fmt.Println("⚠️  请安全保管私钥，丢失后无法恢复！")

		walletData := map[string]string{
			"address":    w.Address(),
			"privateKey": w.PrivateKeyToHex(),
		}

		data, _ := json.MarshalIndent(walletData, "", " ")
		os.WriteFile("my/my_wallet.json", data, 0600) // 0600 表示只有所有者可读写
		fmt.Println("✅ 钱包已保存到 my/my_wallet.json")

	case "import":
		if len(os.Args) < 3 {
			fmt.Println("❌ 请提供私钥")
			os.Exit(1)
		}
		privateKeyHex := os.Args[2]
		privKey, err := wallet.HexToPrivateKey(privateKeyHex)
		if err != nil {
			fmt.Printf("❌ 导入失败：%v\n", err)
			os.Exit(1)
		}
		w := &wallet.Wallet{
			PrivateKey: privKey,
			PublicKey:  &privKey.PublicKey,
		}
		fmt.Println("✅ 私钥导入成功！")
		fmt.Printf("📍 地址：%s\n", w.Address())

	default:
		fmt.Printf("❌ 未知命令：%s\n", command)
		os.Exit(1)
	}
}
