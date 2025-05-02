package admin

import (
	"backend/internal/auth"
	"backend/internal/common"
	"backend/internal/config"
	"backend/internal/user"
	"fmt"

	"github.com/google/uuid"
)

type Service interface {
	GetAllUsers() ([]user.User, error)
	CreateUser(req CreateUserRequest) error
	UpdateUser(userID string, req *UpdateUserRequest) error
	DeleteUsers(currentUserID string, req DeleteUsersRequest) (int64, error)
}

type service struct {
	repo     Repository
	authRepo auth.Repository
	cfg      *config.AppConfig
}

func NewService(ctx *common.AppContext) Service {
	return &service{
		repo:     NewRepository(ctx),
		authRepo: auth.NewRepository(ctx),
		cfg:      ctx.Config,
	}
}

func (s *service) GetAllUsers() ([]user.User, error) {
	return s.repo.GetAllUsers()
}

func (s *service) CreateUser(req CreateUserRequest) error {
	exists, field, err := s.authRepo.CheckUserExists(req.Username, req.Email)
	if err != nil {
		return err
	}

	if exists {
		if field == "username" {
			return auth.ErrUsernameExists
		}
		return auth.ErrEmailExists
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return err
	}

	newUser := &user.User{
		ID:       uuid.NewString(),
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     req.Role,
		Profile: user.Profile{
			ID: uuid.NewString(),
			FirstName: req.FirstName,
			LastName: req.LastName,
			Gender: req.Gender,
			DOB: req.DOB,
			PhoneNumber: req.PhoneNumber,
		},
	}

	if err := s.authRepo.CreateUser(newUser); err != nil {
		return fmt.Errorf("không thể tạo người dùng: %v", err)
	}

	return nil
}

func (s *service) UpdateUser(userID string, req *UpdateUserRequest) error {
	user, err := s.authRepo.GetUserByID(userID)
	if err != nil {
		return auth.ErrUserNotFound
	}

	if req.Email != nil && *req.Email != user.Email {
		exists, err := s.authRepo.CheckUserExistsByEmail(*req.Email) 
		if err != nil {
			return err
		}
		if exists {
			return auth.ErrEmailExists
		}
	}

	if req.Username != nil && *req.Username != user.Username {
		exists, err := s.authRepo.CheckUserExistsByUsername(*req.Username)
		if err != nil {
			return err
		}
		if exists {
			return auth.ErrUsernameExists
		}
	}

	updateUserData := map[string]interface{}{}
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

	updateProfileData := map[string]interface{}{}
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
		if err := s.authRepo.UpdateUserInfo(user, updateUserData); err != nil {
			return auth.ErrUpdateFailed
		}
	}

	if len(updateProfileData) > 0 {
		if err := s.authRepo.UpdateUserProfile(user, updateProfileData); err != nil {
			return auth.ErrUpdateFailed
		}
	}

	return nil 
}

func (s *service) DeleteUsers(currentUserID string, req DeleteUsersRequest) (int64, error) {
	userIDs := req.UserIds
	if len(userIDs) == 0 {
		return 0, fmt.Errorf("không có user_id trong yêu cầu")
	}

	filteredUserIDs := []string{}
	for _, id := range userIDs {
		if id != currentUserID {
			filteredUserIDs = append(filteredUserIDs, id)
		}
	}

	
	if len(filteredUserIDs) == 0 {
		return 0, fmt.Errorf("không thể xóa tài khoản đang đăng nhập")
	}

	rowsAffected, err := s.repo.DeleteUSers(filteredUserIDs); 
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}
