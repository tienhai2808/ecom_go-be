package model

type Cart struct {
	ID            int64   `gorm:"type:bigint;primaryKey" json:"id"`
	TotalPrice    float64 `gorm:"type:decimal(10,2);not null" json:"total_price"`
	TotalQuantity uint    `gorm:"type:int;not null" json:"total_quantity"`
	UserID        int64   `gorm:"type:bigint;not null;unique" json:"user_id"`

	User      *User       `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user"`
	CartItems []*CartItem `gorm:"foreignKey:CartID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"cart_items"`
}

type CartItem struct {
	ID         int64   `gorm:"type:bigint;primaryKey" json:"id"`
	UnitPrice  float64 `gorm:"type:decimal(10,2);not null" json:"unit_price"`
	Quantity   uint    `gorm:"type:int;not null" json:"quantity"`
	TotalPrice float64 `gorm:"type:decimal(10,2);not null" json:"total_price"`
	CartID     int64   `gorm:"type:bigint;not null" json:"cart_id"`
	ProductID  int64   `gorm:"type:bigint;not null" json:"product_id"`

	Cart    *Cart    `gorm:"foreignKey:CartID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"cart"`
	Product *Product `gorm:"foreignKey:ProductID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"product"`
}

func (m *CartItem) SetTotalPrice() {
	m.TotalPrice = m.UnitPrice * float64(m.Quantity)
}