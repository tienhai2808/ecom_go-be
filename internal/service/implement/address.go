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

	"github.com/google/uuid"
)

type addressServiceImpl struct {
	addressRepository repository.AddressRepository
}

func NewAddressService(addressRepository repository.AddressRepository) service.AddressService {
	return &addressServiceImpl{
		addressRepository: addressRepository,
	}
}

func (s *addressServiceImpl) GetMyAddresses(ctx context.Context, userID string) ([]*model.Address, error) {
	addresses, err := s.addressRepository.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("lấy địa chỉ người dùng thất bại: %w", err)
	}

	return addresses, nil
}

func (s *addressServiceImpl) GetAddressDetail(ctx context.Context, userID string, id string) (*model.Address, error) {
	address, err := s.addressRepository.FindByID(ctx, id)
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

func (s *addressServiceImpl) CreateAddress(ctx context.Context, userID string, req request.AddAddressRequest) (*model.Address, error) {
	
	
	count, err := s.addressRepository.CountByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("kiểm tra số lượng địa chỉ người dùng thất bại: %w", err)
	}

	if count >= 10 {
		return nil, customErr.ErrExceedsQuantity
	}

	exists, err := s.addressRepository.CheckDefaultExistsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("kiểm tra tồn tại địa chỉ mặc định thất bại: %w", err)
	}

	if exists && *req.IsDefault {
		if err := s.addressRepository.UpdateNonDefaultByUserID(ctx, userID); err != nil {
			if errors.Is(err, customErr.ErrUserAddressNotFound) {
				return nil, err
			}
			return nil, fmt.Errorf("cập nhật địa chỉ mặc định thất bại: %w", err)
		}
	}

	newAddress := &model.Address{
		ID:          uuid.NewString(),
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		Commune:     req.Commune,
		District:    req.District,
		Province:    req.Province,
		Address:     req.Address,
		IsDefault:   *req.IsDefault,
		UserID:      userID,
	}

	if err := s.addressRepository.Create(ctx, newAddress); err != nil {
		return nil, fmt.Errorf("thêm địa chỉ thất bại: %w", err)
	}

	return newAddress, nil
}

func (s *addressServiceImpl) UpdateAddress(ctx context.Context, userID, id string, req *request.UpdateAddressRequest) (*model.Address, error) {
	address, err := s.addressRepository.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("lấy thông tin địa chỉ thất bại: %w", err)
	}

	if address == nil {
		return nil, customErr.ErrAddressNotFound
	}

	if address.UserID != userID {
		return nil, customErr.ErrUnauthorized
	}

	updateData := map[string]interface{}{}
	if req.FirstName != nil && *req.FirstName != address.FirstName {
		updateData["first_name"] = *req.FirstName
	}
	if req.LastName != nil && *req.LastName != address.LastName {
		updateData["last_name"] = *req.LastName
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
	if req.District != nil && *req.District != address.District {
		updateData["district"] = *req.District
	}
	if req.Province != nil && *req.Province != address.Province {
		updateData["province"] = *req.Province
	}

	if req.IsDefault != nil && *req.IsDefault != address.IsDefault {
		if *req.IsDefault {
			if err := s.addressRepository.UpdateNonDefaultByUserID(ctx, userID); err != nil {
				if errors.Is(err, customErr.ErrUserAddressNotFound) {
					return nil, err
				}
				return nil, fmt.Errorf("cập nhật địa chỉ mặc định thất bại: %w", err)
			}
		} else {
			countAddress, err := s.addressRepository.CountByUserID(ctx, userID)
			if err != nil {
				return nil, fmt.Errorf("kiểm tra số lượng địa chỉ người dùng thất bại: %w", err)
			}
			if countAddress > 1 {
				latestAddress, err := s.addressRepository.FindLatestByUserIDExcludeID(ctx, userID, id)
				if err != nil {
					return nil, fmt.Errorf("lấy địa chỉ thêm gần đây thất bại: %w", err)
				}

				if latestAddress == nil {
					return nil, customErr.ErrAddressNotFound
				}

				if err := s.addressRepository.UpdateDefault(ctx, latestAddress.ID); err != nil {
					if errors.Is(err, customErr.ErrAddressNotFound) {
						return nil, err
					}
					return nil, fmt.Errorf("cập nhật địa chỉ mặc định thất bại: %w", err)
				}
			} else {
				*req.IsDefault = true
			}
			updateData["is_default"] = *req.IsDefault
		}
	}

	if len(updateData) > 0 {
		if err := s.addressRepository.Update(ctx, id, updateData); err != nil {
			if errors.Is(err, customErr.ErrAddressNotFound) {
				return nil, err
			}
			return nil, fmt.Errorf("cập nhật địa chỉ thất bại: %w", err)
		}
	}

	updatedAddress, err := s.addressRepository.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("lấy thông tin địa chỉ thất bại: %w", err)
	}

	if updatedAddress == nil {
		return nil, customErr.ErrAddressNotFound
	}

	return updatedAddress, nil
}

func (s *addressServiceImpl) DeleteAddress(ctx context.Context, userID, id string) error {
	address, err := s.addressRepository.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("lấy thông tin địa chỉ thất bại: %w", err)
	}

	if address == nil {
		return customErr.ErrAddressNotFound
	}

	if address.UserID != userID {
		return customErr.ErrUnauthorized
	}

	countAddress, err := s.addressRepository.CountByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("kiểm tra số lượng địa chỉ người dùng thất bại: %w", err)
	}

	if countAddress > 1 && address.IsDefault {
		latestAddress, err := s.addressRepository.FindLatestByUserIDExcludeID(ctx, userID, id)
		if err != nil {
			return fmt.Errorf("lấy địa chỉ thêm gần đây thất bại: %w", err)
		}

		if latestAddress == nil {
			return customErr.ErrAddressNotFound
		}

		if err := s.addressRepository.UpdateDefault(ctx, latestAddress.ID); err != nil {
			if errors.Is(err, customErr.ErrAddressNotFound) {
				return err
			}
			return fmt.Errorf("cập nhật địa chỉ mặc định thất bại: %w", err)
		}
	}

	if err = s.addressRepository.Delete(ctx, id); err != nil {
		if errors.Is(err, customErr.ErrAddressNotFound) {
			return err
		}
		return fmt.Errorf("xóa địa chỉ thất bại: %w", err)
	}

	return nil
}
