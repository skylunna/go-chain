# go-chain
go-chain


## 🔒 安全演示

### 模拟黑客攻击

1. 启动服务器后，先挖掘几个区块
2. 调用 `/tamper` 接口篡改数据
3. 调用 `/valid` 接口验证，会检测到篡改

```bash
# 篡改区块 1 的数据
curl -X POST http://localhost:8080/tamper \
  -H "Content-Type: application/json" \
  -d '{"index": 1, "data": "HACKED!"}'

# 验证区块链
curl http://localhost:8080/valid
# 输出：区块链无效！检测到篡改！