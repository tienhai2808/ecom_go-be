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
)

type profileServiceImpl struct {
	profileRepo repository.ProfileRepository
	userRepo    repository.UserRepository
}

func NewProfileService(profileRepo repository.ProfileRepository, userRepo repository.UserRepository) service.ProfileService {
	return &profileServiceImpl{
		profileRepo,
		userRepo,
	}
}

func (s *profileServiceImpl) UpdateProfile(ctx context.Context, id int64, req *request.UpdateProfileRequest) (*model.User, error) {
	profile, err := s.profileRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("lấy thông tin người dùng thất bại: %w", err)
	}
	if profile == nil {
		return nil, customErr.ErrProfileNotFound
	}

	updateData := map[string]any{}
	if req.FirstName != nil && *req.FirstName != profile.FirstName {
		updateData["first_name"] = *req.FirstName
	}
	if req.LastName != nil && *req.LastName != profile.LastName {
		updateData["last_name"] = *req.LastName
	}
	if req.Gender != nil && *req.Gender != profile.Gender {
		updateData["gender"] = *req.Gender
	}
	if req.DOB != nil && req.DOB != profile.DOB {
		updateData["dob"] = *req.DOB
	}
	if req.PhoneNumber != nil && *req.PhoneNumber != profile.PhoneNumber {
		updateData["phone_number"] = *req.PhoneNumber
	}

	if len(updateData) > 0 {
		if err := s.profileRepo.Update(ctx, id, updateData); err != nil {
			if errors.Is(err, customErr.ErrProfileNotFound) {
				return nil, err
			}
			return nil, fmt.Errorf("cập nhật thông tin người dùng thất bại: %w", err)
		}
	}

	user, err := s.userRepo.FindByIDWithProfile(ctx, profile.UserID)
	if err != nil {
		return nil, fmt.Errorf("lấy thông tin người dùng thất bại: %w", err)
	}

	if user == nil {
		return nil, customErr.ErrUserNotFound
	}

	return user, nil
}
