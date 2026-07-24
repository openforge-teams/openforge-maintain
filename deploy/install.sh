#!/bin/bash
#
# openforge-maintain 一键安装脚本
# 支持: Ubuntu/Debian/CentOS/AlmaLinux/Rocky Linux
#
set -euo pipefail

# ==================== 配置 ====================
INSTALL_DIR="/opt/openforge-maintain"
BIN_DIR="${INSTALL_DIR}/bin"
DATA_DIR="${INSTALL_DIR}/data"
FRONTEND_DIR="${INSTALL_DIR}/frontend"
LOG_DIR="/var/log/openforge-maintain"
RELEASE_URL="https://github.com/openforge-teams/openforge-maintain/releases/latest/download"
CORE_PORT=9999
AGENT_PORT=10000
VERSION="1.0.0"

# 颜色
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info()  { echo -e "${GREEN}[INFO]${NC} $1"; }
log_warn()  { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# ==================== 环境检测 ====================
detect_os() {
    log_info "检测操作系统..."

    if [ -f /etc/os-release ]; then
        . /etc/os-release
        OS_ID="${ID}"
        OS_VERSION="${VERSION_ID}"
        OS_NAME="${PRETTY_NAME}"
    elif [ -f /etc/redhat-release ]; then
        OS_ID="centos"
        OS_VERSION=$(cat /etc/redhat-release | grep -oE '[0-9]+\.[0-9]+' | head -1)
        OS_NAME=$(cat /etc/redhat-release)
    else
        log_error "不支持的操作系统"
        exit 1
    fi

    log_info "操作系统: ${OS_NAME} (${OS_ID} ${OS_VERSION})"

    ARCH=$(uname -m)
    case "${ARCH}" in
        x86_64)  ARCH_SUFFIX="amd64" ;;
        aarch64) ARCH_SUFFIX="arm64" ;;
        *) log_error "不支持的架构: ${ARCH}"; exit 1 ;;
    esac

    log_info "架构: ${ARCH} (${ARCH_SUFFIX})"
}

check_prerequisites() {
    log_info "检查前置条件..."

    # root 权限
    if [ "$(id -u)" -ne 0 ]; then
        log_error "请使用 root 用户运行安装脚本"
        exit 1
    fi

    # 端口检查
    for port in ${CORE_PORT} ${AGENT_PORT}; do
        if ss -tlnp | grep -q ":${port} "; then
            log_warn "端口 ${port} 已被占用，请确认后继续"
            read -p "是否继续? (y/N): " confirm
            if [ "${confirm}" != "y" ] && [ "${confirm}" != "Y" ]; then
                exit 1
            fi
        fi
    done

    log_info "前置条件检查通过"
}

# ==================== 安装依赖 ====================
install_deps() {
    log_info "安装系统依赖..."

    case "${OS_ID}" in
        ubuntu|debian)
            export DEBIAN_FRONTEND=noninteractive
            apt-get update -qq
            apt-get install -y -qq curl wget jq unzip tar socat sqlite3 > /dev/null
            ;;
        centos|almalinux|rocky)
            yum install -y -q curl wget jq unzip tar socat sqlite > /dev/null 2>&1 || \
            dnf install -y -q curl wget jq unzip tar socat sqlite > /dev/null 2>&1
            ;;
        kylin|uos)
            if command -v apt-get &> /dev/null; then
                export DEBIAN_FRONTEND=noninteractive
                apt-get update -qq
                apt-get install -y -qq curl wget jq unzip tar socat sqlite3 > /dev/null
            elif command -v yum &> /dev/null; then
                yum install -y -q curl wget jq unzip tar socat sqlite > /dev/null 2>&1
            fi
            ;;
        *)
            log_warn "未知发行版，跳过依赖安装，请手动安装: curl wget jq unzip tar socat sqlite"
            ;;
    esac

    log_info "系统依赖安装完成"
}

install_docker() {
    log_info "检测 Docker..."

    if command -v docker &> /dev/null; then
        DOCKER_VERSION=$(docker --version | grep -oE '[0-9]+\.[0-9]+')
        log_info "Docker 已安装 (${DOCKER_VERSION})"
    else
        log_info "安装 Docker CE..."

        if curl -fsSL https://get.docker.com | sh; then
            systemctl enable docker
            systemctl start docker
            log_info "Docker 安装完成"
        else
            log_error "Docker 安装失败，请手动安装"
            exit 1
        fi
    fi

    # 安装 docker-compose
    if ! command -v docker compose &> /dev/null && ! command -v docker-compose &> /dev/null; then
        log_info "安装 Docker Compose 插件..."
        DOCKER_COMPOSE_URL="https://github.com/docker/compose/releases/latest/download/docker-compose-linux-${ARCH_SUFFIX}"
        curl -fsSL "${DOCKER_COMPOSE_URL}" -o /usr/local/bin/docker-compose
        chmod +x /usr/local/bin/docker-compose
    fi
}

install_fail2ban() {
    log_info "检测 Fail2ban..."

    if command -v fail2ban-client &> /dev/null; then
        log_info "Fail2ban 已安装"
    else
        log_info "安装 Fail2ban..."

        case "${OS_ID}" in
            ubuntu|debian|kylin|uos)
                export DEBIAN_FRONTEND=noninteractive
                apt-get install -y -qq fail2ban > /dev/null
                ;;
            centos|almalinux|rocky)
                yum install -y -q fail2ban > /dev/null 2>&1 || \
                dnf install -y -q fail2ban > /dev/null 2>&1
                ;;
        esac

        # openforge-maintain自定义配置
        cat > /etc/fail2ban/jail.d/openforge-maintain.conf <<'EOF'
[openforge-maintain]
enabled = true
port = ${CORE_PORT},${AGENT_PORT}
filter = openforge-maintain
logpath = ${LOG_DIR}/access.log
maxretry = 5
bantime = 3600
findtime = 300
EOF

        cat > /etc/fail2ban/filter.d/openforge-maintain.conf <<'EOF'
[INCLUDES]
before = common.conf

[Definition]
failregex = ^.*"status":"fail".*"ip":"<HOST>".*$
ignoreregex =
EOF

        systemctl enable fail2ban 2>/dev/null || true
        log_info "Fail2ban 安装完成"
    fi
}

# ==================== 安装openforge-maintain ====================
create_directories() {
    log_info "创建安装目录..."
    mkdir -p "${BIN_DIR}" "${DATA_DIR}" "${FRONTEND_DIR}" "${LOG_DIR}"
    log_info "目录创建完成"
}

download_binaries() {
    log_info "下载 openforge-maintain 二进制文件..."

    ARCHIVE="openforge-maintain-linux-${ARCH_SUFFIX}.tar.gz"
    DOWNLOAD_URL="${RELEASE_URL}/${ARCHIVE}"

    # 如果本地已有编译好的文件，直接复制
    if [ -f "./bin/core" ] && [ -f "./bin/agent" ]; then
        log_info "使用本地编译的二进制文件"
        cp -f ./bin/core "${BIN_DIR}/core"
        cp -f ./bin/agent "${BIN_DIR}/agent"
    else
        log_info "从远程下载..."
        TMPDIR=$(mktemp -d)
        curl -fsSL "${DOWNLOAD_URL}" -o "${TMPDIR}/${ARCHIVE}"
        tar xzf "${TMPDIR}/${ARCHIVE}" -C "${BIN_DIR}"
        rm -rf "${TMPDIR}"
    fi

    chmod +x "${BIN_DIR}/core" "${BIN_DIR}/agent"
    log_info "二进制文件安装完成"
}

deploy_frontend() {
    log_info "部署前端资源..."

    # 如果本地有构建好的前端，复制过去
    if [ -d "./frontend/dist" ]; then
        cp -r ./frontend/dist/* "${FRONTEND_DIR}/"
        log_info "前端资源部署完成"
    else
        log_warn "未找到前端构建产物 (frontend/dist)，跳过前端部署"
        log_warn "请单独构建前端: cd frontend && npm install && npm run build"
    fi
}

generate_security_entry() {
    SECURITY_ENTRY=$(head -c 12 /dev/urandom | base64 | tr -dc 'a-z0-9')
    echo "${SECURITY_ENTRY}" > "${DATA_DIR}/security_entry"
    log_info "安全入口已生成: /${SECURITY_ENTRY}/"
}

generate_password() {
    # 生成 16 位随机密码
    ADMIN_PASSWORD=$(head -c 16 /dev/urandom | base64 | tr -dc 'A-Za-z0-9' | head -c 16)
    echo "${ADMIN_PASSWORD}" > "${DATA_DIR}/default_password"
}

setup_systemd() {
    log_info "配置 systemd 服务..."

    # 读取安全入口
    local security_entry="default"
    if [ -f "${DATA_DIR}/security_entry" ]; then
        security_entry=$(cat "${DATA_DIR}/security_entry")
    fi

    # 生成随机 SECRET
    MAINTAIN_SECRET=$(head -c 32 /dev/urandom | base64)

    # 复制 systemd unit 文件
    cp deploy/systemd/openforge-maintain-core.service /etc/systemd/system/
    cp deploy/systemd/openforge-maintain-agent.service /etc/systemd/system/

    # 注入实际配置
    sed -i "s|Environment=MAINTAIN_SECRET=CHANGE_ME_IN_PRODUCTION|Environment=MAINTAIN_SECRET=${MAINTAIN_SECRET}|g" /etc/systemd/system/openforge-maintain-core.service
    sed -i "s|Environment=MAINTAIN_SECURITY_ENTRY=default|Environment=MAINTAIN_SECURITY_ENTRY=${security_entry}|g" /etc/systemd/system/openforge-maintain-core.service

    # 重新加载 systemd
    systemctl daemon-reload

    # 启用服务
    systemctl enable openforge-maintain-core openforge-maintain-agent

    log_info "systemd 服务配置完成"
}

# ==================== 启动服务 ====================
start_services() {
    log_info "启动 openforge-maintain..."

    systemctl start openforge-maintain-core
    sleep 2
    systemctl start openforge-maintain-agent
    sleep 2

    # 检查服务状态
    if systemctl is-active --quiet openforge-maintain-core && systemctl is-active --quiet openforge-maintain-agent; then
        log_info "服务启动成功"
    else
        log_error "服务启动失败，请查看日志:"
        log_error "  journalctl -u openforge-maintain-core -n 50"
        log_error "  journalctl -u openforge-maintain-agent -n 50"
    fi
}

# ==================== 打印安装信息 ====================
print_success() {
    local security_entry=$(cat "${DATA_DIR}/security_entry" 2>/dev/null || echo "default")
    local admin_password=$(cat "${DATA_DIR}/default_password" 2>/dev/null || echo "请查看日志")
    local server_ip=$(hostname -I 2>/dev/null | awk '{print $1}' || echo "<服务器IP>")

    echo ""
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}   openforge-maintain 安装成功!${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo ""
    echo -e "  ${BLUE}openforge-maintain地址:${NC}  http://${server_ip}:${CORE_PORT}/${security_entry}/"
    echo -e "  ${BLUE}Core API:${NC}   http://${server_ip}:${CORE_PORT}/api/v2/core/"
    echo -e "  ${BLUE}Agent API:${NC}  http://${server_ip}:${AGENT_PORT}/api/v2/"
    echo ""
    echo -e "  ${YELLOW}初始账号:${NC}  admin"
    echo -e "  ${YELLOW}初始密码:${NC}  ${admin_password}"
    echo ""
    echo -e "  ${RED}请登录后立即修改默认密码!${NC}"
    echo ""
    echo "  常用命令:"
    echo "    启动:   systemctl start openforge-maintain-core openforge-maintain-agent"
    echo "    停止:   systemctl stop openforge-maintain-core openforge-maintain-agent"
    echo "    重启:   systemctl restart openforge-maintain-core openforge-maintain-agent"
    echo "    状态:   systemctl status openforge-maintain-core"
    echo "    日志:   journalctl -u openforge-maintain-core -f"
    echo "    卸载:   /opt/openforge-maintain/bin/core uninstall"
    echo ""
    echo -e "${GREEN}========================================${NC}"
}

# ==================== 卸载 ====================
uninstall() {
    log_warn "即将卸载 openforge-maintain..."

    # 停止服务
    systemctl stop openforge-maintain-core openforge-maintain-agent 2>/dev/null || true
    systemctl disable openforge-maintain-core openforge-maintain-agent 2>/dev/null || true

    # 删除文件
    rm -rf "${INSTALL_DIR}"
    rm -f /etc/systemd/system/openforge-maintain-core.service
    rm -f /etc/systemd/system/openforge-maintain-agent.service
    systemctl daemon-reload

    log_info "openforge-maintain 已卸载"
}

# ==================== 主流程 ====================
main() {
    echo ""
    echo -e "${BLUE}╔══════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║       openforge-maintain 安装向导 v${VERSION}       ║${NC}"
    echo -e "${BLUE}╚══════════════════════════════════════╝${NC}"
    echo ""

    if [ "${1:-}" = "uninstall" ]; then
        uninstall
        exit 0
    fi

    detect_os
    check_prerequisites
    install_deps
    install_docker
    install_fail2ban
    create_directories
    download_binaries
    deploy_frontend
    generate_security_entry
    generate_password
    setup_systemd
    start_services
    print_success
}

main "$@"
