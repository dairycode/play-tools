# 游戏助手小程序

一个专为移动端设计的游戏助手小程序，目前实现了打牌记账功能，支持多人实时协作记账，自动计算最优结算方案。

## 项目结构

```
play-tools/
├── web/                # 前端代码（uniapp + TypeScript）
│   ├── src/
│   │   ├── api/       # API接口
│   │   ├── pages/     # 页面
│   │   ├── utils/     # 工具类
│   │   ├── types/     # 类型定义
│   │   └── components/# 组件
│   └── package.json
│
└── server/            # 后端代码（Golang + Gin + SQLite）
    ├── api/          # API处理器
    ├── model/        # 数据模型
    ├── database/     # 数据库
    ├── middleware/   # 中间件
    ├── utils/        # 工具类
    └── main.go
```

## 功能特性

### 1. 账号系统
- 用户注册/登录
- JWT Token + Cookie 双重认证
- 用户信息管理
- 自动登录态检查

### 2. 打牌记账助手（多人协作模式）

#### 创建房间
- 房主创建房间，输入房间名称和玩家人数（2-10人）
- 系统生成8位短房间ID，方便分享
- 房主自动加入房间并处于准备状态

#### 邀请好友
- 通过分享链接邀请好友加入房间
- 支持本地和局域网访问
- 自动携带房间ID和加入标识

#### 准备与开始
- 所有成员需要点击"准备"才能开始游戏
- 房主在所有人准备后可以开始游戏
- 游戏中成员自由退出或加入

#### 实时记账
- 每个人只能记录自己向别人的转账
- 点击成员头像快速选择转账对象
- 快捷金额选择（10/20/50/100/200/500）
- 所有记录实时同步给房间内所有成员
- 自动计算每个人的当前余额

#### 智能结算
- 只有房主可以触发结算
- 系统根据所有转账记录自动计算最优结算方案
- 使用贪心算法最小化转账次数
- 所有成员都可以查看结算结果
- 清晰的结算流程展示

### 3. 网络访问
- 支持本地访问（localhost）
- 支持局域网访问（自动检测IP地址）
- 前后端自动适配访问地址

## 技术栈

### 前端
- uniapp（Vue3 + TypeScript）
- 适配移动端H5/小程序
- 响应式设计
- 实时数据刷新（3秒轮询）

### 后端
- Golang 1.21+
- Gin Web框架
- GORM ORM
- SQLite数据库
- JWT认证
- bcrypt密码加密
- CORS跨域支持

## 快速开始

详细的安装和使用说明请查看 [QUICKSTART.md](QUICKSTART.md)

### 前端启动

```bash
cd web

# 安装依赖
npm install

# 开发模式（H5）- 支持局域网访问
npm run dev:h5

# 开发模式（微信小程序）
npm run dev:mp-weixin

# 构建生产版本
npm run build:h5
```

### 后端启动

```bash
cd server

# 下载依赖
go mod download

# 运行服务器 - 支持局域网访问
go run main.go
```

服务器将在 `http://0.0.0.0:8080` 启动，支持：
- 本地访问：http://localhost:8080
- 局域网访问：http://192.168.x.x:8080

## API 接口文档

### 用户相关

#### 注册
```http
POST /api/user/register
Content-Type: application/json

{
  "username": "string",
  "password": "string",
  "nickname": "string"
}
```

#### 登录
```http
POST /api/user/login
Content-Type: application/json

{
  "username": "string",
  "password": "string"
}

响应:
{
  "code": 200,
  "data": {
    "token": "jwt_token_string",
    "user": {
      "id": 1,
      "username": "string",
      "nickname": "string"
    }
  }
}

注意：登录成功后会自动设置 HttpOnly Cookie
```

#### 获取用户信息（需要认证）
```http
GET /api/user/info
Authorization: Bearer {token}
或使用 Cookie 自动携带
```

#### 更新用户信息（需要认证）
```http
PUT /api/user/update
Authorization: Bearer {token}

{
  "nickname": "string"
}
```

### 游戏相关（需要认证）

#### 创建房间
```http
POST /api/game/room/create

{
  "name": "周末德州局",
  "maxPlayers": 5
}

响应:
{
  "code": 200,
  "data": {
    "id": "a1b2c3d4",
    "name": "周末德州局",
    "ownerId": 1,
    "maxPlayers": 5,
    "status": "waiting"
  }
}
```

#### 加入房间
```http
POST /api/game/room/join

{
  "roomId": "a1b2c3d4"
}
```

#### 获取房间信息
```http
GET /api/game/room/:roomId

响应:
{
  "code": 200,
  "data": {
    "id": "a1b2c3d4",
    "name": "周末德州局",
    "ownerId": 1,
    "maxPlayers": 5,
    "status": "playing",
    "members": [
      {
        "userId": 1,
        "nickname": "张三",
        "isReady": true,
        "balance": -100
      },
      {
        "userId": 2,
        "nickname": "李四",
        "isReady": true,
        "balance": 50
      }
    ],
    "transactions": [
      {
        "id": 1,
        "roomId": "a1b2c3d4",
        "fromUserId": 1,
        "fromNickname": "张三",
        "toUserId": 2,
        "toNickname": "李四",
        "amount": 100,
        "createdAt": "2024-01-01T00:00:00Z"
      }
    ],
    "isOwner": true,
    "isMember": true
  }
}
```

#### 切换准备状态
```http
POST /api/game/room/:roomId/ready
```

#### 开始游戏
```http
POST /api/game/room/:roomId/start

注意：只有房主可以开始游戏，且所有成员必须准备
```

#### 添加交易记录
```http
POST /api/game/transaction/add

{
  "roomId": "a1b2c3d4",
  "toUserId": 2,
  "amount": 100
}

注意：fromUserId 自动从登录用户获取
```

#### 结束游戏并结算
```http
POST /api/game/room/finish

{
  "roomId": "a1b2c3d4"
}

响应:
{
  "code": 200,
  "data": [
    {
      "fromUserId": 1,
      "fromNickname": "张三",
      "toUserId": 2,
      "toNickname": "李四",
      "amount": 100
    }
  ]
}

注意：只有房主可以结束游戏
```

#### 获取房间结算结果
```http
GET /api/game/room/:roomId/settlements

注意：只有房间状态为 finished 才能查看
```

#### 获取我的房间列表
```http
GET /api/game/rooms

响应: 返回所有我加入过的房间
```

## 数据模型

### User（用户）
```go
type User struct {
    ID        uint
    Username  string  // 用户名（唯一）
    Password  string  // 加密后的密码
    Nickname  string  // 昵称
    CreatedAt time.Time
}
```

### GameRoom（游戏房间）
```go
type GameRoom struct {
    ID         string    // 房间ID（8位短ID）
    Name       string    // 房间名称
    OwnerID    uint      // 房主ID
    MaxPlayers int       // 最大玩家数（2-10）
    Status     string    // waiting/playing/finished
    CreatedAt  time.Time
}
```

### RoomMember（房间成员）
```go
type RoomMember struct {
    ID       uint
    RoomID   string  // 所属房间
    UserID   uint    // 用户ID
    Nickname string  // 昵称（冗余存储）
    IsReady  bool    // 是否准备
    JoinedAt time.Time
}
```

### Transaction（交易记录）
```go
type Transaction struct {
    ID           uint
    RoomID       string
    FromUserID   uint    // 付款人用户ID
    ToUserID     uint    // 收款人用户ID
    Amount       int
    FromNickname string  // 付款人昵称（冗余）
    ToNickname   string  // 收款人昵称（冗余）
    CreatedAt    time.Time
}
```

### Settlement（结算记录）
```go
type Settlement struct {
    ID           uint
    RoomID       string
    FromUserID   uint
    FromNickname string
    ToUserID     uint
    ToNickname   string
    Amount       int
    CreatedAt    time.Time
}
```

## 核心算法

### 结算算法（贪心算法）

项目使用贪心算法计算最优结算方案，确保最少的转账次数：

1. 根据所有交易记录计算每个人的余额
2. 将玩家分为欠款人（负余额）和收款人（正余额）两组
3. 对两组玩家按余额绝对值从大到小排序
4. 依次匹配最大欠款人和最大收款人
5. 每次转账金额为两者余额绝对值的较小值
6. 更新余额后继续匹配，直到所有账务结清

**示例：**
```
初始余额：A: -150, B: +100, C: +50
结算方案：A 向 B 转 100，A 向 C 转 50
仅需 2 次转账完成结算
```

## 页面说明

1. [登录页](web/src/pages/login/login.vue) - 用户登录/注册，支持跳转回原页面
2. [首页](web/src/pages/index/index.vue) - 游戏助手选择入口
3. [账号管理](web/src/pages/account/account.vue) - 修改昵称、退出登录
4. [创建房间](web/src/pages/poker/create-room.vue) - 设置房间名称和玩家人数
5. [加入房间](web/src/pages/poker/join-room.vue) - 输入房间ID加入
6. [游戏房间](web/src/pages/poker/game-room.vue) - 实时记账、查看成员、准备状态
7. [结算页面](web/src/pages/poker/settlement.vue) - 显示最优结算方案

## 设计特点

### 用户体验
- 首页作为入口，未登录自动跳转登录页
- 登录后返回原页面（支持分享链接）
- 单屏设计，无需滚动
- 移动端优先，触摸友好
- 实时反馈，操作流畅

### 技术特点
- Cookie + JWT 双重认证
- 前后端自动适配局域网IP
- 实时数据同步（3秒轮询）
- 权限清晰（房主/成员权限分离）
- 可追溯（每笔交易都有发起人）

### 安全性
- 密码 bcrypt 加密
- HttpOnly Cookie 防止 XSS
- JWT Token 验证
- CORS 跨域保护

## 多人协作流程

### 1. 房主创建房间
- 设置房间名称和人数上限
- 系统生成8位房间ID
- 房主自动加入并准备

### 2. 邀请好友加入
- 分享房间链接或房间ID
- 好友点击链接自动加入
- 未登录自动跳转登录后返回

### 3. 准备开始
- 所有成员点击"准备"
- 房主确认开始游戏

### 4. 记账过程
- 每个人记录自己的转账
- 选择收款人和金额
- 所有人实时看到更新

### 5. 结算
- 房主点击结束游戏
- 系统自动计算最优方案
- 所有成员查看结算结果

## 后续规划

- [ ] 支持房间密码保护
- [ ] 添加历史记录详情查看
- [ ] 支持导出结算结果
- [ ] 添加统计分析功能
- [ ] 支持更多游戏工具
- [ ] 添加主题切换
- [ ] 微信小程序版本优化

## 许可证

MIT License
