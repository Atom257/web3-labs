#!/bin/bash

# 生成 Go 绑定的脚本
# 使用方法: ./scripts/generate-go-bindings.sh

set -e

# 颜色输出
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 项目根目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CONTRACTS_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

# 合约信息
CONTRACT_NAME="TimeLedgerToken"
CONTRACT_JSON="$CONTRACTS_DIR/out/$CONTRACT_NAME.sol/$CONTRACT_NAME.json"
OUTPUT_DIR="$CONTRACTS_DIR/abi"
OUTPUT_FILE="$OUTPUT_DIR/${CONTRACT_NAME,,}.go"  # 转换为小写

printf "${GREEN}开始生成 Go 绑定...${NC}\n"

# 检查合约 JSON 文件是否存在
if [ ! -f "$CONTRACT_JSON" ]; then
    printf "${YELLOW}警告: 合约 JSON 文件不存在，正在编译合约...${NC}\n"
    cd "$CONTRACTS_DIR"
    forge build
fi

# 创建输出目录
mkdir -p "$OUTPUT_DIR"

# 提取 ABI 和 bytecode
printf "${GREEN}提取 ABI 和 bytecode...${NC}\n"
ABI_FILE="$OUTPUT_DIR/${CONTRACT_NAME}.abi"
BYTECODE_FILE="$OUTPUT_DIR/${CONTRACT_NAME}.bin"

# 使用 jq 提取 ABI（如果没有 jq，使用 Python）
if command -v jq &> /dev/null; then
    jq -r '.abi' "$CONTRACT_JSON" > "$ABI_FILE"
    jq -r '.bytecode.object' "$CONTRACT_JSON" > "$BYTECODE_FILE"
else
    printf "${YELLOW}jq 未安装，使用 Python 提取...${NC}\n"
    python3 << EOF
import json
import sys

with open('$CONTRACT_JSON', 'r') as f:
    data = json.load(f)

with open('$ABI_FILE', 'w') as f:
    json.dump(data['abi'], f, indent=2)

with open('$BYTECODE_FILE', 'w') as f:
    f.write(data['bytecode']['object'])
EOF
fi

# 生成 Go 绑定
printf "${GREEN}使用 abigen 生成 Go 绑定...${NC}\n"
abigen \
    --abi "$ABI_FILE" \
    --bin "$BYTECODE_FILE" \
    --pkg abi \
    --type "$CONTRACT_NAME" \
    --out "$OUTPUT_FILE"

# 清理临时文件
rm -f "$ABI_FILE" "$BYTECODE_FILE"

printf "${GREEN}✓ Go 绑定已生成: $OUTPUT_FILE${NC}\n"
