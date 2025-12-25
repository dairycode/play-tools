# 快速开始指南

## 环境要求

### 前端
- Node.js 16+
- npm 或 yarn

### 后端
- Go 1.21+

## 安装步骤

### 1. 启动后端服务

```bash
# 进入后端目录
cd server

# 初始化Go模块并下载依赖
go mod download

# 运行服务器
go run main.go
```

后端服务将在 `http://localhost:8080` 启动。

首次运行时会自动创建 SQLite 数据库文件 `play-tools.db` 并初始化数据表。

### 2. 启动前端服务

```bash
# 进入前端目录
cd web

# 安装依赖
npm install

# 启动开发服务器（H5模式）
npm run dev:h5
```

前端服务将在浏览器中自动打开，通常是 `http://localhost:80`

### 3. 开始使用

1. 在浏览器中打开前端地址
2. 注册一个新账号
3. 登录后进入首页
4. 点击"打牌记账"卡片
5. 创建房间并开始记账

## 开发模式

### H5 开发
```bash
cd web
npm run dev:h5
```

### 微信小程序开发
```bash
cd web
npm run dev:mp-weixin
```
然后使用微信开发者工具打开 `web/dist/dev/mp-weixin` 目录

## 生产构建

### 前端构建
```bash
cd web
npm run build:h5  # 构建H5版本
# 或
npm run build:mp-weixin  # 构建微信小程序版本
```

### 后端构建
```bash
cd server
go build -o play-tools main.go
./play-tools  # 运行编译后的程序
```

## 常见问题

### 1. 前端无法连接后端

确认：
- 后端服务是否正常运行在 8080 端口
- 检查 [web/src/utils/request.ts](../web/src/utils/request.ts:2) 中的 BASE_URL 配置

### 2. 数据库初始化失败

确认：
- server 目录有写入权限
- Go 版本是否符合要求

### 3. 前端依赖安装失败

尝试：
```bash
# 清除缓存
rm -rf node_modules package-lock.json
npm cache clean --force
npm install
```

## API测试

你可以使用 curl 或 Postman 测试API：

```bash
# 注册用户
curl -X POST http://localhost:8080/api/user/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"123456","nickname":"测试用户"}'

# 登录
curl -X POST http://localhost:8080/api/user/login \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"123456"}'
```

## 项目结构说明

```
play-tools/
├── web/                    # 前端项目
│   ├── src/
│   │   ├── api/           # API调用封装
│   │   ├── pages/         # 页面组件
│   │   │   ├── login/     # 登录页
│   │   │   ├── index/     # 首页
│   │   │   ├── account/   # 账号管理
│   │   │   └── poker/     # 打牌记账相关页面
│   │   ├── utils/         # 工具函数
│   │   └── types/         # TypeScript类型定义
│   └── package.json
│
└── server/                 # 后端项目
    ├── api/               # API处理器
    ├── model/             # 数据模型
    ├── database/          # 数据库初始化
    ├── middleware/        # 中间件（如JWT认证）
    ├── utils/             # 工具函数
    └── main.go            # 入口文件
```

## 下一步

- 阅读 [README.md](../README.md) 了解详细功能
- 查看源码了解实现细节
- 根据需求进行二次开发
