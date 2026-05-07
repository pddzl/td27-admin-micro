package tool

import (
	"context"

	"td27/rpc/basis/internal/logic/tool"
	"td27/rpc/basis/internal/svc"
	"td27/rpc/basis/types/tool/file_pb"
	"td27/rpc/basis/types/common_pb"
)

type FileServer struct {
	svcCtx *svc.ServiceContext
	file_pb.UnimplementedFileServer
}

func NewFileServer(svcCtx *svc.ServiceContext) *FileServer {
	return &FileServer{
		svcCtx: svcCtx,
	}
}

func (fs *FileServer) UploadFile(ctx context.Context, in *file_pb.UploadFileReq) (*file_pb.UploadFileResp, error) {
	fl := tool.NewFileLogic(ctx, fs.svcCtx)
	return fl.UploadFile(in)
}

func (fs *FileServer) GetFile(ctx context.Context, in *common_pb.IdReq) (*file_pb.FileResp, error) {
	fl := tool.NewFileLogic(ctx, fs.svcCtx)
	return fl.GetFile(in)
}

func (fs *FileServer) ListFile(ctx context.Context, in *common_pb.PageReq) (*file_pb.ListFileResp, error) {
	fl := tool.NewFileLogic(ctx, fs.svcCtx)
	return fl.ListFile(in)
}

func (fs *FileServer) DeleteFile(ctx context.Context, in *common_pb.IdReq) (*common_pb.SuccessResp, error) {
	fl := tool.NewFileLogic(ctx, fs.svcCtx)
	return fl.DeleteFile(in)
}
