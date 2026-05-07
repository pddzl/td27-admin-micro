package tool

import (
	"context"
	"td27/rpc/basis/internal/logic/tool"
	"td27/rpc/basis/internal/svc"
	"td27/rpc/basis/types/tool/cron_pb"
	"td27/rpc/basis/types/common_pb"
)

type CronServer struct {
	svcCtx *svc.ServiceContext
	cron_pb.UnimplementedCronServer
}

func NewCronServer(svcCtx *svc.ServiceContext) *CronServer {
	return &CronServer{svcCtx: svcCtx}
}

func (s *CronServer) GetCron(ctx context.Context, in *common_pb.IdReq) (*cron_pb.CronResp, error) {
	return tool.NewCronLogic(ctx, s.svcCtx).GetCron(in)
}
func (s *CronServer) ListCron(ctx context.Context, in *common_pb.PageReq) (*cron_pb.ListCronResp, error) {
	return tool.NewCronLogic(ctx, s.svcCtx).ListCron(in)
}
func (s *CronServer) CreateCron(ctx context.Context, in *cron_pb.CreateCronReq) (*common_pb.SuccessResp, error) {
	return tool.NewCronLogic(ctx, s.svcCtx).CreateCron(in)
}
func (s *CronServer) UpdateCron(ctx context.Context, in *cron_pb.UpdateCronReq) (*cron_pb.CronResp, error) {
	return tool.NewCronLogic(ctx, s.svcCtx).UpdateCron(in)
}
func (s *CronServer) ToggleCronStatus(ctx context.Context, in *cron_pb.ToggleCronStatusReq) (*common_pb.SuccessResp, error) {
	return tool.NewCronLogic(ctx, s.svcCtx).ToggleCronStatus(in)
}
func (s *CronServer) ExecuteCronNow(ctx context.Context, in *cron_pb.ExecuteCronNowReq) (*common_pb.SuccessResp, error) {
	return tool.NewCronLogic(ctx, s.svcCtx).ExecuteCronNow(in)
}
func (s *CronServer) DeleteCron(ctx context.Context, in *common_pb.IdReq) (*common_pb.SuccessResp, error) {
	return tool.NewCronLogic(ctx, s.svcCtx).DeleteCron(in)
}
