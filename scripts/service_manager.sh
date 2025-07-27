#!/bin/bash

# 服务管理脚本
# 用于管理 joyshop_srvs 项目中的所有微服务

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 项目根目录
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# 服务配置
declare -A SERVICES=(
    ["goods_srv"]="50051"
    ["inventory_srv"]="50052"
    ["order_srv"]="50053"
    ["user_srv"]="50054"
)

# 日志目录
LOG_DIR="$PROJECT_ROOT/logs"
PID_DIR="$PROJECT_ROOT/pids"

# 创建必要的目录
mkdir -p "$LOG_DIR" "$PID_DIR"

# 打印帮助信息
show_help() {
    echo -e "${BLUE}Joyshop 微服务管理脚本${NC}"
    echo ""
    echo "用法: $0 [命令] [服务名]"
    echo ""
    echo "命令:"
    echo "  start [服务名]    启动指定服务（不指定则启动所有服务）"
    echo "  stop [服务名]     停止指定服务（不指定则停止所有服务）"
    echo "  restart [服务名]  重启指定服务（不指定则重启所有服务）"
    echo "  status [服务名]   查看指定服务状态（不指定则查看所有服务）"
    echo "  logs [服务名]     查看指定服务的日志"
    echo "  clean             清理日志和PID文件"
    echo "  help              显示此帮助信息"
    echo ""
    echo "服务名:"
    for service in "${!SERVICES[@]}"; do
        echo "  $service (端口: ${SERVICES[$service]})"
    done
    echo ""
    echo "示例:"
    echo "  $0 start                    # 启动所有服务"
    echo "  $0 start goods_srv          # 启动商品服务"
    echo "  $0 stop user_srv            # 停止用户服务"
    echo "  $0 status                   # 查看所有服务状态"
    echo "  $0 logs order_srv           # 查看订单服务日志"
}

# 检查服务是否运行
is_service_running() {
    local service_name=$1
    local pid_file="$PID_DIR/${service_name}.pid"
    
    if [[ -f "$pid_file" ]]; then
        local pid=$(cat "$pid_file")
        if kill -0 "$pid" 2>/dev/null; then
            return 0
        else
            rm -f "$pid_file"
        fi
    fi
    return 1
}

# 获取服务状态
get_service_status() {
    local service_name=$1
    local port=${SERVICES[$service_name]}
    
    if is_service_running "$service_name"; then
        local pid=$(cat "$PID_DIR/${service_name}.pid")
        echo -e "${GREEN}✓${NC} $service_name (PID: $pid, 端口: $port) - 运行中"
    else
        echo -e "${RED}✗${NC} $service_name (端口: $port) - 已停止"
    fi
}

# 启动单个服务
start_service() {
    local service_name=$1
    local port=${SERVICES[$service_name]}
    local service_dir="$PROJECT_ROOT/${service_name}"
    local log_file="$LOG_DIR/${service_name}.log"
    local pid_file="$PID_DIR/${service_name}.pid"
    
    echo -e "${BLUE}正在启动 $service_name...${NC}"
    
    if is_service_running "$service_name"; then
        echo -e "${YELLOW}警告: $service_name 已经在运行中${NC}"
        return 0
    fi
    
    if [[ ! -d "$service_dir" ]]; then
        echo -e "${RED}错误: 服务目录不存在: $service_dir${NC}"
        return 1
    fi
    
    # 检查端口是否被占用
    if lsof -i ":$port" >/dev/null 2>&1; then
        echo -e "${RED}错误: 端口 $port 已被占用${NC}"
        return 1
    fi
    
    # 启动服务
    cd "$service_dir"
    nohup go run main.go > "$log_file" 2>&1 &
    local pid=$!
    
    # 保存PID
    echo $pid > "$pid_file"
    
    # 等待服务启动
    sleep 2
    
    if is_service_running "$service_name"; then
        echo -e "${GREEN}✓ $service_name 启动成功 (PID: $pid)${NC}"
    else
        echo -e "${RED}✗ $service_name 启动失败${NC}"
        rm -f "$pid_file"
        return 1
    fi
}

# 停止单个服务
stop_service() {
    local service_name=$1
    local pid_file="$PID_DIR/${service_name}.pid"
    
    echo -e "${BLUE}正在停止 $service_name...${NC}"
    
    if ! is_service_running "$service_name"; then
        echo -e "${YELLOW}警告: $service_name 未在运行${NC}"
        return 0
    fi
    
    local pid=$(cat "$pid_file")
    
    # 发送SIGTERM信号
    kill -TERM "$pid" 2>/dev/null || true
    
    # 等待进程结束
    local count=0
    while kill -0 "$pid" 2>/dev/null && [[ $count -lt 10 ]]; do
        sleep 1
        ((count++))
    done
    
    # 如果进程仍在运行，强制杀死
    if kill -0 "$pid" 2>/dev/null; then
        echo -e "${YELLOW}强制停止 $service_name...${NC}"
        kill -KILL "$pid" 2>/dev/null || true
    fi
    
    rm -f "$pid_file"
    echo -e "${GREEN}✓ $service_name 已停止${NC}"
}

# 重启单个服务
restart_service() {
    local service_name=$1
    echo -e "${BLUE}正在重启 $service_name...${NC}"
    stop_service "$service_name"
    sleep 2
    start_service "$service_name"
}

# 查看服务日志
show_logs() {
    local service_name=$1
    local log_file="$LOG_DIR/${service_name}.log"
    
    if [[ ! -f "$log_file" ]]; then
        echo -e "${RED}错误: 日志文件不存在: $log_file${NC}"
        return 1
    fi
    
    echo -e "${BLUE}显示 $service_name 的日志 (按 Ctrl+C 退出):${NC}"
    tail -f "$log_file"
}

# 清理日志和PID文件
clean_files() {
    echo -e "${BLUE}正在清理日志和PID文件...${NC}"
    rm -rf "$LOG_DIR"/*.log
    rm -rf "$PID_DIR"/*.pid
    echo -e "${GREEN}✓ 清理完成${NC}"
}

# 主函数
main() {
    local command=$1
    local service_name=$2
    
    case "$command" in
        "start")
            if [[ -n "$service_name" ]]; then
                if [[ -n "${SERVICES[$service_name]}" ]]; then
                    start_service "$service_name"
                else
                    echo -e "${RED}错误: 未知的服务名: $service_name${NC}"
                    exit 1
                fi
            else
                echo -e "${BLUE}正在启动所有服务...${NC}"
                for service in "${!SERVICES[@]}"; do
                    start_service "$service"
                done
                echo -e "${GREEN}✓ 所有服务启动完成${NC}"
            fi
            ;;
        "stop")
            if [[ -n "$service_name" ]]; then
                if [[ -n "${SERVICES[$service_name]}" ]]; then
                    stop_service "$service_name"
                else
                    echo -e "${RED}错误: 未知的服务名: $service_name${NC}"
                    exit 1
                fi
            else
                echo -e "${BLUE}正在停止所有服务...${NC}"
                for service in "${!SERVICES[@]}"; do
                    stop_service "$service"
                done
                echo -e "${GREEN}✓ 所有服务停止完成${NC}"
            fi
            ;;
        "restart")
            if [[ -n "$service_name" ]]; then
                if [[ -n "${SERVICES[$service_name]}" ]]; then
                    restart_service "$service_name"
                else
                    echo -e "${RED}错误: 未知的服务名: $service_name${NC}"
                    exit 1
                fi
            else
                echo -e "${BLUE}正在重启所有服务...${NC}"
                for service in "${!SERVICES[@]}"; do
                    restart_service "$service"
                done
                echo -e "${GREEN}✓ 所有服务重启完成${NC}"
            fi
            ;;
        "status")
            if [[ -n "$service_name" ]]; then
                if [[ -n "${SERVICES[$service_name]}" ]]; then
                    get_service_status "$service_name"
                else
                    echo -e "${RED}错误: 未知的服务名: $service_name${NC}"
                    exit 1
                fi
            else
                echo -e "${BLUE}服务状态:${NC}"
                for service in "${!SERVICES[@]}"; do
                    get_service_status "$service"
                done
            fi
            ;;
        "logs")
            if [[ -n "$service_name" ]]; then
                if [[ -n "${SERVICES[$service_name]}" ]]; then
                    show_logs "$service_name"
                else
                    echo -e "${RED}错误: 未知的服务名: $service_name${NC}"
                    exit 1
                fi
            else
                echo -e "${RED}错误: 请指定服务名${NC}"
                exit 1
            fi
            ;;
        "clean")
            clean_files
            ;;
        "help"|"-h"|"--help")
            show_help
            ;;
        *)
            echo -e "${RED}错误: 未知命令: $command${NC}"
            echo ""
            show_help
            exit 1
            ;;
    esac
}

# 检查参数
if [[ $# -eq 0 ]]; then
    show_help
    exit 1
fi

# 执行主函数
main "$@" 