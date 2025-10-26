package model

type Order struct {
	ID            int64   `gorm:"type:bigint;primaryKey" json:"id"`
	FullName      string  `gorm:"type:varchar(255)" json:"full_name"`
	PhoneNumber   string  `gorm:"type:varchar(20);not null" json:"phone_number"`
	Address       string  `gorm:"type:varchar(255);not null" json:"address"`
	Commune       string  `gorm:"type:varchar(255);not null" json:"commune"`
	Province      string  `gorm:"type:varchar(255);not null" json:"province"`
	TotalPrice    float64 `gorm:"type:decimal(10,2);not null" json:"total_price"`
	TotalQuantity uint    `gorm:"type:int;not null" json:"total_quantity"`
	PaymentMethod string  `gorm:"type:enum('cod','bank','e-wallet');not null" json:"payment_method"`
	Status        string  `gorm:"type:enum('pending','confirmed','shipped');not null" json:"status"`

	OrderItems []*OrderItem `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"order_items"`
}

type OrderItem struct {
	ID         int64   `gorm:"type:bigint;primaryKey" json:"id"`
	UnitPrice  float64 `gorm:"type:decimal(10,2);not null" json:"unit_price"`
	Quantity   uint    `gorm:"type:int;not null" json:"quantity"`
	TotalPrice float64 `gorm:"type:decimal(10,2);not null" json:"total_price"`
	ProductID  int64   `gorm:"type:bigint;not null" json:"product_id"`
	OrderID    int64   `gorm:"type:bigint;not null" json:"order_id"`

	Order   *Order   `gorm:"foreignKey:OrderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"order"`
	Product *Product `gorm:"foreignKey:ProductID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"product"`
}
