package tool

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/sha3"

	"td27/rpc/basis/internal/model/common"
	"td27/rpc/basis/internal/model/tool"
	repoTool "td27/rpc/basis/internal/repository/tool"
)

type ServiceTokenService struct {
	tokenRepo repoTool.ServiceTokenRepository
}

func NewServiceTokenService(tokenRepo repoTool.ServiceTokenRepository) *ServiceTokenService {
	return &ServiceTokenService{
		tokenRepo: tokenRepo,
	}
}

func (s *ServiceTokenService) GetByID(ctx context.Context, id uint) (*tool.ServiceToken, error) {
	return s.tokenRepo.FindOne(ctx, id)
}

func (s *ServiceTokenService) GetByToken(ctx context.Context, token string) (*tool.ServiceToken, error) {
	hash := s.hashToken(token)
	return s.tokenRepo.FindByTokenHash(ctx, hash)
}

func (s *ServiceTokenService) List(ctx context.Context, page *common.PageInfo) ([]*tool.ServiceToken, int64, error) {
	return s.tokenRepo.List(ctx, page)
}

func (s *ServiceTokenService) Create(ctx context.Context, token *tool.ServiceToken, ttl *time.Duration) (string, error) {
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	rawToken := hex.EncodeToString(randomBytes)

	token.TokenHash = s.hashToken(rawToken)

	if ttl != nil {
		expiresAt := time.Now().Add(*ttl).Unix()
		token.ExpiresAt = &expiresAt
	}

	if err := s.tokenRepo.Create(ctx, token); err != nil {
		return "", err
	}

	return rawToken, nil
}

func (s *ServiceTokenService) Update(ctx context.Context, token *tool.ServiceToken) error {
	return s.tokenRepo.Update(ctx, token)
}

func (s *ServiceTokenService) ToggleStatus(ctx context.Context, id uint, status bool) error {
	return s.tokenRepo.ToggleStatus(ctx, id, status)
}

func (s *ServiceTokenService) Delete(ctx context.Context, id uint) error {
	return s.tokenRepo.Delete(ctx, id)
}

func (s *ServiceTokenService) AssignPermissions(ctx context.Context, tokenID uint, permissionIDs []uint) error {
	token, err := s.tokenRepo.FindOne(ctx, tokenID)
	if err != nil {
		return err
	}
	if token == nil {
		return errors.New("service token not found")
	}
	return s.tokenRepo.AssignPermissions(ctx, tokenID, permissionIDs)
}

func (s *ServiceTokenService) GetPermissions(ctx context.Context, tokenID uint) ([]uint, error) {
	return s.tokenRepo.GetPermissions(ctx, tokenID)
}

func (s *ServiceTokenService) hashToken(token string) string {
	hash := sha3.New256()
	hash.Write([]byte(token))
	return hex.EncodeToString(hash.Sum(nil))
}
