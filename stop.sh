#!/bin/bash

# 颜色定义
RED='\033[0;31m'
BLUE='\033[0;34m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

# 获取脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo -e "${BLUE}================================${NC}"
echo -e "${BLUE}停止游戏助手服务${NC}"
echo -e "${BLUE}================================${NC}\n"

# 停止后端服务
echo -e "${RED}[1/2] 停止后端服务...${NC}"
if [ -f "$SCRIPT_DIR/logs/backend.pid" ]; then
    BACKEND_PID=$(cat "$SCRIPT_DIR/logs/backend.pid")
    if ps -p $BACKEND_PID > /dev/null 2>&1; then
        kill $BACKEND_PID
        echo -e "${GREEN}✓ 后端服务已停止 (PID: $BACKEND_PID)${NC}"
    else
        echo -e "${GREEN}  后端服务未运行${NC}"
    fi
    rm "$SCRIPT_DIR/logs/backend.pid"
else
    # 尝试通过端口查找并关闭
    BACKEND_PID=$(lsof -ti :8081)
    if [ -n "$BACKEND_PID" ]; then
        kill $BACKEND_PID
        echo -e "${GREEN}✓ 后端服务已停止 (PID: $BACKEND_PID)${NC}"
    else
        echo -e "${GREEN}  后端服务未运行${NC}"
    fi
fi

# 停止前端服务
echo -e "${RED}[2/2] 停止前端服务...${NC}"
if [ -f "$SCRIPT_DIR/logs/frontend.pid" ]; then
    FRONTEND_PID=$(cat "$SCRIPT_DIR/logs/frontend.pid")
    if ps -p $FRONTEND_PID > /dev/null 2>&1; then
        kill $FRONTEND_PID
        echo -e "${GREEN}✓ 前端服务已停止 (PID: $FRONTEND_PID)${NC}"
    else
        echo -e "${GREEN}  前端服务未运行${NC}"
    fi
    rm "$SCRIPT_DIR/logs/frontend.pid"
else
    # 尝试通过端口查找并关闭
    FRONTEND_PID=$(lsof -ti :5173)
    if [ -n "$FRONTEND_PID" ]; then
        kill $FRONTEND_PID
        echo -e "${GREEN}✓ 前端服务已停止 (PID: $FRONTEND_PID)${NC}"
    else
        echo -e "${GREEN}  前端服务未运行${NC}"
    fi
fi

# 清理可能残留的 node 和 go run 进程
echo -e "\n${RED}清理残留进程...${NC}"
pkill -f "npm run dev:h5" 2>/dev/null && echo -e "${GREEN}✓ 清理 npm 进程${NC}"
pkill -f "go run main.go" 2>/dev/null && echo -e "${GREEN}✓ 清理 go run 进程${NC}"

sleep 1

echo -e "\n${BLUE}================================${NC}"
echo -e "${GREEN}所有服务已停止!${NC}"
echo -e "${BLUE}================================${NC}\n"
