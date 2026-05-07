package tool

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	toolModel "td27/rpc/basis/internal/model/tool"
	"td27/rpc/basis/internal/model/common"
	"td27/rpc/basis/internal/svc"
	"td27/rpc/basis/internal/util"
	"td27/rpc/basis/types/tool/service_token_pb"
	"td27/rpc/basis/types/common_pb"
)

type ServiceTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewServiceTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ServiceTokenLogic {
	return &ServiceTokenLogic{
		ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx),
	}
}

func (l *ServiceTokenLogic) mapToResp(token *toolModel.ServiceToken) *service_token_pb.ServiceTokenResp {
	if token == nil { return nil }
	resp := &service_token_pb.ServiceTokenResp{
		Id: uint64(token.ID), Name: token.Name, Status: token.Status,
		CreatedAt: util.ToProtoTimestamp(token.CreatedAt), UpdatedAt: util.ToProtoTimestamp(token.UpdatedAt),
	}
	if token.ExpiresAt != nil { resp.ExpiresAt = token.ExpiresAt }
	return resp
}

func (l *ServiceTokenLogic) CreateServiceToken(in *service_token_pb.CreateServiceTokenReq) (*service_token_pb.CreateServiceTokenResp, error) {
	if in.Name == "" { return nil, status.Errorf(codes.InvalidArgument, "token name is required") }
	token := &toolModel.ServiceToken{Name: in.Name, Status: true}
	if in.Status != nil { token.Status = *in.Status }
	var ttl *time.Duration
	if in.ExpiresAt != nil { d := time.Duration(*in.ExpiresAt) * time.Second; ttl = &d }
	rawToken, err := l.svcCtx.TokenService.Create(l.ctx, token, ttl)
	if err != nil { return nil, status.Errorf(codes.Internal, "create token failed: %v", err) }
	return &service_token_pb.CreateServiceTokenResp{TokenId: uint64(token.ID), RawToken: rawToken}, nil
}

func (l *ServiceTokenLogic) GetServiceToken(in *common_pb.IdReq) (*service_token_pb.ServiceTokenResp, error) {
	if in.Id == 0 { return nil, status.Errorf(codes.InvalidArgument, "invalid token id") }
	token, err := l.svcCtx.TokenService.GetByID(l.ctx, uint(in.Id))
	if err != nil { return nil, status.Errorf(codes.Internal, "get token failed: %v", err) }
	if token == nil { return nil, status.Errorf(codes.NotFound, "token not found") }
	return l.mapToResp(token), nil
}

func (l *ServiceTokenLogic) ListServiceToken(in *common_pb.PageReq) (*service_token_pb.ListServiceTokenResp, error) {
	page := &common.PageInfo{Page: int(in.Page), PageSize: int(in.PageSize)}
	tokens, count, err := l.svcCtx.TokenService.List(l.ctx, page)
	if err != nil { return nil, status.Errorf(codes.Internal, "list tokens failed: %v", err) }
	resp := &service_token_pb.ListServiceTokenResp{List: make([]*service_token_pb.ServiceTokenResp, 0, len(tokens)), Total: count}
	for _, t := range tokens { resp.List = append(resp.List, l.mapToResp(t)) }
	return resp, nil
}

func (l *ServiceTokenLogic) UpdateServiceToken(in *service_token_pb.UpdateServiceTokenReq) (*service_token_pb.ServiceTokenResp, error) {
	if in.Id == 0 { return nil, status.Errorf(codes.InvalidArgument, "invalid token id") }
	token, err := l.svcCtx.TokenService.GetByID(l.ctx, uint(in.Id))
	if err != nil { return nil, status.Errorf(codes.Internal, "get token failed: %v", err) }
	if token == nil { return nil, status.Errorf(codes.NotFound, "token not found") }
	if in.Name != nil { token.Name = *in.Name }
	if in.Status != nil { token.Status = *in.Status }
	if in.ExpiresAt != nil { token.ExpiresAt = in.ExpiresAt }
	if err := l.svcCtx.TokenService.Update(l.ctx, token); err != nil {
		return nil, status.Errorf(codes.Internal, "update token failed: %v", err)
	}
	return l.mapToResp(token), nil
}

func (l *ServiceTokenLogic) ToggleTokenStatus(in *service_token_pb.ToggleTokenStatusReq) (*common_pb.SuccessResp, error) {
	if in.Id == 0 { return nil, status.Errorf(codes.InvalidArgument, "invalid token id") }
	if err := l.svcCtx.TokenService.ToggleStatus(l.ctx, uint(in.Id), in.Status); err != nil {
		return nil, status.Errorf(codes.Internal, "toggle token status failed: %v", err)
	}
	return &common_pb.SuccessResp{Success: true}, nil
}

func (l *ServiceTokenLogic) AssignTokenPermissions(in *service_token_pb.AssignTokenPermissionsReq) (*common_pb.SuccessResp, error) {
	if in.TokenId == 0 { return nil, status.Errorf(codes.InvalidArgument, "invalid token id") }
	permIDs := make([]uint, len(in.PermissionIds))
	for i, id := range in.PermissionIds { permIDs[i] = uint(id) }
	if err := l.svcCtx.TokenService.AssignPermissions(l.ctx, uint(in.TokenId), permIDs); err != nil {
		return nil, status.Errorf(codes.Internal, "assign permissions failed: %v", err)
	}
	return &common_pb.SuccessResp{Success: true}, nil
}

func (l *ServiceTokenLogic) GetTokenPermissions(in *common_pb.IdReq) (*service_token_pb.GetTokenPermissionsResp, error) {
	if in.Id == 0 { return nil, status.Errorf(codes.InvalidArgument, "invalid token id") }
	permIDs, err := l.svcCtx.TokenService.GetPermissions(l.ctx, uint(in.Id))
	if err != nil { return nil, status.Errorf(codes.Internal, "get permissions failed: %v", err) }
	resp := &service_token_pb.GetTokenPermissionsResp{PermissionIds: make([]uint64, len(permIDs))}
	for i, id := range permIDs { resp.PermissionIds[i] = uint64(id) }
	return resp, nil
}

func (l *ServiceTokenLogic) DeleteServiceToken(in *common_pb.IdReq) (*common_pb.SuccessResp, error) {
	if in.Id == 0 { return nil, status.Errorf(codes.InvalidArgument, "invalid token id") }
	if err := l.svcCtx.TokenService.Delete(l.ctx, uint(in.Id)); err != nil {
		return nil, status.Errorf(codes.Internal, "delete token failed: %v", err)
	}
	return &common_pb.SuccessResp{Success: true}, nil
}

func (l *ServiceTokenLogic) ValidateToken(in *service_token_pb.ValidateTokenReq) (*service_token_pb.ValidateTokenResp, error) {
	if in.Token == "" { return nil, status.Errorf(codes.InvalidArgument, "token is required") }
	token, err := l.svcCtx.TokenService.GetByToken(l.ctx, in.Token)
	if err != nil { return nil, status.Errorf(codes.Internal, "validate token failed: %v", err) }
	resp := &service_token_pb.ValidateTokenResp{Valid: false}
	if token != nil && token.Status {
		valid := true
		if token.ExpiresAt != nil && time.Now().Unix() >= *token.ExpiresAt { valid = false }
		if valid { resp.Valid = true; resp.Token = l.mapToResp(token) }
	}
	return resp, nil
}
