package mapper

import (
	"github.com/tienhai2808/ecom_go/internal/model"
	"github.com/tienhai2808/ecom_go/internal/response"
)

func ToAddressResponse(addr *model.Address) *response.AddressResponse {
	return &response.AddressResponse{
		ID:          addr.ID,
		FullName:    addr.FullName,
		PhoneNumber: addr.PhoneNumber,
		Address:     addr.Address,
		Commune:     addr.Commune,
		Province:    addr.Province,
		CreatedAt:   addr.CreatedAt,
		UpdatedAt:   addr.UpdatedAt,
	}
}

func ToBaseAddressResponse(addr *model.Address) *response.BaseAddressResponse {
	return &response.BaseAddressResponse{
		ID:          addr.ID,
		FullName:    addr.FullName,
		PhoneNumber: addr.PhoneNumber,
		Address:     addr.Address,
		Commune:     addr.Commune,
		Province:    addr.Province,
	}
}

func ToAddressesResponse(addrs []*model.Address) []*response.BaseAddressResponse {
	if len(addrs) == 0 {
		return make([]*response.BaseAddressResponse, 0)
	}

	addrsResp := make([]*response.BaseAddressResponse, 0, len(addrs))
	for _, addr := range addrs {
		addrsResp = append(addrsResp, ToBaseAddressResponse(addr))
	}

	return addrsResp
}
