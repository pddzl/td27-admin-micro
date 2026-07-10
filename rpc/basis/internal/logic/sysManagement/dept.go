package sysManagement

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"td27/rpc/basis/internal/model/common"
	"td27/rpc/basis/internal/model/sysManagement"
	"td27/rpc/basis/internal/svc"
	"td27/rpc/basis/internal/util"
	"td27/rpc/basis/types/common_pb"
	"td27/rpc/basis/types/sysManagement/dept_pb"
)

type DeptLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeptLogic {
	return &DeptLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (dl *DeptLogic) mapDeptToResp(dept *sysManagement.DeptModel) *dept_pb.DeptResp {
	if dept == nil {
		return nil
	}

	return &dept_pb.DeptResp{
		Id:        int64(dept.ID),
		DeptName:  dept.DeptName,
		ParentId:  int64(dept.ParentID),
		Path:      dept.Path,
		Level:     uint32(dept.Level),
		Sort:      uint32(dept.Sort),
		Status:    dept.Status,
		CreatedAt: util.ToProtoTimestamp(dept.CreatedAt),
		UpdatedAt: util.ToProtoTimestamp(dept.UpdatedAt),
	}
}

func (dl *DeptLogic) mapDeptToTree(dept *sysManagement.DeptModel, allDepts []*sysManagement.DeptModel) *dept_pb.DeptTreeResp {
	resp := &dept_pb.DeptTreeResp{
		Dept:     dl.mapDeptToResp(dept),
		Children: make([]*dept_pb.DeptTreeResp, 0),
	}

	for _, child := range allDepts {
		if child.ParentID == dept.ID {
			resp.Children = append(resp.Children, dl.mapDeptToTree(child, allDepts))
		}
	}

	return resp
}

func (dl *DeptLogic) GetDept(in *common_pb.IdReq) (*dept_pb.DeptResp, error) {
	dept, err := dl.svcCtx.DeptService.GetByID(dl.ctx, uint(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get department failed: %v", err)
	}
	if dept == nil {
		return nil, status.Errorf(codes.NotFound, "department not found")
	}

	return dl.mapDeptToResp(dept), nil
}

func (dl *DeptLogic) GetDeptTree(in *common_pb.Empty) (*dept_pb.GetDeptTreeResp, error) {
	allDepts, err := dl.svcCtx.DeptService.GetAll(dl.ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get all departments failed: %v", err)
	}

	tree := make([]*dept_pb.DeptTreeResp, 0)
	for _, dept := range allDepts {
		if dept.ParentID == 0 {
			tree = append(tree, dl.mapDeptToTree(dept, allDepts))
		}
	}

	return &dept_pb.GetDeptTreeResp{
		Tree: tree,
	}, nil
}

func (dl *DeptLogic) GetDeptDescendants(in *dept_pb.GetDeptDescendantsReq) (*dept_pb.GetDeptTreeResp, error) {
	dept, err := dl.svcCtx.DeptService.GetByID(dl.ctx, uint(in.DeptId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get department failed: %v", err)
	}
	if dept == nil {
		return nil, status.Errorf(codes.NotFound, "department not found")
	}

	descendants, err := dl.svcCtx.DeptService.GetDescendants(dl.ctx, uint(in.DeptId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get department descendants failed: %v", err)
	}

	tree := make([]*dept_pb.DeptTreeResp, 0)
	tree = append(tree, dl.mapDeptToTree(dept, descendants))

	return &dept_pb.GetDeptTreeResp{
		Tree: tree,
	}, nil
}

func (dl *DeptLogic) CreateDept(in *dept_pb.CreateDeptReq) (*common_pb.SuccessResp, error) {
	parentID := uint(*in.ParentId)
	dept := &sysManagement.DeptModel{
		DeptName: in.DeptName,
		ParentID: parentID,
		Sort:     uint(*in.Sort),
		Status:   *in.Status,
	}

	err := dl.svcCtx.DeptService.Create(dl.ctx, dept)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create department failed: %v", err)
	}

	return &common_pb.SuccessResp{Success: true}, nil
}

func (dl *DeptLogic) UpdateDept(in *dept_pb.UpdateDeptReq) (*dept_pb.DeptResp, error) {
	dept, err := dl.svcCtx.DeptService.GetByID(dl.ctx, uint(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get department failed: %v", err)
	}
	if dept == nil {
		return nil, status.Errorf(codes.NotFound, "department not found")
	}

	if in.DeptName != nil {
		dept.DeptName = *in.DeptName
	}
	if in.ParentId != nil {
		dept.ParentID = uint(*in.ParentId)
	}
	if in.Sort != nil {
		dept.Sort = uint(*in.Sort)
	}
	if in.Status != nil {
		dept.Status = *in.Status
	}

	err = dl.svcCtx.DeptService.Update(dl.ctx, dept)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update department failed: %v", err)
	}

	updatedDept, err := dl.svcCtx.DeptService.GetByID(dl.ctx, uint(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get updated department failed: %v", err)
	}

	return dl.mapDeptToResp(updatedDept), nil
}

func (dl *DeptLogic) DeleteDept(in *common_pb.IdReq) (*common_pb.SuccessResp, error) {
	if in.Id == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid department id")
	}

	err := dl.svcCtx.DeptService.Delete(dl.ctx, uint(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "delete department failed: %v", err)
	}

	return &common_pb.SuccessResp{Success: true}, nil
}

func (dl *DeptLogic) ListDept(in *dept_pb.ListDeptReq) (*dept_pb.ListDeptResp, error) {
	page := &common.PageInfo{
		Page:     int(in.Page),
		PageSize: int(in.PageSize),
	}

	depts, count, err := dl.svcCtx.DeptService.List(dl.ctx, page)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ListDept failed: %v", err)
	}

	resp := &dept_pb.ListDeptResp{
		List:  make([]*dept_pb.DeptResp, 0, len(depts)),
		Total: count,
	}

	for _, dept := range depts {
		resp.List = append(resp.List, dl.mapDeptToResp(dept))
	}

	return resp, nil
}
