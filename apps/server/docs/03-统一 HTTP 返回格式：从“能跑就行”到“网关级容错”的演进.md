# 统一 HTTP 返回格式：从"能跑就行"到"网关级容错"的演进

在单体项目初期，很多人习惯"先把接口跑通"。但随着业务变复杂，前后端联调时往往会爆发一场血案："你这返回的格式一会儿是 JSON，一会儿是纯文本报错，我的前端拦截器怎么写？"

今天，我们就来彻底解决接口返回格式的问题。但我要提醒你，**这不仅仅是怎么封装一个
JSON 的事，这涉及到网关容错与雪崩防御。**

本项目使用 **Gin** 作为 Web 框架，统一返回封装在 `pkg/result`
包中。我们将经历四次迭代，带你感受架构演进的魅力。

---

## V1：野蛮生长（先跑通再说）

在没有任何规范约束下，我们快速写出了一个登录接口。为了让前端知道结果，很多人会这么写 handler：

```go
func (h *UserHandler) Login(c *gin.Context) {
    var req LoginReq
    if err := c.ShouldBindJSON(&req); err != nil {
        // 痛点：把通讯协议和业务数据强耦合，且错误格式随意
        c.JSON(200, gin.H{"error": err.Error()})
        return
    }
    user, err := h.svc.Login(req)
    if err != nil {
        // 直接抛出原生错误，前端根本没法统一处理
        c.String(500, "密码错误")
        return
    }
    c.JSON(200, gin.H{"token": user.Token})
}
```

### 灾难预警

1. **写死个人**：每个 handler 都要手写 `c.JSON`、拼 `code`/`msg`/`data`，Gopher 会疯掉。
2. **格式失控**：遇到底层报错（比如数据库连不上）直接 `c.String(500, ...)`
   时，前端收到纯文本，拦截器直接崩溃。
3. **牵一发而动全身**：要改错误码规范，得把所有 handler 改一遍。

这，就是典型的"技术债"。

---

## V2：统一封装（引入 result 包）

在分层架构下，我们有一个铁律：**业务逻辑只负责"定罪"（返回错误码），基础设施（`pkg/result`）负责"宣判"（决定如何呈现错误）。**

我们在 `pkg/result` 中定义了统一的响应结构与错误码体系，彻底解耦 handler 与 JSON 格式。

### 第一步：定义统一的响应结构与错误码 (`pkg/result/response.go`)

```go
type ResponseData struct {
    Code ResCode `json:"code"`
    Msg  string  `json:"message"`
    Data any     `json:"data"`
}

// httpStatusForCode 根据码所属域决定 HTTP 状态码：
// 99 域（系统错误）→ 500，其余（业务错误）→ 200。
func httpStatusForCode(code ResCode) int {
    if code/10000 == 99 {
        return http.StatusInternalServerError
    }
    return http.StatusOK
}

func writeJSON(c *gin.Context, code ResCode, msg string, data any) {
    c.JSON(httpStatusForCode(code), &ResponseData{
        Code: code,
        Msg:  msg,
        Data: data,
    })
}

func Error(c *gin.Context, code ResCode) {
    writeJSON(c, code, code.Msg(), nil)
}

func Success(c *gin.Context, data any) {
    writeJSON(c, CodeSuccess, CodeSuccess.Msg(), data)
}

func ErrorWithMsg(c *gin.Context, code ResCode, msg string) {
    writeJSON(c, code, msg, nil)
}
```

### 第二步：定义 6 位错误码与静态文案 (`pkg/result/code.go`)

错误码结构：`6 位 = AA BB CC`

- `AA`（前两位）：业务域 `11=ACL`、`12=Product`、`99=系统`
- `BB`（中两位）：模块 `00=通用 01=用户 02=角色 03=菜单 …`
- `CC`（后两位）：错误序号











```go
type ResCode int

const CodeSuccess ResCode = 200

// 系统错误域（99）：返回 HTTP 500，触发网关熔断
const (
    SysBase         ResCode = 990000
    CodeDBError             = SysBase + 1 // 数据库故障
    CodeRedisError          = SysBase + 2 // 缓存故障
    CodeExternalErr         = SysBase + 3 // 外部服务不可用
    CodePanic               = SysBase + 4 // 运行时恐慌
    CodeServerBusy          = SysBase + 5 // 服务繁忙
)

// ACL 域（11）
const (
    ACLBase              ResCode = 110000
    CodeInvalidParam             = ACLBase + 0001
    CodeNeedLogin                = ACLBase + 0002
    CodeInvalidToken             = ACLBase + 0003
    CodeNoRoute                  = ACLBase + 0004
    CodeUserExist                = ACLBase + 0101
    CodeUserNotExist             = ACLBase + 0102
    CodeInvalidPassword          = ACLBase + 0103
    CodeMenuNodeExist            = ACLBase + 0301
)

// Product 域（12）
const (
    ProductBase     ResCode = 120000
    CodeTrademarkErr         = ProductBase + 0101
)

// 静态文案表：无锁、启动期填充完毕、运行期只读
var codeMsgMap = map[ResCode]string{
    CodeSuccess:         "success",
    CodeInvalidParam:    "请求参数错误",
    CodeNeedLogin:       "需要登录",
    CodeInvalidToken:    "无效的Token",
    CodeNoRoute:         "请求路径不存在",
    CodeUserExist:       "用户名已存在",
    CodeUserNotExist:    "用户名不存在",
    CodeInvalidPassword: "用户名或密码错误",
    CodeMenuNodeExist:   "该节点下有子节点，不可以删除",
    CodeTrademarkErr:    "品牌操作失败",

    CodeDBError:     "数据库服务异常",
    CodeRedisError:  "缓存服务异常",
    CodeExternalErr: "外部服务不可用",
    CodePanic:       "服务内部错误",
    CodeServerBusy:  "服务繁忙",
}

func (c ResCode) Msg() string {
    if msg, ok := codeMsgMap[c]; ok {
        return msg
    }
    return "未知错误"
}
```

#### 重构后的 handler 层

```go
func (h *UserHandler) Login(c *gin.Context) {
    var req LoginReq
    if err := c.ShouldBindJSON(&req); err != nil {
        // 没有任何 JSON 格式化的杂糅代码！
        result.Error(c, result.CodeInvalidParam)
        return
    }
    if err := h.svc.Login(req); err != nil {
        result.Error(c, result.CodeInvalidPassword)
        return
    }
    result.Success(c, gin.H{"token": user.Token})
}
```

前端现在收到的永远是整齐划一的 JSON：

```json
// 成功
{ "code": 200, "message": "success", "data": { "token": "xxx" } }
// 密码错误
{ "code": 110103, "message": "用户名或密码错误", "data": null }
```

### 💡 设计意图：为什么 `httpStatusForCode` 要按 99 域分级

V2 看似只统一了 JSON 格式，但 `httpStatusForCode` 里 `code/10000 == 99`
这一行约定，其实决定了网关容错能力：**系统错误（99 域）返回 500，业务错误返回 200。**

前端同学确实高兴了，但如果你前面挂了一个 **Nginx 作为网关**，Nginx 会依据 HTTP
状态码判断节点健康。

**思考一个场景**：你的 PostgreSQL 进程挂了（或者 OOM 了）。

按照 V2 的逻辑，`CodeDBError`（990001）属于 99 域，`httpStatusForCode` 会返回 `HTTP 500 {"code": 990001,
"message": "数据库服务异常"}`。

Nginx 会怎么想？Nginx 会认为："哦，这个节点返回 500，病了，快熔断摘除！" 于是 Nginx
把流量切到健康节点，**避免雪崩**。

这正是我们设计 6 位错误码 `AA=99` 的用意：**HTTP 状态码是给基础设施（网关）看的，JSON
里的 Code 是给前端业务逻辑看的。** 两者通过 `code/10000 == 99`
这一行约定自动联动，无需任何额外配置。

看起来很完美了对吧？**别急，真正的坑才刚刚开始。**

---

## V3：真实世界的雷区：瞬态故障与三层防御体系

到了
V2，我们的代码已经非常健壮了。但在生产环境，还有一个极其狡猾的幽灵：**瞬态故障**。

比如：运维同学做 PostgreSQL 主从切换，会有短暂的 3 秒钟连接断开。

这 3 秒内进来的请求，GORM 底层报错，被我们的 handler 捕获，调用 `result.Error(c,
result.CodeDBError)`，**毫无疑问会返回 HTTP 500**。

问题来了：**用户辛辛苦苦填了个表单，就因为 PG
抖了一下，你就给他弹个"系统开小差"，这体验能忍吗？**

要完美处理这个问题，很多新手会在 Go
代码里写一堆复杂的"重试逻辑"。错！在微服务架构下，**不要让业务代码去承担基础架构的容错责任**。我们需要建立"三层防线"：

### 第一层防线：Go 底层连接池（吸收微秒级抖动）

其实，对于几百毫秒的网络闪断，你的 Go 代码根本感知不到。Go 标准库的 `database/sql`（GORM
底层基于此）内置了连接池健康检查机制。如果某个连接断了，底层会自动丢弃并重连。只有当连接池里的连接全断了，且超过了配置的超时时间，错误才会真正抛到你的 handler 里。

#### 结论：微秒级抖动，Go 自己消化了，不需要我们管

### 第二层防线：Nginx 网关（解决秒级切换与请求路由）

如果是长达 3 秒的 PG 主从切换，Go 代码确实会返回 500。这时候，**Nginx 必须出场**。

我们在 Nginx 配置里这样写：

```nginx
proxy_next_upstream error timeout http_500 http_502 http_503 http_504;
```

**对于 GET 请求（查询列表）**：节点 A 返回 500，Nginx 瞬间把请求透明地转发给节点
B。如果节点 B 的数据库连接正常，用户毫无察觉，体验完美！

**但对于 POST 请求（提交订单）**：Nginx 默认是极其保守的，**它绝对不会重试
POST**。为什么？想象一个极其隐蔽的灾难场景：后端其实已经扣了钱，数据落盘了，但在返回响应时发生了微秒级网络抖动，TCP 断了。Nginx 没收到响应判定失败，如果它盲目重试 POST，就会导致**重复扣款**（这叫非幂等灾难）。

#### 结论：GET 请求被 Nginx 拯救了，但 POST 请求如果撞上 PG 切换，Nginx 不敢救

### 第三层防线：业务幂等 Token（斩断最后的心魔）

既然 Nginx 不敢重试
POST，遇到这种网络抖动，难道只能让用户看到"支付失败"，然后逼着用户手动再点一次吗？如果用户狂点，依然会导致重复业务！

在微服务架构下有一个铁律：**所有核心写操作（POST），后端都必须当作"可能被重复提交"来处理。**

通用解法是**幂等 Token 机制**：

1. 用户打开"支付"页面前，先 GET 请求获取一个全局唯一 Token。
2. 提交 POST 支付时，带上这个 Token。
3. 后端收到请求，先去 Redis 查 Token：

   - **不存在**：放行扣款，成功后把 Token 存入 Redis。
   - **已存在**：说明这是重试请求（不管是网络抖动导致的重新点击，还是用户手抖狂点），**直接返回之前的成功结果，绝不重复扣款！**

### 🌟 小结：架构师的分层智慧

回过头看，面对"PG
主从切换导致失败"这同一个问题，不同水平的程序员有完全不同的解法：

- **初级**：在 Service 层写个 `for` 循环重试 3
  次。结果雪上加霜，把本来就紧张的数据库连接池彻底耗尽。

- **中级**：统一返回了 500，交给了 Nginx。结果只解决了 GET，POST 照样炸。

- **高级（我们现在的方案）**：各司其职。连接池防微抖，Nginx 防 GET 故障，业务 Token 防

  POST 重复。

只有把这三层防线组合起来，这个"瞬态故障"的雷区，才算真正被我们踩平了。而这一切优雅体验的基石，正是我们在 V2 和 V3 中写下的那个**通过 `code/10000 == 99` 精准区分 200 和 500 的 `result` 包**。

---

## V4：对外 API 的脱敏（ErrorExternal）

内部前端拿到 `"数据库服务异常"` 没问题，但如果这个项目要对外开放
API（第三方开发者），直出 `数据库/缓存` 字样会泄露内部技术栈。

我们在 `pkg/result` 中单独提供了 `ErrorExternal`，对 99 域系统错误脱敏：

```go
const externalMaskedMsg = "系统繁忙，请稍后再试"

// ErrorExternal 对外 API 专用：99 域系统错误统一脱敏，不暴露技术栈
func ErrorExternal(c *gin.Context, code ResCode) {
    msg := code.Msg()
    if code/10000 == 99 {
        msg = externalMaskedMsg
    }
    writeJSON(c, code, msg, nil)
}
```

要点：

1. **统一降级文案**：对外不暴露"数据库/缓存"字样，全部归一成"系统繁忙/依赖服务不可用"。技术栈（DB、Redis、云厂商）对攻击者是无用情报的泄露。
2. **保留 HTTP 状态码分级**：对外同样区分 200（业务错误）和
   500（系统错误），第三方据此决定退避重试，网关据此熔断。
3. **具体细节只进日志**：原始

   error（含连接串、SQL、堆栈）只打服务端日志，绝不进响应体。

用法：对外路由里用 `result.ErrorExternal(c, code)` 替代 `result.Error(c, code)` 即可。`codeMsgMap`
与内部前端通道完全不动。

---

## 💡 深度扩展：业务代码里到底要不要写重试

很多同学学完了上面的三层防线，依然会有疑问："既然 PG 主从切换会导致失败，为什么在
Go 代码里绝对不能写重试？"

答案取决于你面对的**对象**：

1. **查数据库（无状态读）**：不需要写。Go 的 `database/sql`
   底层连接池已经自带了"坏连接丢弃与重试"机制。它对业务层是透明的。
2. **改数据库（事务写）**：**严禁写简单重试！**
   如果事务执行到一半网络断了，应用层处于"不可知状态"（不知道刚才的 SQL
   有没有落盘）。盲目重试会导致数据重复修改（比如重复扣款）。事务中途失败的修复，只能依赖下游的"幂等设计"或离线"对账系统"。
3. **调外部 API（如微信支付）**：**必须写重试！**

   外部接口超时，我们不知道对方是否成功，必须结合上一节讲的"幂等
   Token"，在代码里进行有限次数的重试，以保证最终一致性。

> ### 总监视角的延伸（MQ 消费者场景）

>
> 如果你是在 NATS/RabbitMQ 的消费者里写逻辑，用 `pgx.BeginFunc` 包裹事务。遇到 PG
> 切换报错，**连重试的念头都不要有，直接 return error，坚决不 ACK**。让 NATS
> 去负责重试投递！你的代码只需要保证：只要执行，就结合 `临时表 + ON CONFLICT DO NOTHING`
> 保证幂等。**把重试的责任，甩给最专业的组件。**

### 记住一句话：内部状态变更靠幂等兜底，外部网络调用靠重试保底。各司其职，方为架构

---

## 总结

从 V1 到 V4，我们把一个简单的"统一返回格式"，挖出了五层架构认知：

1. **解耦**：通过 `pkg/result` 包，让 handler 层回归纯粹的业务逻辑，只返回 `ResCode`。
2. **分级**（V2 设计意图）：通过 6 位错误码 `AA=99` 与 `code/10000 == 99` 的精准划分，与 API
   网关达成容错契约（200/500）。
3. **容错**（V3）：理解瞬态故障，善用 Nginx 的重试机制。
4. **防重**（V3）：通过幂等 Token，彻底消灭非幂等请求在分布式环境下的幽灵 Bug。
5. **脱敏**（V4）：对外 API 用 `ErrorExternal` 隔离系统错误细节，保护内部技术栈。

最后，请再看一眼我们精心设计的 `pkg/result` 目录：

- `code.go`：6 位错误码常量 + 无锁静态文案表 `codeMsgMap`。
- `response.go`：`Success` / `Error` / `ErrorWithMsg` / `ErrorExternal` 与 HTTP 状态码分级逻辑。
