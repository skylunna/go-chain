package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// 全局区块链实例（在main中初始化）
var BlockChain *Blockchain

// BlockResponse 用于 API返回的返回结构
type BlockResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"` // 如果值为空，则不输出这个字段
}

// 接收篡改请求的结构
type TamperRequest struct {
	Index int    `json:"index"`
	Data  string `json:"data"`
}

// BlockDTO 用于 API 传输的区块数据结构
type BlockDTO struct {
	Index     int    `json:"Index"`
	Timestamp string `json:"Timestamp"`
	Data      string `json:"Data"`
	PrevHash  string `json:"PrevHash"`
	Hash      string `json:"Hash"`
	Nonce     int    `json:"Nonce"`
}

// BlockToDTO 将 Block 转换为 DTO
func BlockToDTO(b *Block) BlockDTO {
	return BlockDTO{
		Index:     b.Index,
		Timestamp: b.Timestamp,
		Data:      b.Data,
		PrevHash:  b.PrevHash,
		Hash:      b.Hash,
		Nonce:     b.Nonce,
	}
}

// 将 DTO 转换回 Block
func DTOToBlock(dto BlockDTO) *Block {
	return &Block{
		Index:     dto.Index,
		Timestamp: dto.Timestamp,
		Data:      dto.Data,
		PrevHash:  dto.PrevHash,
		Hash:      dto.Hash,
		Nonce:     dto.Nonce,
	}
}

// 处理 GET/blockchain 请求
func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // 告诉浏览器，返回的是Json数据
	blocks := BlockChain.GetBlocks()

	// 转换为 DTO 格式
	dtoBlocks := make([]BlockDTO, len(blocks))
	for i, b := range blocks {
		dtoBlocks[i] = BlockToDTO(b)
	}

	response := BlockResponse{
		Success: true,
		Message: "区块链获取成功",
		Data:    dtoBlocks,
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

	response := BlockResponse{
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

	response := BlockResponse{
		Success: isValid,
		Message: msg,
	}

	json.NewEncoder(w).Encode(response)
}

// handleTamper 处理POST/tamper请求-模拟黑客攻击
func handleTamper(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// ✅ 正确：解析请求体用 Decoder + r.Body
	var req TamperRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := BlockResponse{
			Success: false,
			Message: "请求格式错误，需要 {\"index\": 1, \"data\": \"hacked\"}",
		}
		// ✅ 正确：写入响应体用 Encoder + w
		json.NewEncoder(w).Encode(response)
		return
	}

	// 检查索引是否合法
	if req.Index <= 0 || req.Index >= len(BlockChain.Blocks) {
		response := BlockResponse{
			Success: false,
			Message: fmt.Sprintf("无效的区块索引，有效范围：1-%d", len(BlockChain.Blocks)-1),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// 禁止篡改创世区块
	if req.Index == 0 {
		response := BlockResponse{
			Success: false,
			Message: "⚠️  创世区块不可篡改！",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// 保存原始数据（用于展示对比）
	originalData := BlockChain.Blocks[req.Index].Data
	originalHash := BlockChain.Blocks[req.Index].Hash

	// 🚨 执行篡改：直接修改数据，但不重新计算哈希
	BlockChain.Blocks[req.Index].Data = req.Data

	response := BlockResponse{
		Success: true,
		Message: "⚠️  数据篡改成功！但哈希值未更新，区块链已检测到异常！",
		Data: map[string]interface{}{
			"blockIndex":     req.Index,
			"originalData":   originalData,
			"tamperedData":   req.Data,
			"storedHash":     originalHash,
			"calculatedHash": BlockChain.Blocks[req.Index].CalculateHash(),
			"hashMatch":      originalHash == BlockChain.Blocks[req.Index].CalculateHash(),
			"nextStep":       "请调用 /valid 接口验证区块链完整性",
		},
	}
	// ✅ 正确：写入响应
	json.NewEncoder(w).Encode(response)
}

// 处理POST/block/receive请求 - 接收其他节点的区块
func handleReceiveBlock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var dto BlockDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response := BlockResponse{
			Success: false,
			Message: "无效的区块数据",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	newBlock := DTOToBlock(dto)

	// 验证区块
	BlockChain.mu.Lock()
	defer BlockChain.mu.Unlock()

	lastBlock := BlockChain.Blocks[len(BlockChain.Blocks)-1]

	if newBlock.Index != lastBlock.Index+1 {
		json.NewEncoder(w).Encode(BlockResponse{
			Success: false,
			Message: "区块索引不连续",
		})
		return
	}

	if newBlock.PrevHash != lastBlock.Hash {
		json.NewEncoder(w).Encode(BlockResponse{
			Success: false,
			Message: "前哈希不匹配",
		})
		return
	}

	if newBlock.Hash != newBlock.CalculateHash() {
		json.NewEncoder(w).Encode(BlockResponse{
			Success: false,
			Message: "区块哈希验证失败",
		})
		return
	}

	BlockChain.Blocks = append(BlockChain.Blocks, newBlock)

	json.NewEncoder(w).Encode(BlockResponse{
		Success: true,
		Message: fmt.Sprintf("区块 %d 接收成功", newBlock.Index),
	})
}
