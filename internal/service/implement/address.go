package implement

import (
	"context"
	"fmt"

	customErr "github.com/tienhai2808/ecom_go/internal/errors"
	"github.com/tienhai2808/ecom_go/internal/model"
	"github.com/tienhai2808/ecom_go/internal/repository"
	"github.com/tienhai2808/ecom_go/internal/request"
	"github.com/tienhai2808/ecom_go/internal/service"
	"github.com/tienhai2808/ecom_go/internal/snowflake"
	"gorm.io/gorm"
)

type addressServiceImpl struct {
	db          *gorm.DB
	addressRepo repository.AddressRepository
	sfg         snowflake.SnowflakeGenerator
}

func NewAddressService(db *gorm.DB, addressRepo repository.AddressRepository, sfg snowflake.SnowflakeGenerator) service.AddressService {
	return &addressServiceImpl{
		db,
		addressRepo,
		sfg,
	}
}

func (s *addressServiceImpl) GetMyAddresses(ctx context.Context, userID int64) ([]*model.Address, error) {
	addresses, err := s.addressRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("lấy địa chỉ người dùng thất bại: %w", err)
	}

	return addresses, nil
}

func (s *addressServiceImpl) GetAddressDetail(ctx context.Context, userID int64, id int64) (*model.Address, error) {
	address, err := s.addressRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("lấy thông tin địa chỉ thất bại: %w", err)
	}

	if address == nil {
		return nil, customErr.ErrAddressNotFound
	}

	if address.UserID != userID {
		return nil, customErr.ErrUnauthorized
	}

	return address, nil
}

func (s *addressServiceImpl) CreateAddress(ctx context.Context, userID int64, req request.AddAddressRequest) (*model.Address, error) {
	addresses, err := s.addressRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("lấy danh sách địa chỉ thất bại: %w", err)
	}

	if len(addresses) >= 10 {
		return nil, customErr.ErrExceedsQuantity
	}

	var exists bool
	var existingID int64
	var newAddress *model.Address
	for _, addr := range addresses {
		if addr.IsDefault {
			exists = true
			existingID = addr.ID
		}
	}

	if err = s.db.Transaction(func(tx *gorm.DB) error {
		if exists && existingID != 0 && *req.IsDefault {
			if err = s.addressRepo.UpdateTx(ctx, tx, existingID, map[string]any{"is_default": false}); err != nil {
				return fmt.Errorf("cập nhật địa chỉ mặc định thất bại: %w", err)
			}
		}

		addressID, err := s.sfg.NextID()
		if err != nil {
			return err
		}

		newAddress = &model.Address{
			ID:          addressID,
			FullName:    req.FullName,
			PhoneNumber: req.PhoneNumber,
			Commune:     req.Commune,
			Province:    req.Province,
			Address:     req.Address,
			IsDefault:   *req.IsDefault,
			UserID:      userID,
		}

		if err = s.addressRepo.CreateTx(ctx, tx, newAddress); err != nil {
			return fmt.Errorf("thêm địa chỉ thất bại: %w", err)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return newAddress, nil
}

func (s *addressServiceImpl) UpdateAddress(ctx context.Context, userID, id int64, req request.UpdateAddressRequest) (*model.Address, error) {
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		address, err := s.addressRepo.FindByIDTx(ctx, tx, id)
		if err != nil {
			return fmt.Errorf("lấy thông tin địa chỉ thất bại: %w", err)
		}
		if address == nil {
			return customErr.ErrAddressNotFound
		}

		if address.UserID != userID {
			return customErr.ErrUnauthorized
		}

		updateData := map[string]any{}
		if req.FullName != nil && *req.FullName != address.FullName {
			updateData["full_name"] = *req.FullName
		}
		if req.PhoneNumber != nil && *req.PhoneNumber != address.PhoneNumber {
			updateData["phone_number"] = *req.PhoneNumber
		}
		if req.Address != nil && *req.Address != address.Address {
			updateData["address"] = *req.Address
		}
		if req.Commune != nil && *req.Commune != address.Commune {
			updateData["commune"] = *req.Commune
		}
		if req.Province != nil && *req.Province != address.Province {
			updateData["province"] = *req.Province
		}

		if req.IsDefault != nil && *req.IsDefault != address.IsDefault {
			if *req.IsDefault {
				if err := s.addressRepo.UpdateTx(ctx, tx, id, map[string]any{"is_default": false}); err != nil {
					return fmt.Errorf("cập nhật địa chỉ mặc định thất bại: %w", err)
				}
			} else {
				countAddress, err := s.addressRepo.CountByUserIDTx(ctx, tx, userID)
				if err != nil {
					return fmt.Errorf("kiểm tra số lượng địa chỉ người dùng thất bại: %w", err)
				}
				if countAddress > 1 {
					latestAddress, err := s.addressRepo.FindLatestByUserIDExcludeIDTx(ctx, tx, userID, id)
					if err != nil {
						return fmt.Errorf("lấy địa chỉ thêm gần đây thất bại: %w", err)
					}
					if latestAddress == nil {
						return customErr.ErrAddressNotFound
					}

					if err = s.addressRepo.UpdateTx(ctx, tx, latestAddress.ID, map[string]any{"is_default": true}); err != nil {
						return fmt.Errorf("cập nhật địa chỉ mặc định thất bại: %w", err)
					}
				} else {
					*req.IsDefault = true
				}
				updateData["is_default"] = *req.IsDefault
			}
		}

		if len(updateData) > 0 {
			if err = s.addressRepo.UpdateTx(ctx, tx, id, updateData); err != nil {
				return fmt.Errorf("cập nhật địa chỉ thất bại: %w", err)
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	updatedAddress, err := s.addressRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("lấy thông tin địa chỉ thất bại: %w", err)
	}
	if updatedAddress == nil {
		return nil, customErr.ErrAddressNotFound
	}

	return updatedAddress, nil
}

func (s *addressServiceImpl) DeleteAddress(ctx context.Context, userID, id int64) error {
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		address, err := s.addressRepo.FindByIDTx(ctx, tx, id)
		if err != nil {
			return fmt.Errorf("lấy thông tin địa chỉ thất bại: %w", err)
		}
		if address == nil {
			return customErr.ErrAddressNotFound
		}

		if address.UserID != userID {
			return customErr.ErrUnauthorized
		}

		countAddress, err := s.addressRepo.CountByUserIDTx(ctx, tx, userID)
		if err != nil {
			return fmt.Errorf("kiểm tra số lượng địa chỉ người dùng thất bại: %w", err)
		}

		if countAddress > 1 && address.IsDefault {
			latestAddress, err := s.addressRepo.FindLatestByUserIDExcludeIDTx(ctx, tx, userID, id)
			if err != nil {
				return fmt.Errorf("lấy địa chỉ thêm gần đây thất bại: %w", err)
			}
			if latestAddress == nil {
				return customErr.ErrAddressNotFound
			}

			if err = s.addressRepo.UpdateTx(ctx, tx, latestAddress.ID, map[string]any{"is_default": true}); err != nil {
				return fmt.Errorf("cập nhật địa chỉ mặc định thất bại: %w", err)
			}
		}

		if err = s.addressRepo.DeleteTx(ctx, tx, id); err != nil {
			return fmt.Errorf("xóa địa chỉ thất bại: %w", err)
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}
