package sysTool

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	"td27/rpc/basis/internal/config"
	"td27/rpc/basis/internal/model/common"
	"td27/rpc/basis/internal/model/sysTool"
	sysToolRepo "td27/rpc/basis/internal/repository/sysTool"
)

type FileService struct {
	fileRepo sysToolRepo.FileRepository
	config   config.File
}

func NewFileService(fileRepo sysToolRepo.FileRepository, config config.File) *FileService {
	return &FileService{
		fileRepo: fileRepo,
		config:   config,
	}
}

func (s *FileService) GetByID(ctx context.Context, id uint) (*sysTool.FileModel, error) {
	return s.fileRepo.FindOne(ctx, id)
}

func (s *FileService) List(ctx context.Context, page *common.PageInfo) ([]*sysTool.FileModel, int64, error) {
	return s.fileRepo.List(ctx, page)
}

func (s *FileService) Upload(ctx context.Context, fileName string, fileContent []byte) (*sysTool.FileModel, error) {
	ext := strings.ToLower(filepath.Ext(fileName))
	allowed := false
	for _, allowedExt := range s.config.AllowedExt {
		if ext == allowedExt {
			allowed = true
			break
		}
	}
	if !allowed {
		return nil, fmt.Errorf("file extension %s is not allowed", ext)
	}

	if len(fileContent) > s.config.MaxSize*1024*1024 {
		return nil, fmt.Errorf("file size exceeds maximum allowed size of %d MB", s.config.MaxSize)
	}

	uploadPath := s.config.UploadPath
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	newFileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	fullPath := filepath.Join(uploadPath, newFileName)

	if err := os.WriteFile(fullPath, fileContent, 0644); err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	file := &sysTool.FileModel{
		FileName: fileName,
		FullPath: fullPath,
		Mime:     http.DetectContentType(fileContent),
	}

	if err := s.fileRepo.Create(ctx, file); err != nil {
		os.Remove(fullPath)
		return nil, err
	}

	return file, nil
}

func (s *FileService) DownloadFile(ctx context.Context, id uint) (fileName, mime string, content []byte, err error) {
	file, err := s.fileRepo.FindOne(ctx, id)
	if err != nil {
		return "", "", nil, err
	}
	if file == nil {
		return "", "", nil, errors.New("file not found")
	}

	// Prevent path traversal - reject paths with ".."
	if strings.Contains(file.FullPath, "..") {
		return "", "", nil, errors.New("invalid file path")
	}

	content, err = os.ReadFile(file.FullPath)
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to read file: %w", err)
	}

	return file.FileName, file.Mime, content, nil
}

func (s *FileService) Delete(ctx context.Context, id uint) error {
	file, err := s.fileRepo.FindOne(ctx, id)
	if err != nil {
		return err
	}
	if file == nil {
		return errors.New("file not found")
	}

	if err = os.Remove(file.FullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file from disk: %w", err)
	}

	return s.fileRepo.Delete(ctx, id)
}
