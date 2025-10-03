package implement

import (
	"context"
	"errors"
	"fmt"
	customErr "github.com/tienhai2808/ecom_go/internal/errors"
	"github.com/tienhai2808/ecom_go/internal/model"
	"github.com/tienhai2808/ecom_go/internal/repository"
	"github.com/tienhai2808/ecom_go/internal/request"
	"github.com/tienhai2808/ecom_go/internal/service"
	"github.com/tienhai2808/ecom_go/internal/util"
)

type userServiceImpl struct {
	userRepository    repository.UserRepository
	profileRepository repository.ProfileRepository
}

func NewUserService(userRepository repository.UserRepository, profileRepository repository.ProfileRepository) service.UserService {
	return &userServiceImpl{
		userRepository:    userRepository,
		profileRepository: profileRepository,
	}
}

func (s *userServiceImpl) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	users, err := s.userRepository.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("lấy tất cả người dùng thất bại: %w", err)
	}

	return users, nil
}

func (s *userServiceImpl) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	user, err := s.userRepository.FindByIDWithProfile(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("lấy thông tin người dùng thất bại: %w", err)
	}

	if user == nil {
		return nil, customErr.ErrUserNotFound
	}

	return user, nil
}

func (s *userServiceImpl) CreateUser(ctx context.Context, req request.CreateUserRequest) (*model.User, error) {
	exists, err := s.userRepository.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("kiểm tra người dùng tồn tại thất bại: %w", err)
	}

	if exists {
		return nil, customErr.ErrEmailExists
	}

	exists, err = s.userRepository.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("kiểm tra người dùng tồn tại thất bại: %w", err)
	}

	if exists {
		return nil, customErr.ErrUsernameExists
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("băm mật khẩu thất bại: %w", err)
	}

	userID, err := util.NewSnowflakeID()
	if err != nil {
		return nil, err
	}
	profileID, err := util.NewSnowflakeID()
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
			DOB:         &req.DOB,
			PhoneNumber: req.PhoneNumber,
		},
	}

	if err = s.userRepository.Create(ctx, newUser); err != nil {
		return nil, fmt.Errorf("tạo người dùng thất bại: %w", err)
	}

	return newUser, nil
}

func (s *userServiceImpl) UpdateUser(ctx context.Context, id int64, req *request.UpdateUserRequest) (*model.User, error) {
	user, err := s.userRepository.FindByIDWithProfile(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("lấy thông tin người dùng thất bại: %w", err)
	}

	if user == nil {
		return nil, customErr.ErrUserNotFound
	}

	if req.Email != nil && *req.Email != user.Email {
		exists, err := s.userRepository.ExistsByEmail(ctx, *req.Email)
		if err != nil {
			return nil, fmt.Errorf("kiểm tra người dùng tồn tại thất bại: %w", err)
		}

		if exists {
			return nil, customErr.ErrEmailExists
		}
	}

	if req.Username != nil && *req.Username != user.Username {
		exists, err := s.userRepository.ExistsByUsername(ctx, *req.Username)
		if err != nil {
			return nil, fmt.Errorf("kiểm tra người dùng tồn tại thất bại: %w", err)
		}

		if exists {
			return nil, customErr.ErrUsernameExists
		}
	}

	updateUserData := map[string]any{}
	if req.Username != nil {
		updateUserData["username"] = *req.Username
	}
	if req.Email != nil {
		updateUserData["email"] = *req.Email
	}
	if req.Password != nil {
		updateUserData["password"] = *req.Password
	}
	if req.Role != nil {
		updateUserData["role"] = *req.Role
	}

	updateProfileData := map[string]any{}
	if req.FirstName != nil {
		updateProfileData["first_name"] = *req.FirstName
	}
	if req.LastName != nil {
		updateProfileData["last_name"] = *req.LastName
	}
	if req.PhoneNumber != nil {
		updateProfileData["phone_number"] = *req.PhoneNumber
	}
	if req.DOB != nil {
		updateProfileData["dob"] = *req.DOB
	}
	if req.Gender != nil {
		updateProfileData["gender"] = *req.Gender
	}

	if len(updateUserData) > 0 {
		if err := s.userRepository.Update(ctx, user.ID, updateUserData); err != nil {
			if errors.Is(err, customErr.ErrUserNotFound) {
				return nil, err
			}
			return nil, fmt.Errorf("cập nhật người dùng thất bại: %w", err)
		}
	}

	if len(updateProfileData) > 0 {
		if err := s.profileRepository.UpdateByUserID(ctx, user.ID, updateProfileData); err != nil {
			if errors.Is(err, customErr.ErrUserProfileNotFound) {
				return nil, err
			}
			return nil, fmt.Errorf("cập nhật thông tin người dùng thất bại: %w", err)
		}
	}

	updatedUser, err := s.userRepository.FindByIDWithProfile(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("lấy thông tin người dùng thất bại: %w", err)
	}

	if updatedUser == nil {
		return nil, customErr.ErrUserNotFound
	}

	return updatedUser, nil
}

func (s *userServiceImpl) DeleteUser(ctx context.Context, id int64) error {
	if err := s.userRepository.Delete(ctx, id); err != nil {
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

	rowsAffected, err := s.userRepository.DeleteAllByID(ctx, filteredUserIDs)
	if err != nil {
		return 0, fmt.Errorf("xóa người dùng thất bại: %w", err)
	}

	return rowsAffected, nil
}
