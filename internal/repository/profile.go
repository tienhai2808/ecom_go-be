package repository

import (
	"context"

	"github.com/tienhai2808/ecom_go/internal/model"
)

type ProfileRepository interface {
	FindByID(ctx context.Context, id int64) (*model.Profile, error)

	Update(ctx context.Context, id int64, updateData map[string]any) error 
}