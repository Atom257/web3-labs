# Time Ledger - 多链 ERC20 时间加权积分系统

> 基于持仓时间和余额计算积分的区块链索引系统

> **该项目以“时间加权积分”为业务载体，用于验证多链 Indexer 在并发、Reorg、幂等性和数据一致性场景下的工程实现。**

## 👀 Infra 速读（3–5 分钟）

如果你关注的是 Web3 Infra / Indexer 相关能力，可重点查看以下部分：

- **多链并发 Indexer 实现**：Part 2 → Indexer（事件索引器）
- **区块重组（Reorg）处理与回滚机制**：Part 2 → 区块重组处理流程
- **幂等性与一致性保障**：2.3 并发安全机制
- **大规模数据下的分表设计**：为什么 user_point_log 要分表？

  
## 📖 目录
- [项目简介](#项目简介)
- [核心特性](#核心特性)
- [系统架构](#系统架构)
- [Part 1: 智能合约](#part-1-智能合约)
- [Part 2: Go 后端服务](#part-2-go-后端服务)
- [数据库设计](#数据库设计)
- [快速开始](#快速开始)
- [API 文档](#api-文档)

---

## 项目简介

Time Ledger 是一个多链 ERC20 代币事件追踪与积分计算系统。

**核心功能**：
- ✅ 实时追踪链上 Transfer 事件
- ✅ 基于持仓时间计算积分（时间加权）
- ✅ 多链并发同步（Ethereum、OP Stack）
- ✅ 自动处理区块重组（Reorg）
- ✅ 断点续传和历史回溯

**积分计算公式**：
```
积分 = Σ (余额 × 持有时长 × 费率)
```

**示例**：
- 用户持有 100 代币，持续 24 小时，费率 5%
- 积分 = 100 × (24/24/365) × 0.05 = 0.0137 积分

---

## 核心特性

### 🔗 多链支持
- Ethereum 主网/测试网（Sepolia）
- OP Stack 链（Base Sepolia、Optimism）
- 可扩展至任意 EVM 兼容链

### 🛡️ 数据一致性保障
- **幂等性设计**：唯一索引防止重复记录
- **区块重组处理**：自动检测并回滚到安全区块
- **双重验证**：余额与积分可通过事实表重建验证

### ⚡ 高性能设计
- **多链并发同步**：使用 errgroup 并发处理多条链
- **批量事件拉取**：chunk_size 控制批量大小
- **Redis 缓存**：缓存 pending 区块，减少 RPC 调用
- **数据库连接池**：max_open_conns=50，max_idle_conns=10

### 🔐 并发安全保障
- **数据库事务**：使用 GORM 事务确保原子性
- **唯一索引**：防止并发插入重复数据
- **行级锁**：积分计算时使用 SELECT ... FOR UPDATE
- **幂等操作**：所有写操作支持重复执行

---

## 系统架构

```
┌─────────────────────────────────────────────────────────┐
│                     Blockchain Layer                     │
│  Sepolia (11155111)  │  Base Sepolia (84532)  │  ...    │
└────────────┬─────────────────────┬──────────────────────┘
             │                     │
             ▼                     ▼
┌─────────────────────────────────────────────────────────┐
│                   Indexer Service (Go)                   │
│  • 监听 Transfer 事件                                     │
│  • 维护余额快照                                           │
│  • 检测区块重组                                           │
│  • 并发处理多链                                           │
└────────────┬────────────────────────────────────────────┘
             │
             ▼
┌─────────────────────────────────────────────────────────┐
│                  MySQL Database                          │
│  balance_log │ user_balance │ block_cursor │ ...        │
│  user_point_log_1 (Sepolia) │ user_point_log_2 (Base)  │
└────────────┬────────────────────────────────────────────┘
             │
             ▼
┌─────────────────────────────────────────────────────────┐
│                Calculator Service (Go)                   │
│  • 基于余额变动计算积分                                   │
│  • 支持历史回溯                                           │
│  • 按时间段拆分明细                                       │
│  • 使用行锁防止并发冲突                                   │
└────────────┬────────────────────────────────────────────┘
             │
             ▼
┌─────────────────────────────────────────────────────────┐
│                    HTTP API Server                       │
│  /api/balance  │  /api/points  │  /api/history          │
└─────────────────────────────────────────────────────────┘
```

---

# Part 1: 智能合约

## 1.1 合约概述

**TimeLedgerToken.sol** - 标准 ERC20 代币合约

**核心功能**：
- ✅ Mint/Burn 功能
- ✅ 标准 Transfer/Approve
- ✅ 事件日志完整

**已部署合约**：
- **Sepolia**: `0xBEfe9d9726c3BFD513b6aDd74B243a82b272C073`
- **Base Sepolia**: `0xB8a31EaC0874DC6f5a28FCa601336Ae32c723dF6`

## 1.2 快速部署

```bash
cd timeledger-contracts

# 安装依赖
forge install

# 配置环境变量
cp .env.example .env
# 编辑 .env: DEPLOYER_PRIVATE_KEY, SEPOLIA_RPC_URL, ETHERSCAN_API_KEY

# 部署到 Sepolia
forge script script/Deploy.s.sol:Deploy \
  --rpc-url $SEPOLIA_RPC_URL \
  --broadcast --verify

# 生成 Go 绑定
./scripts/generate-abi.sh
```

## 1.3 合约交互示例

```bash
# Mint 代币
cast send $CONTRACT_ADDRESS \
  "mint(address,uint256)" \
  $USER_ADDRESS 100000000000000000000 \
  --rpc-url $SEPOLIA_RPC_URL --private-key $PRIVATE_KEY

# 转账
cast send $CONTRACT_ADDRESS \
  "transfer(address,uint256)" \
  $RECIPIENT 50000000000000000000 \
  --rpc-url $SEPOLIA_RPC_URL --private-key $PRIVATE_KEY

# 查询余额
cast call $CONTRACT_ADDRESS \
  "balanceOf(address)(uint256)" \
  $USER_ADDRESS --rpc-url $SEPOLIA_RPC_URL
```

---

> 本部分重点展示多链 Indexer 在真实运行环境下的并发模型、  
> 区块重组（Reorg）处理逻辑，以及高并发写入场景下的数据一致性保障。

# Part 2: Go 后端服务

## 2.1 项目结构

```
timeledger-backend/
├── cmd/
│   └── server/
│       └── main.go              # 程序入口
├── internal/
│   ├── config/                  # 配置加载
│   ├── models/                  # 数据模型
│   ├── repository/              # 数据访问层
│   │   ├── db.go                # 数据库初始化
│   │   ├── redis.go             # Redis 初始化
│   │   ├── system_repo.go       # 系统初始化（建表+数据）
│   │   ├── contract_repo.go     # 合约配置
│   │   └── point_rate_repo.go   # 费率配置
│   ├── service/
│   │   ├── indexer/             # 事件索引器
│   │   │   └── indexer.go       # 核心逻辑
│   │   └── calculator/          # 积分计算器
│   │       └── calculator.go    # 核心逻辑
│   └── api/
│       └── server.go            # HTTP API
├── pkg/
│   └── contract/
│       └── erc20/               # 合约 Go 绑定
├── configs/
│   └── config.toml              # 多链配置
├── go.mod
└── go.sum
```

## 2.2 核心服务详解

### 🔍 Indexer（事件索引器）

**职责**：
- 监听链上 Transfer 事件
- 维护用户余额快照
- 检测并处理区块重组

**工作流程**：
```
1. 从 block_cursor 读取上次同步位置
2. 批量拉取区块事件（chunk_size）
3. 解析 Transfer 事件 → 写入 balance_log
4. 更新 user_balance 快照
5. 记录 block_header（用于分叉检测）
6. 更新 block_cursor
```

**配置示例**：

```toml
# 以太坊链（简单确认机制）
[[chains]]
name = "sepolia"
chain_id = 11155111
type = "ethereum"
confirmations = 6        # 等待 6 个确认后入库
chunk_size = 10          # 每次拉取 10 个区块
request_delay_ms = 100   # 请求间隔 100ms

[[chains.contracts]]
address = "0xBEfe9d9726c3BFD513b6aDd74B243a82b272C073"
start_block = 10032808
token_decimals = 18

# OP Stack 链（Reorg Window 机制）
[[chains]]
name = "base-sepolia"
chain_id = 84532
type = "opstack"
reorg_window = 200       # 回溯 200 个区块检测分叉
chunk_size = 10
request_delay_ms = 200

[[chains.contracts]]
address = "0xB8a31EaC0874DC6f5a28FCa601336Ae32c723dF6"
start_block = 36257957
token_decimals = 18
```

**区块重组处理流程**：
```go
// 1. 定期检测（每 reorg_window 个区块）
if currentBlock % reorgWindow == 0 {
    EnsureCanonicalOrRollback(ctx, client, chainID, contract, reorgWindow)
}

// 2. 对比区块哈希
for block := safeBlock; block <= currentBlock; block++ {
    dbHash := getBlockHashFromDB(block)
    chainHash := getBlockHashFromChain(block)
    
    if dbHash != chainHash {
        // 发现分叉，回滚到安全区块
        RollbackTo(ctx, chainID, contract, safeBlock)
        break
    }
}

// 3. 回滚操作（数据库事务）
func RollbackTo(safeBlock) {
    tx := db.Begin()
    defer tx.Rollback()
    
    // 删除 > safeBlock 的所有数据
    tx.Exec("DELETE FROM balance_log WHERE block_number > ?", safeBlock)
    tx.Exec("DELETE FROM block_header WHERE block_number > ?", safeBlock)
    
    // 重建余额快照
    RebuildUserBalance(tx, safeBlock)
    
    // 更新游标
    tx.Exec("UPDATE block_cursor SET block_number = ?", safeBlock)
    
    tx.Commit()
}
```

### 🧮 Calculator（积分计算器）

**职责**：
- 基于余额变动计算积分
- 支持历史回溯（Backfill）
- 生成积分明细日志

**计算逻辑**：
```go
// 伪代码
func CalculatePoints(user) {
    // 使用行锁防止并发计算
    tx := db.Begin()
    tx.Exec("SELECT * FROM user_point WHERE account = ? FOR UPDATE", user)
    
    // 获取余额变动历史
    balanceChanges := getBalanceChanges(user, lastCalcTime)
    
    totalPoints := 0
    for each change in balanceChanges {
        duration = nextChangeTime - currentTime
        rate = getRateAt(currentTime)
        points = balance * duration * rate
        
        totalPoints += points
        
        // 保存明细（幂等性：唯一索引防重）
        savePointLog(balance, from_time, to_time, points, rate)
    }
    
    // 更新总积分
    updateUserPoint(user, totalPoints, currentTime)
    
    tx.Commit()
}
```

**积分计算示例**：
```
场景：用户余额变动历史
2026-01-13 12:00  余额 0   → 100  (Mint)
2026-01-14 10:00  余额 100 → 200  (收到转账)
2026-01-15 03:01  余额 200 → 200  (费率变更 5%→8%)
2026-01-16 08:00  余额 200 → 150  (转出)

计算过程：
时间段 1: 2026-01-13 12:00 ~ 2026-01-14 10:00
  余额=100, 时长=22小时, 费率=5%
  积分 = 100 × (22/24/365) × 0.05 = 0.251

时间段 2: 2026-01-14 10:00 ~ 2026-01-15 03:01
  余额=200, 时长=17.02小时, 费率=5%
  积分 = 200 × (17.02/24/365) × 0.05 = 0.194

时间段 3: 2026-01-15 03:01 ~ 2026-01-16 08:00
  余额=200, 时长=28.98小时, 费率=8%
  积分 = 200 × (28.98/24/365) × 0.08 = 0.532

总积分 = 0.251 + 0.194 + 0.532 + ... = 持续累积
```

### 🌐 API Server

**提供 RESTful API**：
```
GET /api/balance/:chain_id/:contract/:account          # 查询余额
GET /api/points/:chain_id/:contract/:account           # 查询积分
GET /api/balance/history/:chain_id/:contract/:account  # 余额历史
GET /api/points/history/:chain_id/:contract/:account   # 积分明细
```

## 2.3 并发安全机制

### 🔒 数据库层面

**1. 唯一索引（防止重复插入）**：
```sql
-- balance_log: 防止重复记录同一笔转账
UNIQUE KEY idx_balance_log (chain_id, contract_address, tx_hash, log_index, account)

-- user_point_log: 防止重复计算同一时间段
UNIQUE KEY idx_point_log (chain_id, contract_address, account, from_time, to_time)

-- block_header: 防止重复记录同一区块
UNIQUE KEY idx_block_header (chain_id, contract_address, block_number)
```

**2. 数据库事务（保证原子性）**：
```go
// 使用 GORM 事务
tx := db.Begin()
defer func() {
    if r := recover(); r != nil {
        tx.Rollback()
    }
}()

// 执行多个操作
tx.Create(&balanceLog)
tx.Save(&userBalance)
tx.Save(&blockCursor)

// 提交事务
if err := tx.Commit().Error; err != nil {
    return err
}
```

**3. 行级锁（防止并发计算）**：
```go
// 积分计算时使用 FOR UPDATE 锁定用户记录
var userPoint UserPoint
tx.Raw(`
    SELECT * FROM user_point 
    WHERE chain_id = ? AND account = ? 
    FOR UPDATE
`, chainID, account).Scan(&userPoint)

// 计算积分...

// 更新记录（其他事务等待）
tx.Save(&userPoint)
tx.Commit()
```

### 🔄 应用层面

**1. 幂等性设计**：
```go
// 所有写操作支持重复执行
db.Exec(`
    INSERT IGNORE INTO balance_log 
    (chain_id, contract_address, account, delta, tx_hash, log_index, ...)
    VALUES (?, ?, ?, ?, ?, ?, ...)
`)
// 如果唯一索引冲突，自动忽略（不报错）
```

**2. 多链并发隔离**：
```go
// 使用 errgroup 并发处理多条链
g, ctx := errgroup.WithContext(context.Background())

for _, chain := range chains {
    chain := chain // 避免闭包问题
    g.Go(func() error {
        return indexChain(ctx, chain) // 每条链独立处理
    })
}

if err := g.Wait(); err != nil {
    log.Fatal(err)
}
```

---

# 数据库设计

> ⭐ **完整数据示例请查看** [数据库示例.md](./数据库示例.md)
> 
> 该文档包含：
> - 所有表皆为真实数据
> - 详细的字段说明和业务逻辑
> - 数据一致性验证 SQL
> - 幂等性设计说明

## 核心表结构总览

| 表名 | 说明 | 记录数 | 唯一索引 |
|------|------|--------|----------|
| `balance_log` | 余额变动事实表 | 20 | (chain_id, contract, tx_hash, log_index, account) |
| `user_balance` | 用户余额快照 | 6 | (chain_id, contract, account) |
| `user_point` | 用户积分快照 | 6 | (chain_id, contract, account) |
| `user_point_log_1` | 积分明细（Sepolia） | 25 | (chain_id, contract, account, from_time, to_time) |
| `user_point_log_2` | 积分明细（Base Sepolia） | 13 | (chain_id, contract, account, from_time, to_time) |
| `block_cursor` | 同步进度 | 2 | (chain_id, contract) |
| `block_header` | 区块头信息 | 11 | (chain_id, contract, block_number) |
| `point_rate` | 积分费率配置 | 4 | (chain_id, contract, effective_time) |
| `sys_chains` | 链配置 | 2 | (chain_id) |
| `sys_contracts` | 合约配置 | 2 | (chain_id, address) |


> 以下设计基于生产规模模拟，用于说明在多链 Indexer 场景下，  
> 数据规模增长时，数据结构对性能、稳定性和故障隔离能力的影响。

## 为什么 user_point_log 要分表？

### 📊 分表策略：按链分表

**表命名规则**：
- `user_point_log_1` → Sepolia (chain_id=11155111)
- `user_point_log_2` → Base Sepolia (chain_id=84532)
- `user_point_log_N` → 其他链...

### 🎯 分表原因

#### 1. **数据量爆炸**
```
生产环境模拟：
- 1000个活跃用户
- 每天平均10次转账
- 1年 = 1000 × 10 × 365 = 365万条记录
- 10条链 = 3650万条记录

不分表问题：
- 单表 3650万行，索引 2.2GB
- 查询需要扫描全表并过滤 chain_id
- 插入性能随数据量增长而下降
```

#### 2. **查询性能优化**
```sql
-- ❌ 不分表：需要扫描全表并过滤 chain_id
SELECT * FROM user_point_log 
WHERE chain_id = 11155111 AND account = '0x...'
ORDER BY from_time DESC LIMIT 100;
-- 扫描：3650万行 → 过滤 → 返回100行

-- ✅ 分表：直接定位到目标表
SELECT * FROM user_point_log_1 
WHERE account = '0x...'
ORDER BY from_time DESC LIMIT 100;
-- 扫描：365万行 → 返回100行（性能提升10倍）
```

#### 3. **索引效率提升**
```
不分表索引大小：
- 主键索引：(id) → 3650万行 → ~700MB
- 唯一索引：(chain_id, contract, account, from_time, to_time) → ~1.5GB
- 总计：~2.2GB 索引

分表索引大小（单表）：
- 主键索引：(id) → 365万行 → ~70MB
- 唯一索引：(contract, account, from_time, to_time) → ~150MB
- 总计：~220MB 索引（减少90%）
```

#### 4. **维护和备份便利**
```bash
# 按链独立备份
mysqldump timeledger user_point_log_1 > sepolia_points.sql
mysqldump timeledger user_point_log_2 > base_points.sql

# 按链独立归档（删除旧数据）
DELETE FROM user_point_log_1 WHERE from_time < '2025-01-01';

# 按链独立优化
OPTIMIZE TABLE user_point_log_1;
```

#### 5. **隔离故障影响**
```
场景：Sepolia 链发生大规模 Reorg
- 需要删除并重算大量积分明细
- 如果不分表：影响所有链的查询性能
- 分表后：只影响 user_point_log_1，其他链不受影响
```

### 🔧 分表实现方式

```go
// 根据 chain_id 动态选择表名
func GetPointLogTable(chainID int64) string {
    tableMap := map[int64]string{
        11155111: "user_point_log_1",  // Sepolia
        84532:    "user_point_log_2",  // Base Sepolia
    }
    return tableMap[chainID]
}

// 插入积分明细
tableName := GetPointLogTable(chainID)
db.Exec(fmt.Sprintf(`
    INSERT INTO %s (chain_id, contract_address, account, ...)
    VALUES (?, ?, ?, ...)
`, tableName), ...)
```

### 📈 分表效果对比

| 指标 | 不分表 | 分表 | 提升 |
|------|--------|------|------|
| 单次查询扫描行数 | 3650万 | 365万 | 10倍 |
| 索引大小 | 2.2GB | 220MB | 10倍 |
| 插入性能 | 慢（索引大） | 快 | 3-5倍 |
| 备份时间 | 长（全表） | 短（按链） | 按需 |
| 故障隔离 | 全局影响 | 链级隔离 | ✅ |

---

## 数据一致性验证

### 余额一致性
```sql
-- 验证：balance_log 累加 = user_balance
SELECT 
    account,
    SUM(delta) as log_balance,
    (SELECT balance FROM user_balance ub 
     WHERE ub.account = bl.account AND ub.chain_id = bl.chain_id) as snapshot_balance
FROM balance_log bl
WHERE chain_id = 11155111
GROUP BY account;
```

### 积分一致性
```sql
-- 验证：user_point_log 累加 = user_point
SELECT 
    account,
    SUM(points) as log_points,
    (SELECT total_points FROM user_point up 
     WHERE up.account = upl.account AND up.chain_id = upl.chain_id) as snapshot_points
FROM user_point_log_1 upl
WHERE chain_id = 11155111
GROUP BY account;
```

---

# 快速开始

## 环境要求

- Go 1.24+
- MySQL 8.0+
- Redis 6.0+
- Foundry（合约部署）

## 启动服务

### 1. 配置环境变量
```bash
cd timeledger-backend
cp .env.example .env

# 编辑 .env
DB_USER=root
DB_PASSWORD=your_password
DB_HOST=localhost
DB_PORT=3306
DB_NAME=timeledger

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

SEPOLIA_RPC_URL=https://eth-sepolia.g.alchemy.com/v2/YOUR_KEY
BASE_SEPOLIA_RPC_URL=https://base-sepolia.g.alchemy.com/v2/YOUR_KEY
```

### 2. 初始化数据库

**数据库会自动创建表结构和初始化数据**

表结构和初始化逻辑在 `timeledger-backend/internal/repository/system_repo.go` 中：
- `InitSystem()` 函数会自动创建所有表
- 自动同步 `config.toml` 配置到数据库
- 自动创建动态分表 `user_point_log_1`, `user_point_log_2` 等
- 自动初始化默认积分费率（5%）

```bash
# 创建数据库
mysql -u root -p -e "CREATE DATABASE timeledger CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 程序启动时会自动建表和初始化数据
```

### 3. 启动服务

**方式一：启动所有服务（推荐）**
```bash
go run cmd/server/main.go all
```

**方式二：分别启动**
```bash
# 终端 1：启动 Indexer
go run cmd/server/main.go indexer

# 终端 2：启动 Calculator
go run cmd/server/main.go calculator

# 终端 3：启动 API Server
go run cmd/server/main.go api
```

---

# API 文档

## 查询用户余额
```http
GET /api/balance/:chain_id/:contract/:account

Response:
{
  "chain_id": 11155111,
  "contract_address": "0xBEfe9d9726c3BFD513b6aDd74B243a82b272C073",
  "account": "0x8a89a8FA663845284a645f95C5d87Ba1D1d25Dd1",
  "balance": "120000000000000000000",
  "balance_formatted": "120.0",
  "block_number": 10046278,
  "block_time": "2026-01-15T03:01:24Z"
}
```

## 查询用户积分
```http
GET /api/points/:chain_id/:contract/:account

Response:
{
  "chain_id": 11155111,
  "contract_address": "0xBEfe9d9726c3BFD513b6aDd74B243a82b272C073",
  "account": "0x8a89a8FA663845284a645f95C5d87Ba1D1d25Dd1",
  "total_points": "4068.635",
  "last_calc_time": "2026-01-26T08:26:36Z"
}
```

## 查询余额历史
```http
GET /api/balance/history/:chain_id/:contract/:account?limit=10&offset=0

Response:
{
  "total": 8,
  "items": [
    {
      "delta": "-70000000000000000000",
      "delta_formatted": "-70.0",
      "balance_after": "120000000000000000000",
      "balance_after_formatted": "120.0",
      "block_number": 10046278,
      "block_time": "2026-01-15T03:01:24Z",
      "tx_hash": "0x3e6e8a4a767fe80384b203b9494fbdf82e0867e568ad8f16d9385a930ca01e12"
    }
  ]
}
```

## 查询积分明细
```http
GET /api/points/history/:chain_id/:contract/:account?limit=10&offset=0

Response:
{
  "total": 25,
  "items": [
    {
      "balance": "120.0",
      "from_time": "2026-01-25T20:30:00Z",
      "to_time": "2026-01-26T08:26:36Z",
      "points": "114.656",
      "rate": "0.08"
    }
  ]
}
```

---

## 许可证

MIT License
