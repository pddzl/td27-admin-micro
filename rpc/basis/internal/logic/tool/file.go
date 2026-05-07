package tool

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	toolModel "td27/rpc/basis/internal/model/tool"
	"td27/rpc/basis/internal/model/common"
	"td27/rpc/basis/internal/svc"
	"td27/rpc/basis/internal/util"
	"td27/rpc/basis/types/tool/file_pb"
	"td27/rpc/basis/types/common_pb"
)

type FileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileLogic {
	return &FileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (fl *FileLogic) mapFileToResp(file *toolModel.FileModel) *file_pb.FileResp {
	if file == nil {
		return nil
	}

	return &file_pb.FileResp{
		Id:        uint64(file.ID),
		FileName:  file.FileName,
		FullPath:  file.FullPath,
		Mime:      file.Mime,
		CreatedAt: util.ToProtoTimestamp(file.CreatedAt),
		UpdatedAt: util.ToProtoTimestamp(file.UpdatedAt),
	}
}

func (fl *FileLogic) UploadFile(in *file_pb.UploadFileReq) (*file_pb.UploadFileResp, error) {
	file, err := fl.svcCtx.FileService.Upload(fl.ctx, in.FileName, in.FileContent)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "upload file failed: %v", err)
	}

	return &file_pb.UploadFileResp{
		FileId:   uint64(file.ID),
		FileName: file.FileName,
		FullPath: file.FullPath,
		Mime:     file.Mime,
		Size:     int64(len(in.FileContent)),
	}, nil
}

func (fl *FileLogic) GetFile(in *common_pb.IdReq) (*file_pb.FileResp, error) {
	file, err := fl.svcCtx.FileService.GetByID(fl.ctx, uint(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get file failed: %v", err)
	}
	if file == nil {
		return nil, status.Errorf(codes.NotFound, "file not found")
	}

	return fl.mapFileToResp(file), nil
}

func (fl *FileLogic) ListFile(in *common_pb.PageReq) (*file_pb.ListFileResp, error) {
	page := &common.PageInfo{
		Page:     int(in.Page),
		PageSize: int(in.PageSize),
	}

	files, count, err := fl.svcCtx.FileService.List(fl.ctx, page)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list files failed: %v", err)
	}

	resp := &file_pb.ListFileResp{
		List:  make([]*file_pb.FileResp, 0, len(files)),
		Total: count,
	}

	for _, file := range files {
		resp.List = append(resp.List, fl.mapFileToResp(file))
	}

	return resp, nil
}

func (fl *FileLogic) DeleteFile(in *common_pb.IdReq) (*common_pb.SuccessResp, error) {
	if in.Id == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid file id")
	}

	err := fl.svcCtx.FileService.Delete(fl.ctx, uint(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "delete file failed: %v", err)
	}

	return &common_pb.SuccessResp{Success: true}, nil
}
