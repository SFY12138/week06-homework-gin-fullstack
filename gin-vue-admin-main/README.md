# Gin-Vue-Admin 二次开发项目

## 项目版本

- **项目版本**：v2.8.9

## 开发任务索引

### 任务1：环境搭建与初始化

- **任务描述**：放弃 MySQL/PostgreSQL，使用轻量级数据库 SQLite 启动项目，并完成系统初始化
- **完成状态**：已完成
- **关键修改**：
  1. 修改后端 `config.yaml` 配置文件，指定数据库类型为 `sqlite`
  2. 修改数据库初始化逻辑，解决 SQLite DSN 路径拼接问题
  3. 成功启动前后端项目，完成系统初始化

### 任务2：用户行为追踪功能开发

- **任务描述**：在"用户管理"列表页中新增两列，登录IP和登录时间
- **完成状态**：已完成
- **关键修改**：
  1. 在用户模型 `sys_user.go` 中添加 `LoginIP` 和 `LoginTime` 字段
  2. 修改用户服务层 `sys_user.go` 的 `SetUserInfo` 方法，支持更新登录信息
  3. 修改 API 层 `sys_user.go` 的 `TokenNext` 方法，在用户登录成功后更新登录IP和登录时间
  4. 修改前端用户管理页面 `user.vue`，添加登录IP和登录时间列，并格式化显示时间为 YYYY-MM-DD HH:mm 格式

## 核心技术实现

### 1. 数据库迁移（MySQL to SQLite）

**实现思路**：
- 通过修改 `config.yaml` 中的 `db-type` 配置，将数据库类型从 MySQL 切换为 SQLite
- 使用 GORM 的 SQLite 驱动 `github.com/glebarez/sqlite` 实现数据库连接
- 修改 `SqliteEmptyDsn()` 方法，正确处理数据库文件路径，避免路径拼接错误

**关键代码**：
```go
func (i *InitDB) SqliteEmptyDsn() string {
    // 如果DBPath已经包含.db后缀，直接返回
    if len(i.DBPath) > 3 && i.DBPath[len(i.DBPath)-3:] == ".db" {
        return i.DBPath
    }
    // 否则，拼接路径和数据库名
    separator := string(os.PathSeparator)
    return i.DBPath + separator + i.DBName + ".db"
}
```

### 2. 用户行为追踪（登录IP和登录时间）

**实现思路**：
- 在用户模型中扩展 `LoginIP` 和 `LoginTime` 字段，用于记录用户最后一次登录的信息
- 在用户登录成功后，通过中间件或 API 层获取客户端 IP 地址，并更新到用户记录中
- 使用 GORM 的 `UpdatedAt` 或单独的 `LoginTime` 字段记录登录时间
- 前端通过表格列展示这两个字段，并格式化为易读的日期时间格式

**后端关键代码**：
1. 用户模型扩展：
```go
type SysUser struct {
    // ... 其他字段
    LoginIP       string    `json:"loginIp" gorm:"comment:最后登录IP"`
    LoginTime     *time.Time `json:"loginTime" gorm:"comment:最后登录时间"`
}
```

2. 登录时更新登录信息：
```go
// 在 TokenNext 或 Login 方法中
loginIP := utils.GetClientIp(c)
now := time.Now()
global.GVA_DB.Model(&SysUser{}).Where("username = ?", username).Updates(map[string]interface{}{
    "login_ip":   loginIP,
    "login_time": now,
})
```

**前端关键代码**：
```javascript
// 时间格式化
formatDate(row.loginTime) {
    if (!row.loginTime) return '-'
    const date = new Date(row.loginTime)
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    const hours = String(date.getHours()).padStart(2, '0')
    const minutes = String(date.getMinutes()).padStart(2, '0')
    return `${year}-${month}-${day} ${hours}:${minutes}`
}
```

## 项目启动说明

### 前端启动

```bash
cd week06/homework/gin-fullstack/gin-vue-admin-main/web
npm run dev
```

### 后端启动

```bash
cd week06/homework/gin-fullstack/gin-vue-admin-main/server
# 设置 CGO_ENABLED 环境变量
$env:CGO_ENABLED=1
go run main.go
```

### 系统初始化

1. 访问 http://localhost:8080/
2. 在初始化页面中选择数据库类型为 `sqlite`
3. 填写 dbPath（例如：`./gva.db`）和 dbName（例如：`gva`）
4. 点击"立即初始化"完成系统初始化
5. 使用管理员账号登录（用户名：`admin`，密码：`123456`）

## 功能演示

### 用户行为追踪功能演示

1. 使用管理员账户登录
2. 进入"用户管理"页面
3. 新增一个用户（例如：`user_a`）
4. 使用第二个浏览器或隐私模式登录新创建的用户 `user_a`
5. 回到管理员界面，刷新用户管理列表
6. 可以看到用户 `user_a` 的"登录IP"和"登录时间"已经显示出来

## 技术栈

- **后端**：Golang 1.20 + Gin 1.9.1 + Gorm 1.25.2
- **前端**：Vue 3.3.4 + Element Plus 2.3.8 + Vite
- **数据库**：SQLite
- **认证**：JWT + Casbin
- **文档**：Swagger

## 项目结构

```
week06/homework/gin-fullstack/gin-vue-admin-main/
├── server/                 # 后端项目
│   ├── api/               # API 层
│   ├── config/            # 配置包
│   ├── model/             # 数据模型
│   ├── service/           # 业务逻辑层
│   ├── router/            # 路由层
│   ├── middleware/        # 中间件
│   └── utils/             # 工具包
├── web/                   # 前端项目
│   ├── src/
│   │   ├── api/           # API 调用
│   │   ├── view/          # 页面组件
│   │   └── permission.js # 权限管理
│   └── package.json
├── README_OLD.md          # 原项目文档（历史参考）
└── README.md              # 本次开发项目文档
```

## 注意事项

1. SQLite 数据库文件路径应使用相对路径或绝对路径，确保后端服务有写入权限
2. 确保 CGO_ENABLED=1 环境变量已设置，否则 SQLite 驱动无法正常工作
3. 登录IP获取使用 `utils.GetClientIp(c)` 方法，需要根据实际部署环境调整代理设置
4. 前端时间格式化需要处理 null 值情况，避免显示 NaN
