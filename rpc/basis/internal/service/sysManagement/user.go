package sysManagement

import (
	"context"
	"errors"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"

	"td27/pkg/tool"

	"td27/rpc/basis/internal/model/common"
	"td27/rpc/basis/internal/model/sysManagement"
	sysManagementRepo "td27/rpc/basis/internal/repository/sysManagement"
)

type UserService struct {
	userRepo sysManagementRepo.UserRepository
	roleRepo sysManagementRepo.RoleRepository
}

func NewUserService(userRepo sysManagementRepo.UserRepository, roleRepo sysManagementRepo.RoleRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

func (s *UserService) GetByID(ctx context.Context, id uint) (*sysManagement.UserModel, error) {
	return s.userRepo.FindOne(ctx, id)
}

func (s *UserService) GetByIDWithRoles(ctx context.Context, id uint) (*sysManagement.UserModel, error) {
	return s.userRepo.FindOneWithRoles(ctx, id)
}

func (s *UserService) GetByUsername(ctx context.Context, username string) (*sysManagement.UserModel, error) {
	return s.userRepo.FindOneByUsername(ctx, username)
}

func (s *UserService) List(ctx context.Context, page *common.PageInfo, deptID *uint) ([]*sysManagement.UserModel, int64, error) {
	return s.userRepo.List(ctx, page, deptID)
}

func (s *UserService) Create(ctx context.Context, user *sysManagement.UserModel, password string) error {
	existing, err := s.userRepo.FindOneByUsername(ctx, user.Username)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("username already exists")
	}

	hash, err := tool.HashPassword(password)
	if err != nil {
		return fmt.Errorf("hash password failed: %w", err)
	}
	user.Password = hash

	return s.userRepo.Create(ctx, user)
}

func (s *UserService) Update(ctx context.Context, user *sysManagement.UserModel) error {
	return s.userRepo.Update(ctx, user)
}

func (s *UserService) Delete(ctx context.Context, id uint) error {
	user, err := s.userRepo.FindOne(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	return s.userRepo.Delete(ctx, id)
}

func (s *UserService) ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error {
	user, err := s.userRepo.FindOne(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	if !tool.CheckPassword(user.Password, oldPassword) {
		return errors.New("old password is incorrect")
	}

	hash, err := tool.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("hash password failed: %w", err)
	}

	return s.userRepo.UpdatePassword(ctx, userID, hash)
}

func (s *UserService) ToggleActive(ctx context.Context, userID uint, active bool) error {
	user, err := s.userRepo.FindOne(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	return s.userRepo.ToggleActive(ctx, userID, active)
}

func (s *UserService) AssignRoles(ctx context.Context, userID uint, roleIDs []uint) error {
	for _, roleID := range roleIDs {
		role, err := s.roleRepo.FindOne(ctx, roleID)
		if err != nil {
			return err
		}
		if role == nil {
			return fmt.Errorf("role with id %d not found", roleID)
		}
	}

	return s.roleRepo.AssignUserRoles(ctx, userID, roleIDs)
}

func (s *UserService) GetUserRoles(ctx context.Context, userID uint) ([]*sysManagement.RoleModel, error) {
	return s.roleRepo.GetUserRoles(ctx, userID)
}

func (s *UserService) GetRolesByUserIDs(ctx context.Context, userIDs []uint) (map[uint][]*sysManagement.RoleModel, error) {
	return s.roleRepo.GetRolesByUserIDs(ctx, userIDs)
}

// VerifyPassword verifies the plaintext password against the stored hash.
// Legacy MD5 hashes are accepted and transparently upgraded to bcrypt after a
// successful check; upgrade failures are logged but do not block the login.
func (s *UserService) VerifyPassword(ctx context.Context, user *sysManagement.UserModel, password string) bool {
	if !tool.CheckPassword(user.Password, password) {
		return false
	}

	if tool.IsLegacyHash(user.Password) {
		if hash, err := tool.HashPassword(password); err == nil {
			if err := s.userRepo.UpdatePassword(ctx, user.ID, hash); err != nil {
				logx.Errorf("upgrade legacy password hash for user %d failed: %v", user.ID, err)
			}
		} else {
			logx.Errorf("hash password for legacy upgrade of user %d failed: %v", user.ID, err)
		}
	}

	return true
}
