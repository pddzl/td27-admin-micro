package monitor

import (
	"td27/rpc/basis/internal/model/common"
)

// OperationLogModel System operation audit log entity
type OperationLogModel struct {
	common.Td27Model
	Ip        string `json:"ip" gorm:"comment:请求ip"`
	Method    string `json:"method" gorm:"comment:请求方法"`
	Path      string `json:"path" gorm:"comment:请求路径"`
	Status    int    `json:"status" gorm:"comment:请求状态"`
	UserAgent string `json:"userAgent"`
	ReqParam  string `json:"reqParam" gorm:"type:text;comment:请求Body"`
	RespData  string `json:"respData" gorm:"type:text;comment:响应数据"`
	RespTime  int64  `json:"respTime"`
	UserID    uint   `json:"userID" gorm:"comment:用户id"`
	UserName  string `json:"userName" gorm:"comment:用户名称"`
}

func (ol *OperationLogModel) TableName() string {
	return "sys_monitor_operation_log"
}
