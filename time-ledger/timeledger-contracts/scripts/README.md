# 生成 Go 绑定

这个脚本用于从 Solidity 合约生成 Go 语言绑定。

## 前置要求

1. **abigen**: Go Ethereum 工具，用于生成绑定
   ```bash
   # 安装 abigen
   go install github.com/ethereum/go-ethereum/cmd/abigen@latest
   ```

2. **jq** (可选): 用于提取 JSON 数据
   ```bash
   # Ubuntu/Debian
   sudo apt-get install jq
   
   # macOS
   brew install jq
   ```
   如果没有 jq，脚本会自动使用 Python 来提取数据。

3. **forge**: Foundry 工具，用于编译合约
   ```bash
   # 如果合约未编译，脚本会自动运行 forge build
   ```

## 使用方法

### 方法 1: 使用脚本（推荐）

```bash
cd timeledger-contracts
./scripts/generate-go-bindings.sh
```

### 方法 2: 手动生成

如果你想要手动生成绑定，可以按照以下步骤：

1. **编译合约**（如果还没有编译）:
   ```bash
   cd timeledger-contracts
   forge build
   ```

2. **提取 ABI 和 bytecode**:
   ```bash
   # 使用 jq
   jq -r '.abi' out/TimeLedgerToken.sol/TimeLedgerToken.json > abi/TimeLedgerToken.abi
   jq -r '.bytecode.object' out/TimeLedgerToken.sol/TimeLedgerToken.json > abi/TimeLedgerToken.bin
   ```

3. **生成 Go 绑定**:
   ```bash
   abigen \
     --abi abi/TimeLedgerToken.abi \
     --bin abi/TimeLedgerToken.bin \
     --pkg abi \
     --type TimeLedgerToken \
     --out abi/timeledgertoken.go
   ```

## 输出

生成的 Go 绑定文件位于：
```
timeledger-contracts/abi/timeledgertoken.go
```

## 在 Go 代码中使用

```go
package main

import (
    "context"
    "math/big"
    
    "github.com/ethereum/go-ethereum/accounts/abi/bind"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/ethclient"
    
    "github.com/Atom257/web3-labs/timeledger-contracts/abi"
)

func main() {
    // 连接到以太坊节点
    client, err := ethclient.Dial("https://your-rpc-url")
    if err != nil {
        panic(err)
    }
    
    // 合约地址
    contractAddress := common.HexToAddress("0x...")
    
    // 创建合约实例
    token, err := abi.NewTimeLedgerToken(contractAddress, client)
    if err != nil {
        panic(err)
    }
    
    // 调用只读方法
    name, err := token.Name(nil)
    if err != nil {
        panic(err)
    }
    println("Token name:", name)
    
    // 调用需要交易的方法
    auth, _ := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
    tx, err := token.Mint(auth, toAddress, amount)
    if err != nil {
        panic(err)
    }
    println("Transaction hash:", tx.Hash().Hex())
}
```

## 注意事项

- 生成的绑定文件是自动生成的，不要手动编辑
- 每次合约更新后，需要重新运行脚本生成新的绑定
- 确保 Go 模块已安装必要的依赖（`github.com/ethereum/go-ethereum`）
