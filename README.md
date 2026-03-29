# Go-Chain 🚀

一个使用 Go 语言实现的去中心化区块链原型，包含工作量证明 (PoW)、P2P 网络通信和防篡改验证机制。

## ✨ 核心特性

- 🔒 **防篡改**：基于 SHA256 哈希链，任何数据修改都会导致验证失败
- ⛏️ **工作量证明 (PoW)**：模拟比特币挖矿机制，可调整难度
- 🌐 **P2P 网络**：支持多节点部署，自动同步最长链，实时广播新区块
- 🛡️ **安全演示**：内置黑客攻击模拟接口，直观展示区块链安全性
- 🐳 **Docker 支持**：一键启动多节点集群

## 🚀 快速开始

### 本地运行
```bash
go mod tidy
go run .
```

### Docker 启动多节点集群
```bash
docker-compose up
```

### API 接口
GET     /blockchain   获取完整区块链数据
POST    /mine         挖掘新区块
GET     /valid        验证链条完整性
POST    /tamper       模拟黑客篡改数据


### 安全演示
尝试篡改数据并观察验证失败:
```bash
curl -X POST http://localhost:8080/tamper -d '{"index":1, "data":"Hacked"}'
curl http://localhost:8080/valid
```

### 技术栈
- Go 1.21+
- HTTP/JSON API
- SHA256 Cryptography
- Goroutine Concurrency