package repository

import "context"

type ProfileRepository interface {
	UpdateByUserID(ctx context.Context, userID int64, updateData map[string]any) error 
}