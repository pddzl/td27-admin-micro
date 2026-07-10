package logic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"td27/rpc/basis/internal/svc"
	"td27/rpc/basis/types/basis_pb"
)

type PingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingLogic {
	return &PingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PingLogic) Ping(_ *basis_pb.PingReq) (*basis_pb.PingResp, error) {
	dbStatus := "reachable"
	sqlDB := l.svcCtx.DB.DB
	if err := sqlDB.Ping(); err != nil {
		dbStatus = "unreachable: " + err.Error()
	}
	return &basis_pb.PingResp{DbStatus: dbStatus}, nil
}
