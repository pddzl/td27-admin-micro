package sysTool

import (
	"context"

	"td27/rpc/basis/internal/logic/sysTool"
	"td27/rpc/basis/internal/svc"
	"td27/rpc/basis/types/common_pb"
	"td27/rpc/basis/types/sysTool/file_pb"
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
	fl := sysTool.NewFileLogic(ctx, fs.svcCtx)
	return fl.UploadFile(in)
}

func (fs *FileServer) GetFile(ctx context.Context, in *common_pb.IdReq) (*file_pb.FileResp, error) {
	fl := sysTool.NewFileLogic(ctx, fs.svcCtx)
	return fl.GetFile(in)
}

func (fs *FileServer) ListFile(ctx context.Context, in *common_pb.PageReq) (*file_pb.ListFileResp, error) {
	fl := sysTool.NewFileLogic(ctx, fs.svcCtx)
	return fl.ListFile(in)
}

func (fs *FileServer) DeleteFile(ctx context.Context, in *common_pb.IdReq) (*common_pb.SuccessResp, error) {
	fl := sysTool.NewFileLogic(ctx, fs.svcCtx)
	return fl.DeleteFile(in)
}

func (fs *FileServer) DownloadFile(ctx context.Context, in *common_pb.IdReq) (*file_pb.DownloadFileResp, error) {
	fl := sysTool.NewFileLogic(ctx, fs.svcCtx)
	return fl.DownloadFile(in)
}
