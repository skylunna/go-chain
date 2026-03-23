package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// p2p 节点通信管理器
type P2PManager struct {
	config *NodeConfig
}

// NewP2PManager 创建 P2P 管理器
func NewP2PManager(config *NodeConfig) *P2PManager {
	// 构造方法
	return &P2PManager{config: config}
}

// 向所有已知节点广播新区块
func (p *P2PManager) BroadcastBlock(block *Block) {
	fmt.Printf("🔈正在向 %d 个节点广播区块 %d ... \n", len(p.config.Peers), block.Index)

	for _, peer := range p.config.Peers {
		go p.sendBlockToPeer(peer, block)
	}
}

// 发送区块到单个节点
func (p *P2PManager) sendBlockToPeer(peer string, block *Block) {
	url := "http://" + peer + "/block/receive"

	/**
	把这个对象转成JSON字符串 block := &Block{Index: 1, Timestamp: 123456, Data: "转账100"}
						-> {"index":1,"timestamp":123456,"data":"转账100"}

	再把字符串转成 []byte 字节类型 (网络传输必须用字节)
						-> []byte(`{"index":1,"timestamp":123456,"data":"转账100"}`)
	*/
	dto := BlockToDTO(block)
	data, _ := json.Marshal(dto)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Printf("❌ 发送失败到 %s: %v\n", peer, err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("✅ 区块 %d 已发送到 %s\n", block.Index, peer)
}

// SyncWithPeers 启动时从其他节点同步区块链
func (p *P2PManager) SyncWithPeers() {
	fmt.Println("🔄 正在从其他节点同步区块链...")

	var longestChain []*Block

	for _, peer := range p.config.Peers {
		chain := p.fetchChainFromPeer(peer)
		if len(chain) > len(longestChain) {
			longestChain = chain
		}
	}

	if len(longestChain) > len(BlockChain.Blocks) {
		fmt.Printf("📥 发现更长的链 (%d 个区块)，正在同步...\n", len(longestChain))
		BlockChain.Blocks = longestChain
		fmt.Println("✅ 同步完成！")
	} else {
		fmt.Println("✅ 当前链已是最长链")
	}
}

// 从单个节点获取区块链
func (p *P2PManager) fetchChainFromPeer(peer string) []*Block {
	url := "http://" + peer + "/blockchain"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("❌ 无法连接 %s: %v\n", peer, err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("❌ 读取响应失败 %s: %v\n", peer, err)
		return nil
	}

	// 解析API响应
	var apiResp BlockResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		fmt.Printf("❌ 解析 JSON 失败 %s: %v\n", peer, err)
		return nil
	}

	// 将 Data 转换为 BlockDTO 数组
	dataBytes, _ := json.Marshal(apiResp.Data)
	var dtoBlocks []BlockDTO
	if err := json.Unmarshal(dataBytes, &dtoBlocks); err != nil {
		fmt.Printf("❌ 转换区块数据失败 %s: %v\n", peer, err)
		return nil
	}

	// 转换为 Block 数组
	result := make([]*Block, len(dtoBlocks))
	for i, dto := range dtoBlocks {
		result[i] = DTOToBlock(dto)
	}

	fmt.Printf("📥 从 %s 获取到 %d 个区块\n", peer, len(result))
	return result
}
