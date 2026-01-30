#!/bin/bash

# RuleBack 项目初始化脚本
# 用于将模板项目重命名为新项目

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 打印带颜色的消息
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 显示使用说明
usage() {
    echo "用法: $0 <新项目名称>"
    echo ""
    echo "示例:"
    echo "  $0 myapp"
    echo "  $0 my-awesome-project"
    echo ""
    echo "这将把所有 'ruleback' 引用替换为新的项目名称"
    exit 1
}

# 检查参数
if [ -z "$1" ]; then
    print_error "请提供新项目名称"
    usage
fi

NEW_NAME=$1
OLD_NAME="ruleback"

# 验证项目名称格式
if [[ ! "$NEW_NAME" =~ ^[a-z][a-z0-9_-]*$ ]]; then
    print_error "项目名称只能包含小写字母、数字、下划线和连字符，且必须以字母开头"
    exit 1
fi

# 获取脚本所在目录的父目录（项目根目录）
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

print_info "项目目录: $PROJECT_DIR"
print_info "旧名称: $OLD_NAME"
print_info "新名称: $NEW_NAME"
echo ""

# 确认操作
read -p "确认要将项目重命名为 '$NEW_NAME' 吗? (y/N) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    print_warn "操作已取消"
    exit 0
fi

echo ""
print_info "开始重命名项目..."

# 进入项目目录
cd "$PROJECT_DIR"

# 1. 替换 go.mod 中的模块名
print_info "更新 go.mod..."
if [ -f "go.mod" ]; then
    sed -i '' "s/module $OLD_NAME/module $NEW_NAME/g" go.mod
fi

# 2. 替换所有 Go 文件中的 import 路径
print_info "更新 Go 文件中的 import 路径..."
find . -type f -name "*.go" -not -path "./vendor/*" | while read -r file; do
    sed -i '' "s|\"$OLD_NAME/|\"$NEW_NAME/|g" "$file"
done

# 3. 替换配置文件中的项目名
print_info "更新配置文件..."
if [ -f "configs/config.yaml" ]; then
    sed -i '' "s/name: \"$OLD_NAME\"/name: \"$NEW_NAME\"/g" configs/config.yaml
fi
if [ -f "configs/config.yaml.example" ]; then
    sed -i '' "s/name: \"$OLD_NAME\"/name: \"$NEW_NAME\"/g" configs/config.yaml.example
fi

# 4. 更新 README.md 中的项目名
print_info "更新 README.md..."
if [ -f "README.md" ]; then
    sed -i '' "s/# RuleBack/# ${NEW_NAME^}/g" README.md
    sed -i '' "s/$OLD_NAME/$NEW_NAME/g" README.md
fi

# 5. 更新 CLAUDE.md 中的项目名
print_info "更新 CLAUDE.md..."
if [ -f "CLAUDE.md" ]; then
    sed -i '' "s/$OLD_NAME/$NEW_NAME/g" CLAUDE.md
fi

# 6. 清理并重新下载依赖
print_info "清理 Go 模块缓存..."
rm -f go.sum
go mod tidy

# 7. 重新生成 Wire 代码
print_info "重新生成 Wire 代码..."
if command -v wire &> /dev/null; then
    wire ./internal/wire/...
elif [ -f "$HOME/go/bin/wire" ]; then
    "$HOME/go/bin/wire" ./internal/wire/...
else
    print_warn "Wire 未安装，请手动运行: go install github.com/google/wire/cmd/wire@latest"
    print_warn "然后运行: wire ./internal/wire/..."
fi

echo ""
print_info "项目初始化完成!"
echo ""
echo "后续步骤:"
echo "  1. 编辑 configs/config.yaml 配置数据库连接"
echo "  2. 运行项目: go run cmd/server/main.go cmd/server/bootstrap.go"
echo ""
echo "配合 AI 使用:"
echo "  1. 让 AI 阅读 CLAUDE.md 了解项目规范"
echo "  2. 描述你需要的功能，AI 会自动生成符合规范的代码"
echo "  3. 详细指南请参考 docs/AI_USAGE.md"
