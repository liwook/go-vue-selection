# 分层结构体使用规范：handler / service / repository / types

> 本文讲清各层之间**参数与返回值用什么类型**，避免 `model`（数据库实体）泄漏到 API 响应，或 `types.ParamXxx` 被误用作返回值。读完记住一句话：**入参走 `types.ParamXxx`，数据库走 `model.Xxx`，出参走 `types.Xxx`/`types.ResponseXxx`，`model→types` 转换只在 service 发生，响应永远过 `pkg/result`，`ctx` 从 handler 一路透传到 repo。**

---

## 一、分层与依赖方向

```text
handler  ──►  service  ──►  repository  ──►  dal/model（数据库实体）
   │           │              │
   └───────────┴──────────────┴── 四层都依赖 types（公共 DTO），types 不依赖任何业务层
```

- **依赖只能单向向下**，下层绝不反向依赖上层。
- **handler（接口适配层）**：绑定 `types.ParamXxx`、调 service、用 `pkg/result` 包装响应。不写业务、不碰 DB、不构造 `model`。
- **service（业务逻辑层）**：业务规则、事务边界、跨表编排、以及**唯一的 `model↔types` 转换点**。不知 HTTP 存在。
- **repository（数据访问层）**：只管 `model.Xxx` 的存/取（经 `dal/query`）。不含业务判断，错误用 `fmt.Errorf("...: %w", err)` 上抛。
- **types（公共 DTO）**：四层共享的契约层，不依赖任何业务层，是依赖图最底层。

**铁律**：`model.Xxx` 只活在 service 与 repository 之间，**绝不直接作为 HTTP 响应返回**。

---

## 二、`types` 包的四种结构体（命名即契约）

所有 API 相关结构体集中在 `types/` 包，按用途分四类：

| 类型 | 用途 | 关键约定 | 示例 |
|------|------|---------|------|
| `Xxx` | 对外可读的领域对象 | 嵌入 `BaseModel`；JSON **驼峰**；ID 一律 `string`；无 gorm tag | `User`、`Category3`、`Menu` |
| `ParamXxx` | 请求入参 | 带 `binding` tag；**禁止作为返回值** | `ParamUserLogin`、`ParamDoAssignRole` |
| `ResponseXxx` | 响应视图 | 用于拼装/聚合/分页；分页形如 `Records`+`Total/Size/Current/Pages` | `ResponseUserInfo`、`ResponseUserList` |
| `BaseModel` | 时间基类 | 嵌入所有 `Xxx`；`CreateTime/UpdateTime` 用 `time.Time`，JSON `,omitzero` | —— |

```go
type User struct {
    BaseModel
    UserID   string `json:"userId"`      // ID 必须 string，避免大整数精度丢失
    Username string `json:"username"`
    Status   bool   `json:"status"`
}

type ParamUserLogin struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

type BaseModel struct {
    CreateTime time.Time `json:"createTime,omitzero"`
    UpdateTime time.Time `json:"updateTime,omitzero"`
}
```

---

## 三、ctx 透传规范（贯穿三层）

`context.Context` 随一次请求下传，用于**超时控制、链路追踪、请求取消**。强制要求 `ctx` 从 HTTP 入口透传到最底层 DB 查询，**中间层不得凭空造 ctx**。

### 约定（4 条）

1. **来源唯一**：源头只能是 handler 的 `c.Request.Context()`。
2. **签名首位**：service/repository 每个方法**第一个参数是 `ctx context.Context`**。
3. **到底为止**：repo 内必须透传到 `WithContext(ctx)`，禁止 `context.Background()`。
4. **不存储、不改义**：ctx 只存在于调用参数链上，不存进结构体字段；需附加 traceID 等用 `context.WithValue` 派生子 ctx 后继续下传。

### 传递路径

```text
handler.Login(c)
   └─ ctx := c.Request.Context()        ← 唯一源头
        └─ service.Login(ctx, p)        ← ctx 作为首参透传
             └─ repo.Login(ctx, ...)    ← ctx 继续透传
                  └─ d.q.User.WithContext(ctx).First()   ← 落到 DB 调用
```

**禁止**：❌ repo 内 `context.Background()`；❌ 方法签名省略 `ctx`；❌ 把 ctx 存进结构体字段长期持有。

### 超时控制（handler 层示例）

```go
func (u *userHandler) Login(c *gin.Context) {
    ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
    defer cancel()
    token, err := u.userSvc.Login(ctx, p)
    // ...
}
```

> 即使 handler 不设 `WithTimeout`，也**必须**保证 `ctx` 透传到底，否则上游中间件设置的 deadline 在 repo 层失效。

---

## 四、一次请求的完整数据流（以登录为例）

```text
HTTP JSON
  │ c.ShouldBindJSON(p)               p := new(types.ParamUserLogin)
  ▼
handler.Login
  │ ctx := c.Request.Context()
  ▼
service.Login(ctx, p)
  │ user, _ := repo.Login(ctx, p.Username, p.Password)   // 返回 *model.User
  │ token := jwt.GenToken(user.UserID, user.Username)     // 用 model 字段，不出 types
  ▼ 返回 token string
handler
  │ result.Success(c, token)          // 响应永远过 pkg/result
  ▼
响应 { code, message, data: "eyJ..." }
```

可见：**入参 `types.ParamXxx`** 在 handler 止步，**`model.Xxx`** 在 service↔repo 之间流动，**出参 `types.Xxx`/`ResponseXxx`** 由 service 经 `convert.go` 产出，最终由 `pkg/result` 包装。

---

## 五、各层硬规则

### handler

- **入参**：`p := new(types.ParamXxx)` → `c.ShouldBindJSON(p)`
- **出参**：一律过 `pkg/result`（`result.Success` / `result.Error`）
- **ctx**：从 `c.Request.Context()` 取出，作为首参传给 service
- **不构造 `model.Xxx`**，只把 `types.ParamXxx` 往下传

### service

- **入参**：收 `*types.ParamXxx`（或基础类型）
- **内部**：用 `model.Xxx` 与 repository 交互
- **出参**：返回 `*types.ResponseXxx` / `*types.Xxx` / 普通类型
- **转换（唯一转换点）**：所有 `model→types` 必须走 `service/convert.go` 的 `modelToXxx(...)`，**禁止手写字段拷贝**

```go
func (u *userService) GetUserInfo(ctx context.Context, userID int64) (*types.ResponseUserInfo, error) {
    user, err := u.userRepo.GetUserById(ctx, userID)   // 返回 *model.User
    if err != nil {
        return nil, err
    }
    return convert.ModelToResponseUserInfo(user), nil  // 转换只在此处发生
}
```

### repository

- **入参**：`model.Xxx`；关联写入只收基础类型（如 `DoAssign(ctx, userId string, roleIds []string)`），**不收 `types.ParamXxx`**
- **出参**：`model.Xxx`
- **ctx**：首参 `ctx`，透传至 `WithContext(ctx)`；禁止 `context.Background()`
- **错误**：`fmt.Errorf("动作说明: %w", err)` 上抛

```go
func (d *UserRepo) InsertUser(ctx context.Context, user *model.User) (err error) {
    hashed, _ := hashPassword(user.Password)
    user.Password = hashed
    return d.q.User.WithContext(ctx).Create(user).Error
}
```

---

## 六、选型速查表

| 场景 | 该用的类型 |
|------|-----------|
| 绑定 HTTP 请求体 | `types.ParamXxx` |
| handler 返回单个对象 / 列表分页 | `types.Xxx` / `types.ResponseXxx` |
| service 与 repo 之间传 DB 实体 | `model.Xxx` |
| service 把 DB 结果转 API 对象 | `convert.go` 的 `modelToXxx(...)` |
| 关联表写入的入参 | `model.Xxx` 或基础类型，不收 `types.ParamXxx` |
| 统一 HTTP 响应包装 | `result.Success` / `result.Error` |
| 任意跨层调用的方法首参 | `ctx context.Context` |

---

## 七、常见错误（务必避免）

1. ❌ `result.Success(c, modelUser)`——`model` 直接返回会泄漏 gorm tag、`*string` 指针、DB 内部字段。
2. ❌ 用 `types.ParamXxx` 作函数返回值——它是输入契约，不应出现在响应里。
3. ❌ service 里手写 `types.User{UserID: strconv.Itoa(...)}` 做转换——应调用 `convert.go` 的 `modelToUser`。
4. ❌ repository 返回 `types.Xxx` 或收 `types.ParamXxx`——repo 只认 `model.Xxx`。
5. ❌ `types.Xxx` 里 ID 用 `int64`——必须是 `string`（见 `idconv.ToStr`）。
6. ❌ repo 内 `context.Background()` 或方法省略 `ctx`——`ctx` 必须从 handler 透传到底。
