package models

import "time"

type CartRequest struct {
	Username  string `json:"username"`
	ProductID int    `json:"productId" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required"`
}
type Cart struct {
	ID          int        `json:"id" gorm:"primaryKey"`
	Owner       string     `json:"owner"`
	IsCompleted bool       `json:"isCompleted" gorm:"default:false"`
	CartItems   []CartItem `json:"cartItems" gorm:"foreignKey:CartID"`
}

type CartItem struct {
	ID        int  `json:"id" gorm:"primaryKey"`
	CartID    int  `json:"cartId" gorm:"not null"`
	ProductID int  `json:"productId"`
	Quantity  int  `json:"quantity"`
	Cart      Cart `json:"cart" gorm:"foreignKey:CartID"`
}

type GetCartResponse struct {
	CartId    int                   `json:"cartId"`
	Owner     string                `json:"owner"`
	CartItems []GetCartItemResponse `json:"cartItems"`
}

type GetCartItemResponse struct {
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
	Subtotal float64 `json:"subtotal"`
}

type Order struct {
	ID          int         `json:"id" gorm:"primaryKey"`
	OrderNumber string      `json:"orderNumber" gorm:"unique;not null"`
	Username    string      `json:"username"`
	TotalPrice  float64     `json:"totalPrice"`
	CreatedAt   time.Time   `json:"createdAt" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time   `json:"updatedAt" gorm:"default:CURRENT_TIMESTAMP"`
	OrderItems  []OrderItem `json:"orderItems,omitempty" gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	ID        int     `json:"id" gorm:"primaryKey"`
	OrderID   int     `json:"orderId" gorm:"not null"`
	ProductID int     `json:"productId"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
	Order     Order   `json:"order,omitempty" gorm:"foreignKey:OrderID"`
}
