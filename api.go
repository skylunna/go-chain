package main

import (
	"encoding/json"
	"net/http"
)

// 全局区块链实例（在main中初始化）
var BlockChain *Blockchain

// BlockResponse 用于 API返回的返回结构
type BlockResponse struct {
	Success bool	`json:"success"`
	Message string	`json:"message"`
	Data    interface{}	`json:"data,omitempty"`	// 如果值为空，则不输出这个字段
}


// 处理 GET/blockchain 请求
func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")	// 告诉浏览器，返回的是Json数据
	
	// 获取数据，从全局区块链变量中获取所有区块
	blocks := BlockChain.GetBlocks()

	response := BlockResponse{
		Success: true,
		Message: "区块链获取成功",
		Data: blocks,
	}

	/*
		HTTP/1.1 200 OK
		Content-Type: application/json
		Date: Mon, 01 Jan 2024 00:00:00 GMT
		Content-Length: 358

		{
		"success": true,
		"message": "区块链获取成功",
		"data": [
			{
			"Index": 0,
			"Timestamp": "2024-01-01 00:00:00",
			"Data": "Genesis Block",
			"PrevHash": "",
			"Hash": "3a5d8f2e...",
			"Nonce": 12345
			},
			{
			"Index": 1,
			"Timestamp": "2024-01-01 00:01:00",
			"Data": "Transaction 1",
			"PrevHash": "3a5d8f2e...",
			"Hash": "8b2c4f1a...",
			"Nonce": 67890
			}
		]
		}
	*/
	json.NewEncoder(w).Encode(response)
}

// 处理 POST/mine 请求
// 用法：curl -X POST http://localhost:8080/mine -d "data=hello"
func handleMineBlock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 简单获取POST数据
	// 实际项目中应该解析JSON body，这里为了简化直接读取
	// 如果没有数据，默认使用 "Transaction"
	data := "New Transaction"
	// 这里简化处理，实际应该读取 r.Body
	
	BlockChain.AddBlock(data)

	response := BlockResponse {
		Success: true,
		Message: "区块挖掘成功",
	}

	json.NewEncoder(w).Encode(response)
}

// 处理GET/Valid请求
func handleIsValid(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	isValid := BlockChain.IsChainValid()

	msg := "区块链有效"
	if !isValid {
		msg = "区块链无效！监测到篡改"
	}

	response := BlockResponse {
		Success: isValid,
		Message: msg,
	}

	json.NewEncoder(w).Encode(response)
}