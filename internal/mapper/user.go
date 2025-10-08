package mapper

import (
	"github.com/tienhai2808/ecom_go/internal/model"
	"github.com/tienhai2808/ecom_go/internal/response"
	"github.com/tienhai2808/ecom_go/internal/types"
)

func ToUserResponse(user *model.User) *response.UserResponse {
	return &response.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		Profile: &response.ProfileResponse{
			ID:          user.Profile.ID,
			FirstName:   user.Profile.FirstName,
			LastName:    user.Profile.LastName,
			PhoneNumber: user.Profile.PhoneNumber,
			DOB:         user.Profile.DOB,
			Gender:      user.Profile.Gender,
			UpdatedAt:   user.Profile.UpdatedAt,
		},
	}
}

func ToBaseUserResponse(user *model.User) *response.BaseUserResponse {
	return &response.BaseUserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}
}

func ToUserData(user *model.User) *types.UserData {
	return &types.UserData{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		Profile: types.ProfileData{
			ID:          user.Profile.ID,
			FirstName:   user.Profile.FirstName,
			LastName:    user.Profile.LastName,
			PhoneNumber: user.Profile.PhoneNumber,
			DOB:         user.Profile.DOB,
			Gender:      user.Profile.Gender,
			UpdatedAt:   user.Profile.UpdatedAt,
		},
	}
}

func ToUsersResponse(usrs []*model.User) []*response.BaseUserResponse {
	if len(usrs) == 0 {
		return make([]*response.BaseUserResponse, 0)
	}

	usrsResp := make([]*response.BaseUserResponse, 0, len(usrs))
	for _, usr := range usrs {
		usrsResp = append(usrsResp, ToBaseUserResponse(usr))
	}

	return usrsResp
}
