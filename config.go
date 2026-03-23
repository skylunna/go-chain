package main

import (
	"os"
)

type NodeConfig struct {
	Port  string   // 当前节点端口
	Peers []string // 已知的其他节点地址
}

// LoadConfig 从环境变量加载配置
func LoadConfig() *NodeConfig {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 从环境变量获取其它节点地址，用逗号分隔
	peersStr := os.Getenv("PEERS")
	var peers []string
	if peersStr != "" {
		// 简单分割
		// 生产环境需要更健壮的解析
		for _, p := range splitString(peersStr, ",") {
			if p != "" {
				peers = append(peers, p)
			}
		}
	}

	return &NodeConfig{
		Port:  ":" + port,
		Peers: peers,
	}
}

// 简单的字符串分割函数 (避免引入额外依赖)
func splitString(s, sep string) []string {
	var result []string
	start := 0
	for i := 0; i <= len(s)-len(sep); i++ {
		if s[i:i+len(sep)] == sep {
			result = append(result, s[start:i])
			start = i + len(sep)
		}
	}

	result = append(result, s[start:])
	return result
}
