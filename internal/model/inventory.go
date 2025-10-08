package model

type Inventory struct {
	ID        int64 `gorm:"type:int;primaryKey" json:"id"`
	Quantity  uint  `gorm:"type:int;not null" json:"quantity"`
	Purchased uint  `gorm:"type:int;not null" json:"purchased"`
	Stock     uint  `gorm:"type:int;not null" json:"stock"`
	IsStock   bool  `gorm:"type:boolean;not null" json:"is_stock"`
	ProductID int64 `gorm:"type:int;not null;unique" json:"product_id"`

	Product *Product `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"product"`
}

func(m *Inventory) SetStock() {
	m.Stock = m.Quantity - m.Purchased
	if m.Stock <= 5 {
		m.IsStock = false
	}
}
