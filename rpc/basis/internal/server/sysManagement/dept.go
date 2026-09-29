package sysManagement

import (
	"context"

	"td27/rpc/basis/internal/logic/sysManagement"
	"td27/rpc/basis/internal/svc"
	"td27/rpc/basis/types/common_pb"
	"td27/rpc/basis/types/sysManagement/dept_pb"
)

type DeptServer struct {
	svcCtx *svc.ServiceContext
	dept_pb.UnimplementedDeptServer
}

func NewDeptServer(svcCtx *svc.ServiceContext) *DeptServer {
	return &DeptServer{
		svcCtx: svcCtx,
	}
}

func (ds *DeptServer) GetDept(ctx context.Context, in *common_pb.IdReq) (*dept_pb.DeptResp, error) {
	dl := sysManagement.NewDeptLogic(ctx, ds.svcCtx)
	return dl.GetDept(in)
}

func (ds *DeptServer) GetDeptTree(ctx context.Context, in *common_pb.Empty) (*dept_pb.GetDeptTreeResp, error) {
	dl := sysManagement.NewDeptLogic(ctx, ds.svcCtx)
	return dl.GetDeptTree(in)
}

func (ds *DeptServer) GetDeptDescendants(ctx context.Context, in *dept_pb.GetDeptDescendantsReq) (*dept_pb.GetDeptTreeResp, error) {
	dl := sysManagement.NewDeptLogic(ctx, ds.svcCtx)
	return dl.GetDeptDescendants(in)
}

func (ds *DeptServer) CreateDept(ctx context.Context, in *dept_pb.CreateDeptReq) (*common_pb.SuccessResp, error) {
	dl := sysManagement.NewDeptLogic(ctx, ds.svcCtx)
	return dl.CreateDept(in)
}

func (ds *DeptServer) UpdateDept(ctx context.Context, in *dept_pb.UpdateDeptReq) (*dept_pb.DeptResp, error) {
	dl := sysManagement.NewDeptLogic(ctx, ds.svcCtx)
	return dl.UpdateDept(in)
}

func (ds *DeptServer) ListDept(ctx context.Context, in *dept_pb.ListDeptReq) (*dept_pb.ListDeptResp, error) {
	dl := sysManagement.NewDeptLogic(ctx, ds.svcCtx)
	return dl.ListDept(in)
}

func (ds *DeptServer) DeleteDept(ctx context.Context, in *common_pb.IdReq) (*common_pb.SuccessResp, error) {
	dl := sysManagement.NewDeptLogic(ctx, ds.svcCtx)
	return dl.DeleteDept(in)
}
