# SM2 协同签名服务

基于 Go 语言实现的 SM2 协同签名服务后端，实现一个基于协同签名模式的安全密码服务系统。服务端存储密钥分量 D2，客户端持有密钥分量 D1，通过多轮交互完成签名和解密操作。

## 产品概述

重构 SM2 协同签名服务后端，实现一个基于协同签名模式的安全密码服务系统。服务端存储密钥分量 D2，客户端持有密钥分量 D1，通过多轮交互完成签名和解密操作。

## 核心功能

- **用户注册**：自动生成 SM2 密钥对（协同签名模式），服务端存储 D2 分量
- **协同签名**：服务端参与签名计算，返回签名分量 r, s2, s3
- **协同解密**：服务端参与解密计算，返回中间密文 T2
- **Token 认证**：基于 Token 的会话管理，所有状态持久化到 SQLite3
- **管理接口**：用户管理（增删改查）、密钥管理、系统管理（健康检查、日志查询）

## 技术栈

- **语言**: Go 1.24.12
- **Web 框架**: Fiber v2
- **数据库**: SQLite3 (使用 modernc.org/sqlite 纯 Go 实现)
- **国密算法**: github.com/emmansun/gmsm (纯 Go 实现的 SM2/SM3/SM4)
- **配置管理**: Viper
- **日志**: Zap (结构化日志)
- **API 文档**: OpenAPI 3.0 (YAML + Markdown)

## 项目结构

```
backend/
├── cmd/
│   └── server/          # 服务入口
│       └── main.go
├── internal/
│   ├── config/          # 配置管理
│   │   └── config.go
│   ├── handler/         # 请求处理
│   │   ├── handler.go
│   │   ├── user.go
│   │   ├── cosign.go
│   │   └── admin.go
│   ├── middleware/      # 中间件
│   │   └── auth.go
│   ├── model/           # 数据模型
│   │   ├── user.go
│   │   ├── key.go
│   │   ├── session.go
│   │   └── audit_log.go
│   ├── repository/      # 数据访问
│   │   ├── repository.go
│   │   ├── user_repo.go
│   │   ├── key_repo.go
│   │   ├── session_repo.go
│   │   └── audit_log_repo.go
│   ├── service/         # 业务逻辑
│   │   ├── user_service.go
│   │   ├── cosign_service.go
│   │   └── crypto_service.go
│   └── crypto/          # 密码服务
│       └── sm2_coop.go
├── pkg/
│   ├── response/        # 统一响应格式
│   │   └── response.go
│   └── utils/           # 工具函数
│       └── utils.go
├── docs/                # API 文档
│   ├── api.yaml         # OpenAPI 3.0 文档
│   └── api.md           # Markdown API 文档
├── scripts/             # 脚本文件
│   └── schema.sql       # 数据库初始化脚本
├── config.yaml          # 配置文件
├── go.mod               # Go 模块定义
├── go.sum               # 依赖校验
├── Makefile             # 构建脚本
└── README.md            # 项目说明
```

## SM2 协同签名协议

### 密钥生成流程

1. 客户端生成 d1 → 计算 P1 = d1 * G → 发送 P1 到服务端
2. 服务端生成 d2 → 计算 d2Inv = d2^(-1) mod n
3. 服务端计算 P2 = d2Inv * G
4. 服务端计算 Pa = d2Inv * P1 + (n-1) * G（协同公钥）
5. 服务端存储 (d2, d2Inv, Pa)，返回 (P2, Pa)
6. 客户端计算完整私钥 d = d1 * d2 - 1

### 协同签名流程

1. 客户端发送 (Q1, E) 到服务端
2. 服务端生成随机 (k2, k3)
3. 服务端计算 Q2 = k2 * G, x1 = k3 * Q1 + Q2
4. 服务端计算 r = E + x1 mod n
5. 服务端计算 s2 = d2 * k3 mod n, s3 = d2 * (r + k2) mod n
6. 服务端返回给客户端
7. 客户端组合生成完整签名

### 协同解密流程

1. 客户端发送 T1（密文点 C1 的变换）到服务端
2. 服务端计算 T2 = d2Inv * T1
3. 服务端返回 T2 给客户端
4. 客户端使用 T2 计算对称密钥解密明文

## 构建和运行

### 构建

```bash
# 构建二进制文件
make build
```

### 运行

```bash
# 运行服务
make run

# 或直接运行二进制文件
./bin/sm2-co-sign-server
```

### 测试

```bash
# 运行所有测试
make test
```

### 依赖管理

```bash
# 下载依赖
make deps
```

### 清理

```bash
# 清理构建产物
make clean
```

## 配置

配置文件位于 `config.yaml`，可以通过环境变量覆盖配置项。

### 主要配置项

- `server.port`: 服务端口
- `database.path`: SQLite3 数据库文件路径
- `jwt.secret`: JWT 签名密钥
- `jwt.expiresIn`: Token 过期时间
- `crypto.masterKey`: 主密钥（用于加密存储密钥分量）

## API 接口

### 接口前缀

- 管理接口前缀：`/mapi`
- 业务接口前缀：`/api`

### 文档

- OpenAPI 3.0 文档：`docs/api.yaml`
- Markdown API 文档：`docs/api.md`

## 安全注意事项

- 确保使用 HTTPS 协议
- 配置合适的访问控制
- 定期更新密钥
- 监控异常访问
- 密钥分量使用主密钥加密存储

## 部署

### 二进制部署

1. 构建二进制文件
2. 复制配置文件到相应目录
3. 运行二进制文件

### Docker 部署

（可选）使用 Docker 构建和运行服务。

## 许可证

Apache License 2.0
