# Graduation Project Deployment Documentation

### Description

> 1. Default Account Password  
    1.1 Username: admin  
    2.2 Password: 123456  

## Deployment Overview

This document provides detailed instructions on how to deploy a full-stack project based on Gin + Vue3 to production environment. The project adopts a frontend-backend separation architecture and supports multiple deployment methods.

## Deployment Architecture

```
Production Environment Architecture:
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Nginx (80)    │    │   Backend       │    │   MySQL         │
│   (Load Balancer)│────│   (8888)        │────│   (3306)        │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │
         │                       │
         ▼                       ▼
┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   Redis         │
│   (Static Files)│    │   (Cache)       │
└─────────────────┘    └─────────────────┘
```

## Environment Requirements

### Server Configuration
- **CPU**: 2+ cores
- **Memory**: 4GB+
- **Storage**: 50GB+
- **Operating System**: Ubuntu 20.04+ / CentOS 7+ / Windows Server 2019+

### Software Requirements
- **Go**: 1.23+
- **Node.js**: 16+
- **MySQL**: 8.0+
- **Nginx**: 1.18+
- **Redis**: 6.0+ (optional)

## Deployment Methods

### Method 1: Traditional Deployment

#### 1. Server Environment Preparation

##### Ubuntu/Debian System
```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install basic tools
sudo apt install -y curl wget git unzip

# Install Go
wget https://go.dev/dl/go1.23.6.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.23.6.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Install Node.js
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs

# Install MySQL
sudo apt install -y mysql-server
sudo systemctl start mysql
sudo systemctl enable mysql

# Install Nginx
sudo apt install -y nginx
sudo systemctl start nginx
sudo systemctl enable nginx
```

##### CentOS/RHEL System
```bash
# Update system
sudo yum update -y

# Install basic tools
sudo yum install -y curl wget git unzip

# Install Go
wget https://go.dev/dl/go1.23.6.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.23.6.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Install Node.js
curl -fsSL https://rpm.nodesource.com/setup_18.x | sudo bash -
sudo yum install -y nodejs

# Install MySQL
sudo yum install -y mysql-server
sudo systemctl start mysqld
sudo systemctl enable mysqld

# Install Nginx
sudo yum install -y nginx
sudo systemctl start nginx
sudo systemctl enable nginx
```

#### 2. Database Configuration

```bash
# Login to MySQL
sudo mysql -u root -p

# Create database and user
CREATE DATABASE design_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'design_user'@'localhost' IDENTIFIED BY 'your_secure_password';
GRANT ALL PRIVILEGES ON design_db.* TO 'design_user'@'localhost';
FLUSH PRIVILEGES;
EXIT;
```

#### 3. Backend Deployment

```bash
# Clone project
git clone <project-url>
cd project/backend

# Install dependencies
go mod tidy
go mod download

# Build project
go build -o design-server main.go

# Create configuration file
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

# Create service directory
sudo mkdir -p /opt/design
sudo cp design-server /opt/design/
sudo cp config.yaml /opt/design/
sudo mkdir -p /opt/design/uploads

# Create systemd service
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

# Start service
sudo systemctl daemon-reload
sudo systemctl enable design-server
sudo systemctl start design-server
sudo systemctl status design-server
```

#### 4. Frontend Deployment

```bash
# Enter frontend directory
cd ../frontend

# Install dependencies
npm install

# Modify API address
sed -i 's|http://localhost:8888|http://your-domain.com|g' src/utils/request.ts

# Build project
npm run build

# Deploy to Nginx
sudo cp -r dist/* /var/www/html/
sudo chown -R www-data:www-data /var/www/html/
```

#### 5. Nginx Configuration

```bash
# Create Nginx configuration
sudo tee /etc/nginx/sites-available/design << EOF
server {
    listen 80;
    server_name your-domain.com;
    
    # Frontend static files
    location / {
        root /var/www/html;
        try_files \$uri \$uri/ /index.html;
        index index.html;
    }
    
    # Backend API proxy
    location /api/ {
        proxy_pass http://127.0.0.1:8888;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }
    
    # File upload
    location /uploads/ {
        alias /opt/design/uploads/;
        expires 30d;
        add_header Cache-Control "public, immutable";
    }
    
    # Security configuration
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header Referrer-Policy "no-referrer-when-downgrade" always;
    add_header Content-Security-Policy "default-src 'self' http: https: data: blob: 'unsafe-inline'" always;
}
EOF

# Enable configuration
sudo ln -s /etc/nginx/sites-available/design /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### Method 2: Docker Deployment

#### 1. Install Docker

```bash
# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/download/v2.20.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

#### 2. Create Dockerfile

**Backend Dockerfile**
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

**Frontend Dockerfile**
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

#### 3. Create docker-compose.yml

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

#### 4. Deployment Commands

```bash
# Build and start all services
docker-compose up -d

# Check service status
docker-compose ps

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

### Method 3: Cloud Platform Deployment

#### 1. Alibaba Cloud ECS Deployment

```bash
# Use Alibaba Cloud ECS, follow traditional deployment method
# Configure security groups, open ports 80, 443, 8888
```

#### 2. Tencent Cloud CVM Deployment

```bash
# Use Tencent Cloud CVM, follow traditional deployment method
# Configure security groups, open corresponding ports
```

#### 3. AWS EC2 Deployment

```bash
# Use AWS EC2, follow traditional deployment method
# Configure security groups, open corresponding ports
```



## Monitoring and Maintenance

### 1. Log Management

```bash
# View backend logs
sudo journalctl -u design-server -f

# View Nginx logs
sudo tail -f /var/log/nginx/access.log
sudo tail -f /var/log/nginx/error.log

# View Docker logs
docker-compose logs -f backend
```

### 2. Performance Monitoring

```bash
# Install monitoring tools
sudo apt install -y htop iotop nethogs

# Monitor system resources
htop
iotop
nethogs

# Monitor database
mysql -u root -p -e "SHOW PROCESSLIST;"
```

### 3. Backup Strategy

```bash
# Database backup script
#!/bin/bash
BACKUP_DIR="/backup/mysql"
DATE=$(date +%Y%m%d_%H%M%S)
mysqldump -u design_user -p design_db > $BACKUP_DIR/design_db_$DATE.sql

# File backup script
#!/bin/bash
BACKUP_DIR="/backup/uploads"
DATE=$(date +%Y%m%d_%H%M%S)
tar -czf $BACKUP_DIR/uploads_$DATE.tar.gz /opt/design/uploads/
```

### 4. Auto Update

```bash
# Create update script
#!/bin/bash
cd /opt/design
git pull origin main
go build -o design-server main.go
sudo systemctl restart design-server
```

## Security Configuration

### 1. Firewall Configuration

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

### 2. SSL Certificate Configuration

```bash
# Install Certbot
sudo apt install certbot python3-certbot-nginx

# Obtain SSL certificate
sudo certbot --nginx -d your-domain.com

# Auto renewal
sudo crontab -e
# Add: 0 12 * * * /usr/bin/certbot renew --quiet
```

### 3. Security Hardening

```bash
# Disable root login
sudo sed -i 's/#PermitRootLogin yes/PermitRootLogin no/' /etc/ssh/sshd_config

# Change SSH port
sudo sed -i 's/#Port 22/Port 2222/' /etc/ssh/sshd_config

# Restart SSH service
sudo systemctl restart sshd
```

## Troubleshooting

### 1. Common Issues

**Issue 1: Database Connection Failed**
```bash
# Check MySQL service status
sudo systemctl status mysql

# Check database connection
mysql -u design_user -p design_db

# Check configuration file
cat /opt/design/config.yaml
```

**Issue 2: Frontend Cannot Access Backend API**
```bash
# Check backend service status
sudo systemctl status design-server

# Check if port is open
netstat -tlnp | grep 8888

# Check Nginx configuration
sudo nginx -t
```

**Issue 3: File Upload Failed**
```bash
# Check upload directory permissions
ls -la /opt/design/uploads/

# Modify permissions
sudo chown -R www-data:www-data /opt/design/uploads/
sudo chmod -R 755 /opt/design/uploads/
```

### 2. Performance Optimization

```bash
# Database optimization
mysql -u root -p
SHOW VARIABLES LIKE 'max_connections';
SET GLOBAL max_connections = 200;

# Nginx optimization
sudo nano /etc/nginx/nginx.conf
# Adjust worker_processes and worker_connections

# System optimization
echo 'vm.swappiness=10' | sudo tee -a /etc/sysctl.conf
sudo sysctl -p
```

### 3. Log Analysis

```bash
# Analyze access logs
awk '{print $1}' /var/log/nginx/access.log | sort | uniq -c | sort -nr

# Analyze error logs
grep "ERROR" /var/log/nginx/error.log

# Analyze backend logs
grep "ERROR" /var/log/syslog | grep design-server
```

## Extended Deployment

### 1. Load Balancing

```nginx
# Nginx load balancing configuration
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

### 2. Cluster Deployment

```yaml
# Kubernetes deployment configuration
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

### 3. Microservice Architecture

```yaml
# Microservice Docker Compose configuration
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

## Deployment Checklist

### Pre-deployment Check
- [ ] Server environment meets requirements
- [ ] Database is properly configured
- [ ] Domain is resolved to server
- [ ] SSL certificate is configured
- [ ] Firewall rules are set

### Post-deployment Check
- [ ] Backend service starts normally
- [ ] Frontend pages are accessible
- [ ] API interfaces respond normally
- [ ] Database connection is normal
- [ ] File upload function works
- [ ] Log recording is normal
- [ ] Monitoring alerts are configured

### Security Check
- [ ] Sensitive information is encrypted
- [ ] Database password is changed
- [ ] JWT secret is set
- [ ] CORS configuration is correct
- [ ] File permissions are set correctly
- [ ] SSL certificate is valid
- [ ] Firewall rules are correct

## Technical Support

**Technical Support Email**: 1653765889@qq.com
