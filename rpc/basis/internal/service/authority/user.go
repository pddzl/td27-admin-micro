package authority

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"td27/rpc/basis/internal/model/authority"
	"td27/rpc/basis/internal/model/common"
	repoAuthority "td27/rpc/basis/internal/repository/authority"
)

type UserService struct {
	userRepo repoAuthority.UserRepository
	roleRepo repoAuthority.RoleRepository
}

func NewUserService(userRepo repoAuthority.UserRepository, roleRepo repoAuthority.RoleRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

func (s *UserService) GetByID(ctx context.Context, id uint) (*authority.UserModel, error) {
	return s.userRepo.FindOne(ctx, id)
}

func (s *UserService) GetByIDWithRoles(ctx context.Context, id uint) (*authority.UserModel, error) {
	return s.userRepo.FindOneWithRoles(ctx, id)
}

func (s *UserService) GetByUsername(ctx context.Context, username string) (*authority.UserModel, error) {
	return s.userRepo.FindOneByUsername(ctx, username)
}

func (s *UserService) List(ctx context.Context, page *common.PageInfo, deptID *uint) ([]*authority.UserModel, int64, error) {
	return s.userRepo.List(ctx, page, deptID)
}

func (s *UserService) Create(ctx context.Context, user *authority.UserModel, password string) error {
	existing, err := s.userRepo.FindOneByUsername(ctx, user.Username)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("username already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = string(hashedPassword)

	return s.userRepo.Create(ctx, user)
}

func (s *UserService) Update(ctx context.Context, user *authority.UserModel) error {
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

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.New("old password is incorrect")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	return s.userRepo.UpdatePassword(ctx, userID, string(hashedPassword))
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

func (s *UserService) GetUserRoles(ctx context.Context, userID uint) ([]*authority.RoleModel, error) {
	return s.roleRepo.GetUserRoles(ctx, userID)
}

func (s *UserService) VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
