package result

type ResCode int

// 错误码结构：6 位 = AA BB CC
//   AA（前两位）：业务域     11=ACL  12=Product  99=系统
//   BB（中两位）：模块       00=通用 01=用户 02=角色 03=菜单 …（系统域无模块，恒 00）
//   CC（后两位）：错误序号   从 01 起递增
// 该结构自带可观测性：看日志里的码即可定位“域/模块/错误”，无需查表。
// 注意：httpStatusForCode 仍按 code/10000 == 99 判定系统域（前两位 AA），
// 因此所有 AA=99 的码返回 HTTP 500，其余返回 200，无需改动。

// —— 成功码：全局统一，不分域 ——
const CodeSuccess ResCode = 200

// —— 系统错误域（99，全局共用，不属任何业务，BB=00）——
// 该域下的错误码会让 HTTP 状态码返回 500，触发网关熔断。
const (
	SysBase         ResCode = 990000
	CodeDBError             = SysBase + 1 // 数据库故障
	CodeRedisError          = SysBase + 2 // 缓存故障
	CodeExternalErr         = SysBase + 3 // 外部服务不可用
	CodePanic               = SysBase + 4 // 运行时恐慌
	CodeServerBusy          = SysBase + 5 // 服务繁忙（系统级，回 500）
)

// —— ACL 域（11）——
// 模块划分：00=通用  01=用户  02=角色  03=菜单
const (
	ACLBase ResCode = 110000 // 域基址

	// 通用模块（BB=00）：参数 / 登录 / 鉴权 / 路由
	ACLCommonBase    = ACLBase + 0000
	CodeInvalidParam = ACLCommonBase + 1 // 请求参数错误
	CodeNeedLogin    = ACLCommonBase + 2 // 需要登录
	CodeInvalidToken = ACLCommonBase + 3 // 无效的 Token
	CodeNoRoute      = ACLCommonBase + 4 // 请求路径不存在
	CodeNoPermission = ACLCommonBase + 5 // 无操作权限

	// 用户模块（BB=01）
	ACLUserBase         = ACLBase + 0100
	CodeUserExist       = ACLUserBase + 1 // 用户名已存在
	CodeUserDBErr       = ACLUserBase + 2 // 用户数据访问失败
	CodeInvalidPassword = ACLUserBase + 3 // 用户名或密码错误

	// 角色模块（BB=02）
	ACLRoleBase   = ACLBase + 0200
	CodeRoleDBErr = ACLRoleBase + 1 // 角色数据访问失败
	CodeRoleExist = ACLRoleBase + 2 // 角色已存在

	// 菜单模块（BB=03）
	ACLMenuBase       = ACLBase + 0300
	CodeMenuNodeExist = ACLMenuBase + 1 // 该节点下有子节点，不可以删除
	CodeMenuDBErr     = ACLMenuBase + 2 // 菜单数据访问失败
	CodeMenuExist     = ACLMenuBase + 3 // 菜单名称已存在
)

// —— Product 域（12）——
// 模块划分：00=通用  01=品牌  02=分类  03=属性  04=SPU  05=SKU
const (
	ProductBase ResCode = 120000 // 域基址

	// 通用模块（BB=00）
	ProductCommonBase = ProductBase + 0000
	// CodeProductCommonXxx 待补充

	// 品牌模块（BB=01）
	ProductBrandBase   = ProductBase + 0100
	CodeTrademarkErr   = ProductBrandBase + 1 // 品牌操作失败
	CodeTrademarkDBErr = ProductBrandBase + 2 // 品牌数据访问失败

	// 分类模块（BB=02）
	ProductCategoryBase = ProductBase + 0200
	CodeCategoryDBErr   = ProductCategoryBase + 1 // 分类数据访问失败

	// 属性模块（BB=03）
	ProductAttrBase = ProductBase + 0300
	CodeAttrDBErr   = ProductAttrBase + 1 // 属性数据访问失败

	// SPU 模块（BB=04）
	ProductSPUBase = ProductBase + 0400
	CodeSpuDBErr   = ProductSPUBase + 1 // SPU 数据访问失败

	// SKU 模块（BB=05）
	ProductSKUBase = ProductBase + 0500
	CodeSkuDBErr   = ProductSKUBase + 1 // SKU 数据访问失败
)

// —— 预留业务域（13+）——

var codeMsgMap = map[ResCode]string{
	CodeSuccess:         "success",
	CodeInvalidParam:    "请求参数错误",
	CodeNeedLogin:       "需要登录",
	CodeInvalidToken:    "无效的Token",
	CodeNoRoute:         "请求路径不存在",
	CodeNoPermission:    "无操作权限",
	CodeUserExist:       "用户名已存在",
	CodeUserDBErr:       "用户数据访问失败",
	CodeInvalidPassword: "用户名或密码错误",
	CodeRoleDBErr:       "角色数据访问失败",
	CodeRoleExist:       "角色已存在",
	CodeMenuNodeExist:   "该节点下有子节点，不可以删除",
	CodeMenuDBErr:       "菜单数据访问失败",
	CodeMenuExist:       "菜单名称已存在",
	CodeTrademarkErr:    "品牌操作失败",
	CodeTrademarkDBErr:  "品牌数据访问失败",
	CodeCategoryDBErr:   "分类数据访问失败",
	CodeAttrDBErr:       "属性数据访问失败",
	CodeSpuDBErr:        "SPU数据访问失败",
	CodeSkuDBErr:        "SKU数据访问失败",

	CodeDBError:     "数据库服务异常",
	CodeRedisError:  "缓存服务异常",
	CodeExternalErr: "外部服务不可用",
	CodePanic:       "服务内部错误",
	CodeServerBusy:  "服务繁忙",
}

func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		return "未知错误"
	}
	return msg
}

// RegisterCodeMsg 用于运行时注册自定义错误码文案。
//
// 当前不对外提供注册能力，原因：
//  1. codeMsgMap 是无锁的普通 map，仅在包初始化阶段（并发安全）填充完毕，
//     运行期之后一直是“只读”使用。若允许在运行期调用注册，会触发并发写/读写，
//     在 Go 中属于未定义行为，可能导致 panic 或数据竞态。
//  2. 若为了支持运行期注册而给 codeMsgMap 加读写锁，则每次 ResCode.Msg()
//     查询都要付出 RLock 开销。当前错误码体系是编译期常量、静态可枚举，
//     完全可以在启动前（包 init/常量定义）就写死，没有运行期注册的真实需求，
//     因此不希望为这个伪需求引入锁开销。
//
// 结论：保持静态 map，错误码及其文案在编译期确定即可。
// 若未来确有运行期注册需求（如插件化/配置驱动），再改为带锁的注册器或 sync.Map。
// func RegisterCodeMsg(code ResCode, msg string) {
// 	if _, ok := codeMsgMap[code]; ok {
// 		panic("code already registered: " + strconv.Itoa(int(code)) + " " + code.Msg())
// 	}
// 	codeMsgMap[code] = msg
// }
