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
	data, _ := json.Marshal(block)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Printf("❌ 发送失败到 %s: %v\n", peer, err)
		return
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("❌ 节点 %s 返回错误 [%d]: %s\n", peer, resp.StatusCode, string(body))
		return
	}

	// 解析响应体确认业务逻辑成功
	var apiResp BlockResponse
	json.NewDecoder(resp.Body).Decode(&apiResp)
	if apiResp.Success {
		fmt.Printf("✅ 区块 %d 已被 %s 确认接收\n", block.Index, peer)
	} else {
		fmt.Printf("⚠️  区块 %d 被 %s 拒绝: %s\n", block.Index, peer, apiResp.Message)
	}
}

// SyncWithPeers 启动时从其他节点同步区块链
func (p *P2PManager) SyncWithPeers() {
	fmt.Println("🔄 正在从其他节点同步区块链...")

	// 1. 获取本地链长度
	localLength := BlockChain.GetChainLength()
	fmt.Printf("📊 本地链长度: %d\n", localLength)

	var longestChain []*Block
	var bestPeer string

	// 2. 向所有已知节点请求链
	for _, peer := range p.config.Peers {
		fmt.Printf("🔍 查询节点 %s 的链...\n", peer)
		chain := p.fetchChainFromPeer(peer)

		if chain != nil && len(chain) > len(longestChain) {
			longestChain = chain
			bestPeer = peer
		}
	}

	// 3. 如果找到了更长的链，尝试同步
	if longestChain != nil && len(longestChain) > localLength {
		fmt.Printf("📥 从 %s 发现更长的链 (%d 个区块 vs 本地 %d 个)\n",
			bestPeer, len(longestChain), localLength)

		if BlockChain.SyncChain(longestChain) {
			fmt.Println("✅ 同步完成！")
		} else {
			fmt.Println("❌ 同步失败，保留本地链")
		}
	} else {
		fmt.Println("✅ 当前链已是最长链，无需同步")
	}
}

// 从单个节点获取区块链
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

	// 解析 API 响应
	var apiResp BlockResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		fmt.Printf("❌ 解析 JSON 失败 %s: %v\n", peer, err)
		return nil
	}

	// ✅ 优化：直接解析为 []BlockDTO，避免中间转换
	// 因为 apiResp.Data 已经是 json.RawMessage 或类似类型
	// 我们可以直接复用 json 解码
	dtoBlocks := make([]BlockDTO, 0)

	// 方法：将 Data 字段重新编码再解码
	if apiResp.Data != nil {
		dataBytes, err := json.Marshal(apiResp.Data)
		if err != nil {
			fmt.Printf("❌ 序列化区块数据失败 %s: %v\n", peer, err)
			return nil
		}
		if err := json.Unmarshal(dataBytes, &dtoBlocks); err != nil {
			fmt.Printf("❌ 解析区块数组失败 %s: %v\n", peer, err)
			return nil
		}
	}

	// 转换为 Block 数组
	result := make([]*Block, len(dtoBlocks))
	for i, dto := range dtoBlocks {
		result[i] = DTOToBlock(dto)
	}

	fmt.Printf("📥 从 %s 获取到 %d 个区块\n", peer, len(result))
	return result
}
