# Time Ledger - 多链 ERC20 积分系统

基于 Go 开发的多链 ERC20 代币事件追踪与积分计算系统，支持以太坊和 OP Stack 链。

## 📋 核心功能

- **智能合约**：ERC20 代币合约（支持 mint/burn）
- **事件索引器**：实时追踪链上 Transfer 事件，维护用户余额
- **积分计算器**：基于持仓时间和余额计算积分（时间加权）
- **多链支持**：支持 Sepolia、Base Sepolia 等多条链并发同步
- **分叉处理**：OP Stack 链的区块重组检测与自动回滚
- **容错恢复**：断点续传和历史积分回溯

## 🏗️ 项目架构

```
time-ledger/
├── timeledger-contracts/          # 智能合约（Foundry）
│   ├── src/TimeLedgerToken.sol   # ERC20 合约
│   ├── script/Deploy.s.sol       # 部署脚本
│   └── abi/timeledgertoken.go    # Go 绑定
│
└── timeledger-backend/            # Go 后端服务
    ├── cmd/server/                # 主程序入口
    ├── internal/
    │   ├── service/indexer/       # 事件索引器
    │   ├── service/calculator/    # 积分计算器
    │   ├── repository/            # 数据访问层
    │   └── api/                   # HTTP API
    └── configs/config.toml        # 多链配置
```

## 🚀 快速开始

### 前置要求

- Go 1.24+、Foundry、MySQL 8.0+、Redis 6.0+
- RPC 节点访问权限（Alchemy、Infura 等）

### 1. 部署智能合约

```bash
cd timeledger-contracts
forge install && forge build

# 配置环境变量
cp .env.example .env
# 编辑 .env，填入 DEPLOYER_PRIVATE_KEY 和 RPC_URL

# 部署到 Sepolia
forge script script/Deploy.s.sol:Deploy \
  --rpc-url $SEPOLIA_RPC_URL \
  --broadcast --verify
```

### 2. 生成 Go 绑定

```bash
abigen --abi out/TimeLedgerToken.sol/TimeLedgerToken.json \
  --pkg erc20 \
  --type TimeLedgerToken \
  --out timeledger-backend/pkg/contract/erc20/timeledgertoken.go
```

### 3. 配置后端服务

```bash
cd timeledger-backend
go mod download

# 配置数据库和 RPC
cp .env.example .env
# 编辑 .env 和 configs/config.toml
```

配置示例 `configs/config.toml`：

```toml
[[chains]]
name = "sepolia"
chain_id = 11155111
type = "ethereum"
rpc_env_key = "SEPOLIA_RPC_URL"
confirmations = 6
chunk_size = 10
request_delay_ms = 100

[[chains.contracts]]
address = "0xYourContractAddress"
start_block = 10032808
token_decimals = 18
```

### 4. 启动服务

```bash
# 启动所有服务（推荐）
go run cmd/server/main.go all

# 或分别启动
go run cmd/server/main.go api        # 仅 API
go run cmd/server/main.go indexer    # 仅索引器
go run cmd/server/main.go calculator # 仅积分计算器
```

### 5. 测试合约交互

```bash
# Mint 代币
cast send $CONTRACT_ADDRESS \
  "mint(address,uint256)" \
  $USER_ADDRESS 100000000000000000000 \
  --rpc-url $SEPOLIA_RPC_URL --private-key $PRIVATE_KEY

# 转账代币
cast send $CONTRACT_ADDRESS \
  "transfer(address,uint256)" \
  $RECIPIENT_ADDRESS 50000000000000000000 \
  --rpc-url $SEPOLIA_RPC_URL --private-key $PRIVATE_KEY

# 查询余额
cast call $CONTRACT_ADDRESS \
  "balanceOf(address)(uint256)" \
  $USER_ADDRESS --rpc-url $SEPOLIA_RPC_URL
```

## 📊 数据库设计

### 核心表结构

| 表名 | 说明 | 关键字段 |
|------|------|----------|
| `block_cursor` | 同步进度 | chain_id, contract_address, block_number |
| `balance_log` | 余额变动事实表 | account, delta, balance_after, tx_hash |
| `user_balance` | 用户余额快照 | account, balance, block_number |
| `user_point` | 用户积分快照 | account, total_points, last_calc_time |
| `user_point_log` | 积分计算明细 | account, from_time, to_time, points |

> 📖 **详细示例**：查看 [数据库示例.md](./数据库示例.md) 了解完整的数据示例、业务逻辑和一致性验证

## 🔧 核心机制

### 1. 区块确认机制

**以太坊链**：使用 `confirmations` 参数（默认 6 个区块），简单线性同步

**OP Stack 链**：使用 `reorg_window` 参数（默认 200 个区块）
- Redis 缓存 pending 区块
- 定期检测区块重组
- 自动回滚到 safe block

### 2. 积分计算公式

```
积分 = Σ (余额 × 持有时长 × 费率)
```

**示例**：
```
15:00 - 余额 0
15:10 - 余额 100
15:30 - 余额 200

16:00 计算积分：
= 100 × (20分钟/60) × 0.05 + 200 × (30分钟/60) × 0.05
= 1.665 + 5.0 = 6.665 积分
```

### 3. 分叉检测与回滚

```go
// 检测分叉：对比数据库与链上区块哈希
EnsureCanonicalOrRollback(ctx, client, chainID, contract, reorgWindow)

// 执行回滚：删除 > safeBlock 的数据，重建余额快照
RollbackTo(ctx, chainID, contract, safeBlock)
```

### 4. 历史积分回溯

程序中断后重启，自动补算历史积分：
- 读取 `last_calc_time`（如 1月1日）
- 读取 `safe_block_time`（如 1月5日）
- 计算区间积分并按小时拆分生成明细

## 📝 API 接口

```bash
# 查询用户余额
GET /api/balance/:chain_id/:contract/:account

# 查询用户积分
GET /api/points/:chain_id/:contract/:account

# 查询余额变动历史
GET /api/balance/history/:chain_id/:contract/:account

# 查询积分明细
GET /api/points/history/:chain_id/:contract/:account
```

## 🔐 安全性保障

- **余额一致性**：数据库事务 + 唯一索引防重 + 负余额检测
- **积分准确性**：基于事实表计算 + 行锁防并发 + 幂等性设计
- **分叉安全**：safe block 机制 + 定期检测 + 自动回滚
- **RPC 容错**：限流器 + 指数退避重试 + 错误分类处理

## 🛠️ 运维建议

### 监控指标

- **Indexer**：同步延迟、RPC 错误率、分叉检测次数
- **Calculator**：积分计算延迟、计算失败次数、回溯时长
- **数据库**：连接池使用率、慢查询、表大小增长

### 故障排查

**问题 1：Indexer 停止同步（RPC 限流）**
- 增加 `request_delay_ms` 参数
- 降低 `chunk_size` 参数
- 升级 RPC 节点套餐

**问题 2：积分计算延迟**
- 检查 Indexer 是否正常运行
- 检查 Calculator 是否启动
- 手动触发积分计算

**问题 3：余额不一致（分叉未检测）**
- 检查 `block_header` 表是否有重复区块
- 手动触发回滚：`RollbackTo(safeBlock)`

## 📈 性能优化

### 已实现

- 多链并发同步（errgroup）
- 批量拉取事件（chunk_size）
- Redis 缓存 pending 区块
- 数据库索引优化
- 连接池配置（max_open_conns=50）

### 待优化方向

> 本项目为学习实践项目，以下为可迭代方向：

1. **日志系统模块化**：引入 zap/zerolog 实现结构化日志
2. **数据访问层抽象**：引入 Repository 模式解耦业务逻辑
3. **链级并发优化**：实现合约级并发和区块级并发
4. **监控告警**：集成 Prometheus + Grafana
5. **缓存优化**：引入本地缓存减少 Redis 访问
6. **RPC 优化**：实现节点池和 WebSocket 订阅

## 🧪 测试

```bash
# 合约测试
cd timeledger-contracts && forge test -vvv

# 后端测试
cd timeledger-backend && go test ./...
```

## 📄 许可证

MIT License
