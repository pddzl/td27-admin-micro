package sysMonitor

import (
	"td27/rpc/basis/internal/model/common"
)

// OperationLogModel System operation audit log entity
type OperationLogModel struct {
	common.Td27Model
	Ip        string `json:"ip" db:"ip"`
	Method    string `json:"method" db:"method"`
	Path      string `json:"path" db:"path"`
	Status    int    `json:"status" db:"status"`
	UserAgent string `json:"userAgent" db:"user_agent"`
	ReqParam  string `json:"reqParam" db:"req_param"`
	RespData  string `json:"respData" db:"resp_data"`
	RespTime  int64  `json:"respTime" db:"resp_time"`
	UserID    uint   `json:"userID" db:"user_id"`
	UserName  string `json:"userName" db:"user_name"`
}

func (ol *OperationLogModel) TableName() string {
	return "sys_monitor_operation_log"
}
