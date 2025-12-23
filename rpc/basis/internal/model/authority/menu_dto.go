package authority

type MenuDTO struct {
	Pid       uint       `json:"pid"`                          // 默认0 根目录
	Name      string     `json:"name"`                         // 名称
	Path      string     `json:"path" binding:"required"`      // 路径
	Redirect  string     `json:"redirect"`                     // 重定向
	Component string     `json:"component" binding:"required"` // 前端组件
	Sort      uint       `json:"sort" binding:"required"`      // 排序
	Meta      MetaEntity `json:"meta"`
}

type NewMenuDTO struct {
	ID uint
	MenuDTO
}
