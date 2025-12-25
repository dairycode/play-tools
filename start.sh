#!/bin/bash

# 颜色定义
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 获取脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo -e "${BLUE}================================${NC}"
echo -e "${BLUE}启动游戏助手服务${NC}"
echo -e "${BLUE}================================${NC}\n"

# 检查后端进程是否已在运行
if lsof -i :8080 > /dev/null 2>&1; then
    echo -e "${YELLOW}警告: 后端服务已在运行 (端口 8080)${NC}"
else
    # 启动后端服务
    echo -e "${GREEN}[1/2] 启动后端服务...${NC}"
    cd "$SCRIPT_DIR/server"
    nohup go run main.go > ../logs/backend.log 2>&1 &
    BACKEND_PID=$!
    echo $BACKEND_PID > ../logs/backend.pid
    echo -e "${GREEN}✓ 后端服务已启动 (PID: $BACKEND_PID, 端口: 8080)${NC}"
    echo -e "${GREEN}  日志文件: logs/backend.log${NC}\n"
fi

# 检查前端进程是否已在运行
if lsof -i :80 > /dev/null 2>&1; then
    echo -e "${YELLOW}警告: 前端服务已在运行 (端口 80)${NC}"
else
    # 启动前端服务
    echo -e "${GREEN}[2/2] 启动前端服务...${NC}"
    cd "$SCRIPT_DIR/web"
    nohup npm run dev:h5 > ../logs/frontend.log 2>&1 &
    FRONTEND_PID=$!
    echo $FRONTEND_PID > ../logs/frontend.pid
    echo -e "${GREEN}✓ 前端服务已启动 (PID: $FRONTEND_PID, 端口: 80)${NC}"
    echo -e "${GREEN}  日志文件: logs/frontend.log${NC}\n"
fi

sleep 2

echo -e "${BLUE}================================${NC}"
echo -e "${GREEN}所有服务启动完成!${NC}"
echo -e "${BLUE}================================${NC}"
echo -e "前端地址: ${GREEN}http://localhost:80${NC}"
echo -e "后端地址: ${GREEN}http://localhost:8080${NC}"
echo -e "\n查看日志:"
echo -e "  后端: tail -f logs/backend.log"
echo -e "  前端: tail -f logs/frontend.log"
echo -e "\n停止服务: ${YELLOW}./stop.sh${NC}\n"
