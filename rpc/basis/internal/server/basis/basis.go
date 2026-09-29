package basis

import (
	"context"

	"td27/rpc/basis/internal/logic"
	"td27/rpc/basis/internal/svc"
	"td27/rpc/basis/types/basis_pb"
)

type BasisServer struct {
	svcCtx *svc.ServiceContext
	basis_pb.UnimplementedBasisServer
}

func NewBasisServer(svcCtx *svc.ServiceContext) *BasisServer {
	return &BasisServer{
		svcCtx: svcCtx,
	}
}

func (s *BasisServer) Ping(ctx context.Context, req *basis_pb.PingReq) (*basis_pb.PingResp, error) {
	l := logic.NewPingLogic(ctx, s.svcCtx)
	return l.Ping(req)
}
