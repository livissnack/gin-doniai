# gin-doniai

## Windows环境下一键打包linux二进制运行包

```shell
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gin-doniai main.go
```

## 创建服务文件

```shell
vim /etc/systemd/system/discuss-web.service
```

```shell
[Unit]
Description=System Discuss Service
After=network.target
Wants=network.target

[Service]
Type=simple
User=www-data
Group=www-data
WorkingDirectory=/var/www/gin-doniai
ExecStart=/var/www/gin-doniai/gin-doniai 8080
ExecReload=/bin/kill -HUP $MAINPID
Restart=always
RestartSec=10

# 安全设置
NoNewPrivileges=yes
PrivateTmp=yes
ProtectSystem=strict
ProtectHome=yes
ReadWritePaths=/var/www/gin-doniai /var/log
ProtectKernelTunables=yes
ProtectKernelModules=yes
ProtectControlGroups=yes

# 资源限制
MemoryMax=100M
CPUQuota=30%

# 日志
StandardOutput=journal
StandardError=journal
SyslogIdentifier=gin-doniai


[Install]
WantedBy=multi-user.target
```

## 启用服务

```shell
# 重新加载 systemd 配置
sudo systemctl daemon-reload

# 启用开机自启
sudo systemctl enable discuss-web.service

# 启动服务
sudo systemctl start discuss-web.service

# 查看状态
sudo systemctl status discuss-web.service

# 查看日志
sudo journalctl -u discuss-web.service -f
```

![电脑1](https://pic.114156.xyz/uploads/0EXDab7wpK64.webp)
![电脑2](https://pic.114156.xyz/uploads/GVpQ3nTdSmlO.webp)
![电脑3](https://pic.114156.xyz/uploads/FO87l7BCs4tv.webp)
![电脑4](https://pic.114156.xyz/uploads/6f8jUrl0Ra4X.webp)
![电脑5](https://pic.114156.xyz/uploads/jj4-UENslx2m.webp)
![电脑6](https://pic.114156.xyz/uploads/Vy7ylFAyjrvt.webp)
![电脑7](https://pic.114156.xyz/uploads/_C5yxyLmCE77.webp)
![电脑8](https://pic.114156.xyz/uploads/6WF1S7YVOiAR.webp)
![电脑9](https://pic.114156.xyz/uploads/nS8RacZ8GJtF.webp)


bash /root/.acme.sh/acme.sh --issue --dns dns_ali -d livissnack.com -d "*.livissnack.com" --set-default-ca --server letsencrypt