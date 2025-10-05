package implement

import (
	"context"
	"errors"
	"fmt"

	customErr "github.com/tienhai2808/ecom_go/internal/errors"
	"github.com/tienhai2808/ecom_go/internal/model"
	"github.com/tienhai2808/ecom_go/internal/repository"
	"github.com/tienhai2808/ecom_go/internal/request"
	"github.com/tienhai2808/ecom_go/internal/security"
	"github.com/tienhai2808/ecom_go/internal/service"
	"github.com/tienhai2808/ecom_go/internal/snowflake"
)

type userServiceImpl struct {
	userRepo    repository.UserRepository
	profileRepo repository.ProfileRepository
	sfg         snowflake.SnowflakeGenerator
}

func NewUserService(userRepo repository.UserRepository, profileRepo repository.ProfileRepository, sfg snowflake.SnowflakeGenerator) service.UserService {
	return &userServiceImpl{
		userRepo,
		profileRepo,
		sfg,
	}
}

func (s *userServiceImpl) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	users, err := s.userRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("lấy tất cả người dùng thất bại: %w", err)
	}

	return users, nil
}

func (s *userServiceImpl) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	user, err := s.userRepo.FindByIDWithProfile(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("lấy thông tin người dùng thất bại: %w", err)
	}

	if user == nil {
		return nil, customErr.ErrUserNotFound
	}

	return user, nil
}

func (s *userServiceImpl) CreateUser(ctx context.Context, req request.CreateUserRequest) (*model.User, error) {
	exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("kiểm tra người dùng tồn tại thất bại: %w", err)
	}
	if exists {
		return nil, customErr.ErrEmailExists
	}

	exists, err = s.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("kiểm tra người dùng tồn tại thất bại: %w", err)
	}
	if exists {
		return nil, customErr.ErrUsernameExists
	}

	hashedPassword, err := security.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("băm mật khẩu thất bại: %w", err)
	}

	userID, err := s.sfg.NextID()
	if err != nil {
		return nil, err
	}
	profileID, err := s.sfg.NextID()
	if err != nil {
		return nil, err
	}

	newUser := &model.User{
		ID:       userID,
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     req.Role,
		Profile: &model.Profile{
			ID:          profileID,
			FirstName:   req.FirstName,
			LastName:    req.LastName,
			Gender:      req.Gender,
			DOB:         req.DOB,
			PhoneNumber: req.PhoneNumber,
		},
	}

	if err = s.userRepo.Create(ctx, newUser); err != nil {
		return nil, fmt.Errorf("tạo người dùng thất bại: %w", err)
	}

	return newUser, nil
}

func (s *userServiceImpl) UpdateUser(ctx context.Context, id int64, req *request.UpdateUserRequest) (*model.User, error) {
	user, err := s.userRepo.FindByIDWithProfile(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("lấy thông tin người dùng thất bại: %w", err)
	}
	if user == nil {
		return nil, customErr.ErrUserNotFound
	}

	updateUserData := map[string]any{}
	if req.Username != nil && *req.Username != user.Username {
		exists, err := s.userRepo.ExistsByUsername(ctx, *req.Username)
		if err != nil {
			return nil, fmt.Errorf("kiểm tra người dùng tồn tại thất bại: %w", err)
		}
		if exists {
			return nil, customErr.ErrUsernameExists
		}
		updateUserData["username"] = *req.Username
	}
	if req.Email != nil && *req.Email != user.Email {
		exists, err := s.userRepo.ExistsByEmail(ctx, *req.Email)
		if err != nil {
			return nil, fmt.Errorf("kiểm tra người dùng tồn tại thất bại: %w", err)
		}
		if exists {
			return nil, customErr.ErrEmailExists
		}
		updateUserData["email"] = *req.Email
	}
	if req.Password != nil {
		hashedPw, err := security.HashPassword(*req.Password)
		if err != nil {
			return nil, err
		}
		updateUserData["password"] = hashedPw
	}
	if req.Role != nil && *req.Role != user.Role {
		updateUserData["role"] = *req.Role
	}

	updateProfileData := map[string]any{}
	if req.FirstName != nil && *req.FirstName != user.Profile.FirstName {
		updateProfileData["first_name"] = *req.FirstName
	}
	if req.LastName != nil && *req.LastName != user.Profile.LastName {
		updateProfileData["last_name"] = *req.LastName
	}
	if req.PhoneNumber != nil && *req.PhoneNumber != user.Profile.PhoneNumber {
		updateProfileData["phone_number"] = *req.PhoneNumber
	}
	if req.DOB != nil && req.DOB != user.Profile.DOB {
		updateProfileData["dob"] = *req.DOB
	}
	if req.Gender != nil && *req.Gender != user.Profile.Gender {
		updateProfileData["gender"] = *req.Gender
	}

	if len(updateUserData) > 0 {
		if err := s.userRepo.Update(ctx, user.ID, updateUserData); err != nil {
			if errors.Is(err, customErr.ErrUserNotFound) {
				return nil, err
			}
			return nil, fmt.Errorf("cập nhật người dùng thất bại: %w", err)
		}
	}

	if len(updateProfileData) > 0 {
		if err := s.profileRepo.Update(ctx, user.Profile.ID, updateProfileData); err != nil {
			if errors.Is(err, customErr.ErrProfileNotFound) {
				return nil, err
			}
			return nil, fmt.Errorf("cập nhật thông tin người dùng thất bại: %w", err)
		}
	}

	updatedUser, err := s.userRepo.FindByIDWithProfile(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("lấy thông tin người dùng thất bại: %w", err)
	}

	if updatedUser == nil {
		return nil, customErr.ErrUserNotFound
	}

	return updatedUser, nil
}

func (s *userServiceImpl) DeleteUser(ctx context.Context, id int64) error {
	if err := s.userRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, customErr.ErrUserNotFound) {
			return err
		}
		return fmt.Errorf("xóa người dùng thất bại: %w", err)
	}

	return nil
}

func (s *userServiceImpl) DeleteUsers(ctx context.Context, currentUserID int64, req request.DeleteManyRequest) (int64, error) {
	userIDs := req.IDs
	filteredUserIDs := []int64{}

	for _, id := range userIDs {
		if id != currentUserID {
			filteredUserIDs = append(filteredUserIDs, id)
		}
	}

	if len(filteredUserIDs) == 0 {
		return 0, customErr.ErrUserConflict
	}

	rowsAffected, err := s.userRepo.DeleteAllByID(ctx, filteredUserIDs)
	if err != nil {
		return 0, fmt.Errorf("xóa người dùng thất bại: %w", err)
	}

	return rowsAffected, nil
}
