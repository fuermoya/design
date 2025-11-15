# 毕业设计部署文档

### 说明

> 1. 默认账号密码  
    1.1 账号：admin  
    2.2 密码：123456  

## 部署概述

本文档详细说明了如何将基于 Gin + Vue3 的全栈项目部署到生产环境。项目采用前后端分离架构，支持多种部署方式。

## 部署架构

```
生产环境架构:
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Nginx (80)    │    │   Backend       │    │   MySQL         │
│   (负载均衡)     │────│   (8888)        │────│   (3306)        │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │
         │                       │
         ▼                       ▼
┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   Redis         │
│   (静态文件)     │    │   (缓存)        │
└─────────────────┘    └─────────────────┘
```

## 环境要求

### 服务器配置
- **CPU**: 2核心以上
- **内存**: 4GB以上
- **存储**: 50GB以上
- **操作系统**: Ubuntu 20.04+ / CentOS 7+ / Windows Server 2019+

### 软件要求
- **Go**: 1.23+
- **Node.js**: 16+
- **MySQL**: 8.0+
- **Nginx**: 1.18+
- **Redis**: 6.0+ (可选)

## 部署方式

### 方式一: 传统部署

#### 1. 服务器环境准备

##### Ubuntu/Debian 系统
```bash
# 更新系统
sudo apt update && sudo apt upgrade -y

# 安装基础工具
sudo apt install -y curl wget git unzip

# 安装Go
wget https://go.dev/dl/go1.23.6.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.23.6.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 安装Node.js
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs

# 安装MySQL
sudo apt install -y mysql-server
sudo systemctl start mysql
sudo systemctl enable mysql

# 安装Nginx
sudo apt install -y nginx
sudo systemctl start nginx
sudo systemctl enable nginx
```

##### CentOS/RHEL 系统
```bash
# 更新系统
sudo yum update -y

# 安装基础工具
sudo yum install -y curl wget git unzip

# 安装Go
wget https://go.dev/dl/go1.23.6.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.23.6.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 安装Node.js
curl -fsSL https://rpm.nodesource.com/setup_18.x | sudo bash -
sudo yum install -y nodejs

# 安装MySQL
sudo yum install -y mysql-server
sudo systemctl start mysqld
sudo systemctl enable mysqld

# 安装Nginx
sudo yum install -y nginx
sudo systemctl start nginx
sudo systemctl enable nginx
```

#### 2. 数据库配置

```bash
# 登录MySQL
sudo mysql -u root -p

# 创建数据库和用户
CREATE DATABASE design_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'design_user'@'localhost' IDENTIFIED BY 'your_secure_password';
GRANT ALL PRIVILEGES ON design_db.* TO 'design_user'@'localhost';
FLUSH PRIVILEGES;
EXIT;
```

#### 3. 后端部署

```bash
# 克隆项目
git clone <项目地址>
cd project/backend

# 安装依赖
go mod tidy
go mod download

# 编译项目
go build -o design-server main.go

# 创建配置文件
cat > config.yaml << EOF
system:
  env: "production"
  addr: 8888
  db-type: "mysql"
  oss-type: "local"

jwt:
  signing-key: "your_jwt_secret_key"
  expire: 7200
  buffer-time: 86400
  issuer: "qmPlus"

zap:
  level: "info"
  format: "json"
  prefix: "[DESIGN]"
  director: "log"
  link-name: "latest_log"
  show-line: true
  encode-level: "LowercaseColorLevelEncoder"
  stacktrace-key: "stacktrace"
  log-in-console: true

mysql:
  path: "127.0.0.1:3306"
  port: "3306"
  config: "charset=utf8mb4&parseTime=True&loc=Local"
  db-name: "design_db"
  username: "design_user"
  password: "your_secure_password"
  prefix: ""
  singular-table: false
  engine: ""
  max-idle-conns: 10
  max-open-conns: 100
  log-mode: ""
  log-zap: false

local:
  path: "uploads/file"
  store-path: "uploads/file"

cors:
  mode: "whitelist"
  whitelist:
    - "http://localhost:3000"
    - "http://your-domain.com"
EOF

# 创建服务目录
sudo mkdir -p /opt/design
sudo cp design-server /opt/design/
sudo cp config.yaml /opt/design/
sudo mkdir -p /opt/design/uploads

# 创建systemd服务
sudo tee /etc/systemd/system/design-server.service << EOF
[Unit]
Description=Design Backend Server
After=network.target mysql.service

[Service]
Type=simple
User=www-data
Group=www-data
WorkingDirectory=/opt/design
ExecStart=/opt/design/design-server
Restart=always
RestartSec=5
Environment=GIN_MODE=release

[Install]
WantedBy=multi-user.target
EOF

# 启动服务
sudo systemctl daemon-reload
sudo systemctl enable design-server
sudo systemctl start design-server
sudo systemctl status design-server
```

#### 4. 前端部署

```bash
# 进入前端目录
cd ../frontend

# 安装依赖
npm install

# 修改API地址
sed -i 's|http://localhost:8888|http://your-domain.com|g' src/utils/request.ts

# 构建项目
npm run build

# 部署到Nginx
sudo cp -r dist/* /var/www/html/
sudo chown -R www-data:www-data /var/www/html/
```

#### 5. Nginx配置

```bash
# 创建Nginx配置
sudo tee /etc/nginx/sites-available/design << EOF
server {
    listen 80;
    server_name your-domain.com;
    
    # 前端静态文件
    location / {
        root /var/www/html;
        try_files \$uri \$uri/ /index.html;
        index index.html;
    }
    
    # 后端API代理
    location /api/ {
        proxy_pass http://127.0.0.1:8888;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }
    
    # 文件上传
    location /uploads/ {
        alias /opt/design/uploads/;
        expires 30d;
        add_header Cache-Control "public, immutable";
    }
    
    # 安全配置
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header Referrer-Policy "no-referrer-when-downgrade" always;
    add_header Content-Security-Policy "default-src 'self' http: https: data: blob: 'unsafe-inline'" always;
}
EOF

# 启用配置
sudo ln -s /etc/nginx/sites-available/design /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### 方式二: Docker部署

#### 1. 安装Docker

```bash
# 安装Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# 安装Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/download/v2.20.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

#### 2. 创建Dockerfile

**后端Dockerfile**
```dockerfile
# backend/Dockerfile
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/config.yaml .
RUN mkdir -p uploads

EXPOSE 8888
CMD ["./main"]
```

**前端Dockerfile**
```dockerfile
# frontend/Dockerfile
FROM node:18-alpine AS builder

WORKDIR /app
COPY package*.json ./
RUN npm ci

COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf

EXPOSE 80
```

#### 3. 创建docker-compose.yml

```yaml
# docker-compose.yml
version: '3.8'

services:
  mysql:
    image: mysql:8.0
    container_name: design-mysql
    restart: always
    environment:
      MYSQL_ROOT_PASSWORD: root_password
      MYSQL_DATABASE: design_db
      MYSQL_USER: design_user
      MYSQL_PASSWORD: design_password
    volumes:
      - mysql_data:/var/lib/mysql
      - ./init.sql:/docker-entrypoint-initdb.d/init.sql
    ports:
      - "3306:3306"
    networks:
      - design-network

  redis:
    image: redis:7-alpine
    container_name: design-redis
    restart: always
    volumes:
      - redis_data:/data
    networks:
      - design-network

  backend:
    build: ./backend
    container_name: design-backend
    restart: always
    environment:
      - GIN_MODE=release
    volumes:
      - ./uploads:/root/uploads
    ports:
      - "8888:8888"
    depends_on:
      - mysql
      - redis
    networks:
      - design-network

  frontend:
    build: ./frontend
    container_name: design-frontend
    restart: always
    ports:
      - "80:80"
    depends_on:
      - backend
    networks:
      - design-network

  nginx:
    image: nginx:alpine
    container_name: design-nginx
    restart: always
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./uploads:/var/www/uploads
    ports:
      - "80:80"
      - "443:443"
    depends_on:
      - frontend
      - backend
    networks:
      - design-network

volumes:
  mysql_data:
  redis_data:

networks:
  design-network:
    driver: bridge
```

#### 4. 部署命令

```bash
# 构建并启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

### 方式三: 云平台部署

#### 1. 阿里云ECS部署

```bash
# 使用阿里云ECS，按照传统部署方式进行
# 配置安全组，开放80、443、8888端口
```

#### 2. 腾讯云CVM部署

```bash
# 使用腾讯云CVM，按照传统部署方式进行
# 配置安全组，开放相应端口
```

#### 3. AWS EC2部署

```bash
# 使用AWS EC2，按照传统部署方式进行
# 配置安全组，开放相应端口
```



## 监控和维护

### 1. 日志管理

```bash
# 查看后端日志
sudo journalctl -u design-server -f

# 查看Nginx日志
sudo tail -f /var/log/nginx/access.log
sudo tail -f /var/log/nginx/error.log

# 查看Docker日志
docker-compose logs -f backend
```

### 2. 性能监控

```bash
# 安装监控工具
sudo apt install -y htop iotop nethogs

# 监控系统资源
htop
iotop
nethogs

# 监控数据库
mysql -u root -p -e "SHOW PROCESSLIST;"
```

### 3. 备份策略

```bash
# 数据库备份脚本
#!/bin/bash
BACKUP_DIR="/backup/mysql"
DATE=$(date +%Y%m%d_%H%M%S)
mysqldump -u design_user -p design_db > $BACKUP_DIR/design_db_$DATE.sql

# 文件备份脚本
#!/bin/bash
BACKUP_DIR="/backup/uploads"
DATE=$(date +%Y%m%d_%H%M%S)
tar -czf $BACKUP_DIR/uploads_$DATE.tar.gz /opt/design/uploads/
```

### 4. 自动更新

```bash
# 创建更新脚本
#!/bin/bash
cd /opt/design
git pull origin main
go build -o design-server main.go
sudo systemctl restart design-server
```

## 安全配置

### 1. 防火墙配置

```bash
# Ubuntu/Debian
sudo ufw allow 22/tcp
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable

# CentOS/RHEL
sudo firewall-cmd --permanent --add-service=ssh
sudo firewall-cmd --permanent --add-service=http
sudo firewall-cmd --permanent --add-service=https
sudo firewall-cmd --reload
```

### 2. SSL证书配置

```bash
# 安装Certbot
sudo apt install certbot python3-certbot-nginx

# 获取SSL证书
sudo certbot --nginx -d your-domain.com

# 自动续期
sudo crontab -e
# 添加: 0 12 * * * /usr/bin/certbot renew --quiet
```

### 3. 安全加固

```bash
# 禁用root登录
sudo sed -i 's/#PermitRootLogin yes/PermitRootLogin no/' /etc/ssh/sshd_config

# 更改SSH端口
sudo sed -i 's/#Port 22/Port 2222/' /etc/ssh/sshd_config

# 重启SSH服务
sudo systemctl restart sshd
```

## 故障排除

### 1. 常见问题

**问题1: 数据库连接失败**
```bash
# 检查MySQL服务状态
sudo systemctl status mysql

# 检查数据库连接
mysql -u design_user -p design_db

# 检查配置文件
cat /opt/design/config.yaml
```

**问题2: 前端无法访问后端API**
```bash
# 检查后端服务状态
sudo systemctl status design-server

# 检查端口是否开放
netstat -tlnp | grep 8888

# 检查Nginx配置
sudo nginx -t
```

**问题3: 文件上传失败**
```bash
# 检查上传目录权限
ls -la /opt/design/uploads/

# 修改权限
sudo chown -R www-data:www-data /opt/design/uploads/
sudo chmod -R 755 /opt/design/uploads/
```

### 2. 性能优化

```bash
# 数据库优化
mysql -u root -p
SHOW VARIABLES LIKE 'max_connections';
SET GLOBAL max_connections = 200;

# Nginx优化
sudo nano /etc/nginx/nginx.conf
# 调整worker_processes和worker_connections

# 系统优化
echo 'vm.swappiness=10' | sudo tee -a /etc/sysctl.conf
sudo sysctl -p
```

### 3. 日志分析

```bash
# 分析访问日志
awk '{print $1}' /var/log/nginx/access.log | sort | uniq -c | sort -nr

# 分析错误日志
grep "ERROR" /var/log/nginx/error.log

# 分析后端日志
grep "ERROR" /var/log/syslog | grep design-server
```

## 扩展部署

### 1. 负载均衡

```nginx
# Nginx负载均衡配置
upstream backend {
    server 127.0.0.1:8888;
    server 127.0.0.1:8889;
    server 127.0.0.1:8890;
}

server {
    listen 80;
    server_name your-domain.com;
    
    location /api/ {
        proxy_pass http://backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### 2. 集群部署

```yaml
# Kubernetes部署配置
apiVersion: apps/v1
kind: Deployment
metadata:
  name: design-backend
spec:
  replicas: 3
  selector:
    matchLabels:
      app: design-backend
  template:
    metadata:
      labels:
        app: design-backend
    spec:
      containers:
      - name: backend
        image: design-backend:latest
        ports:
        - containerPort: 8888
```

### 3. 微服务架构

```yaml
# 微服务Docker Compose配置
version: '3.8'
services:
  user-service:
    build: ./services/user
    ports:
      - "8001:8001"
  
  article-service:
    build: ./services/article
    ports:
      - "8002:8002"
  
  file-service:
    build: ./services/file
    ports:
      - "8003:8003"
```

## 部署检查清单

### 部署前检查
- [ ] 服务器环境满足要求
- [ ] 数据库已正确配置
- [ ] 域名已解析到服务器
- [ ] SSL证书已配置
- [ ] 防火墙规则已设置

### 部署后检查
- [ ] 后端服务正常启动
- [ ] 前端页面正常访问
- [ ] API接口正常响应
- [ ] 数据库连接正常
- [ ] 文件上传功能正常
- [ ] 日志记录正常
- [ ] 监控告警已配置

### 安全检查
- [ ] 敏感信息已加密
- [ ] 数据库密码已更改
- [ ] JWT密钥已设置
- [ ] CORS配置正确
- [ ] 文件权限设置正确
- [ ] SSL证书有效
- [ ] 防火墙规则正确

## 联系支持

**技术支持邮箱**: 1653765889@qq.com
