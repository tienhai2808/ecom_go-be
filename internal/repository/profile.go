package repository

import "context"

type ProfileRepository interface {
	UpdateProfileByUserID(ctx context.Context, userID string, updateData map[string]interface{}) error 
}