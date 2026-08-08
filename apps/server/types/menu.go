package types

// 菜单类型（对应 menu.type）：1=目录，2=菜单/按钮，0=其它
const (
	MenuTypeOther int = 0
	MenuTypeDir   int = 1
	MenuTypeMenu  int = 2
)

// 菜单层级（对应 menu.level）：1=一级，2=二级，3=三级，4=按钮
const (
	MenuLevel1 int = 1
	MenuLevel2 int = 2
	MenuLevel3 int = 3
	MenuLevel4 int = 4
)

type Menu struct {
	BaseModel
	MenuID   string `json:"menuId"`
	Name     string `json:"name"`
	ParentID string `json:"parentId"`
	CODE     string `json:"code"`
	TOCODE   string `json:"toCode"`
	TYPE     int    `json:"type"` // 菜单类型：1=目录 2=菜单/按钮 0=其它（见 MenuType* 常量）
	STATUS   bool   `json:"status"`
	LEVEL    int    `json:"level"` // 菜单层级：1=一级 2=二级 3=三级 4=按钮（见 MenuLevel* 常量）
	CHILDREN []Menu `json:"children"`
	SELECT   bool   `json:"select"`
}

type ParamMenuSave struct {
	Name     string `json:"name"`
	ParentID string `json:"parentId"`
	CODE     string `json:"code"`
	TYPE     int    `json:"type"`  // 菜单类型：1=目录 2=菜单/按钮 0=其它
	LEVEL    int    `json:"level"` // 菜单层级：1=一级 2=二级 3=三级 4=按钮
}

type ParamMenuUpdate struct {
	MenuID   string `json:"menuId"`
	Name     string `json:"name"`
	ParentID string `json:"parentId"`
	CODE     string `json:"code"`
	LEVEL    int    `json:"level"` // 菜单层级：1=一级 2=二级 3=三级 4=按钮
}
