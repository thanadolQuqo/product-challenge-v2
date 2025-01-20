package repository

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"product-challenge/internal/models"
	"product-challenge/pkg/config"
)

type OrderRepository interface {
	GetCart(ctx context.Context, username string) (*models.GetCartResponse, error)
	AddProductToCart(ctx context.Context, req *models.CartRequest) (bool, error)
	RemoveProductFromCart(ctx context.Context, req *models.CartRequest) (bool, error)
	MakeOrder(ctx context.Context, username string) (*models.Order, error)
}

type orderRepository struct {
	db          *gorm.DB
	cfg         config.Config
	orderNumber string
}

func NewOrderRepository(db *gorm.DB, cfg *config.Config) OrderRepository {
	return &orderRepository{
		db:  db,
		cfg: *cfg,
	}
}

func (r *orderRepository) GetCart(ctx context.Context, username string) (*models.GetCartResponse, error) {
	var (
		cart      models.Cart
		cartId    int
		cartItems []models.CartItem
		products  []models.Products
		totalCart float64
	)
	// get all item that store in cart with that specific username
	err := r.db.First(&cart, "owner = ? AND is_completed = ?", username, false).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("Find Cart by owner error : " + err.Error())
			return nil, err
		}
		fmt.Println("the cart is empty")
		return nil, nil
	}
	cartId = cart.ID

	// find all item in cart
	err = r.db.Where("cart_id = ?", cartId).Find(&cartItems).Error
	if err != nil {
		fmt.Println("Find CartItems error: " + err.Error())
		return nil, err
	}

	// Collect all product IDs from the cartItems
	productIDs := make([]int, 0)
	for _, item := range cartItems {
		productIDs = append(productIDs, item.ProductID)
	}

	// Query the product details for all product IDs
	err = r.db.Where("id IN ?", productIDs).Find(&products).Error
	if err != nil {
		fmt.Println("Find Products error: " + err.Error())
		return nil, err
	}

	// Store product data to map with ID as key
	productMap := make(map[int]models.Products)
	for _, product := range products {
		productMap[product.ID] = product
	}

	// Sum the total cost of the cart
	cartDetails := make([]models.GetCartItemResponse, 0)
	for _, item := range cartItems {
		productName := productMap[item.ProductID].Name
		productPrice := productMap[item.ProductID].Price
		subtotal := float64(item.Quantity) * productPrice
		totalCart += subtotal

		cartDetails = append(cartDetails, models.GetCartItemResponse{
			Name:     productName,
			Quantity: item.Quantity,
			Price:    productPrice,
			Subtotal: subtotal,
		})
	}

	// Prepare the response
	response := &models.GetCartResponse{
		CartId:    cartId,
		Owner:     username,
		CartItems: cartDetails,
	}

	return response, nil
}

func (r *orderRepository) AddProductToCart(ctx context.Context, req *models.CartRequest) (bool, error) {
	var (
		cart             models.Cart
		cartId           int
		existingCartItem models.CartItem
	)
	//	1. check if cart exist. if exist, get cart_id. If not, create new cart and get cart_id

	err := r.db.First(&cart, "owner = ? AND is_completed = ?", req.Username, false).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("Find Cart by owner error : " + err.Error())
			return false, err
		}
		newCart := models.Cart{
			Owner:       req.Username,
			IsCompleted: false,
		}

		err = r.db.Create(&newCart).Error
		if err != nil {
			fmt.Println("Create new Cart error : " + err.Error())
			return false, err
		}
		fmt.Println("newCart ID : ", newCart.ID)
		cartId = newCart.ID
	} else {
		cartId = cart.ID
	}

	// check if product Id is already in cart or not. If it is, add the request quantity to that product ID
	err = r.db.First(&existingCartItem, "cart_id = ? AND product_id = ?", cartId, req.ProductID).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("find cart items error : ", err)
			return false, err
		}

	}

	// If the product is already in the cart, update its quantity
	if existingCartItem.ID != 0 {
		existingCartItem.Quantity += req.Quantity
		err = r.db.Save(&existingCartItem).Error
		if err != nil {
			fmt.Println("Update CartItem Quantity error: " + err.Error())
			return false, err
		}
	} else {
		// If the product is not in the cart, create a new cart item
		newCartItem := models.CartItem{
			CartID:    cartId,
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
		}

		err = r.db.Create(&newCartItem).Error
		if err != nil {
			fmt.Println("Create CartItem error: " + err.Error())
			return false, err
		}
	}

	return true, nil
}

func (r *orderRepository) RemoveProductFromCart(ctx context.Context, req *models.CartRequest) (bool, error) {
	var (
		cart             models.Cart
		cartId           int
		existingCartItem models.CartItem
	)
	//	1. check if cart exist. if exist, get cart_id. If not, create new cart and get cart_id

	err := r.db.First(&cart, "owner = ? AND is_completed = ?", req.Username, false).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("Find Cart by owner error : " + err.Error())
			return false, err
		}
		newCart := models.Cart{
			Owner:       req.Username,
			IsCompleted: false,
		}

		err = r.db.Create(&newCart).Error
		if err != nil {
			fmt.Println("Create new Cart error : " + err.Error())
			return false, err
		}
		fmt.Println("newCart ID : ", newCart.ID)
		cartId = newCart.ID
	} else {
		cartId = cart.ID
	}

	// check if product Id is already in cart or not. If it is, add the request quantity to that product ID
	err = r.db.First(&existingCartItem, "cart_id = ? AND product_id = ?", cartId, req.ProductID).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("find cart items error : ", err)
			return false, err
		}
		fmt.Println("Product not found in cart")
		return false, nil
	}

	// If the product is already in the cart, update its quantity

	remainingQuantity := existingCartItem.Quantity - req.Quantity
	fmt.Println("remainingQuantity : ", remainingQuantity)
	if remainingQuantity <= 0 {
		fmt.Println("delete product from cart")
		err = r.db.Delete(&existingCartItem).Error
		if err != nil {
			fmt.Println("Delete CartItem error : " + err.Error())
			return false, err
		}
	} else {
		existingCartItem.Quantity -= req.Quantity
		err = r.db.Save(&existingCartItem).Error
		if err != nil {
			fmt.Println("Update CartItem Quantity error: " + err.Error())
			return false, err
		}
	}

	return true, nil
}

func (r *orderRepository) MakeOrder(ctx context.Context, username string) (*models.Order, error) {

	var (
		cart       models.Cart
		cartId     int
		cartItems  []models.CartItem
		products   []models.Products
		totalPrice float64
	)

	// 1. get cart id by username where it is not completed yet
	// 1.1 If cart is not available or item in cart is empty, return error
	err := r.db.First(&cart, "owner = ? AND is_completed = ?", username, false).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("Find Cart by owner error : " + err.Error())
			return nil, err
		}
		fmt.Println("the cart is empty. Can't make order")
		return nil, nil
	}

	cartId = cart.ID

	// find all items in cart
	err = r.db.Where("cart_id = ?", cartId).Find(&cartItems).Error
	if err != nil {
		fmt.Println("Find CartItems error: " + err.Error())
		return nil, err
	}

	if len(cartItems) == 0 {
		return nil, errors.New("cart is empty")
	}

	// Collect all product IDs from the cartItems
	productIDs := make([]int, 0)
	for _, item := range cartItems {
		productIDs = append(productIDs, item.ProductID)
	}

	// Query the product details for all product IDs
	err = r.db.Where("id IN ?", productIDs).Find(&products).Error
	if err != nil {
		fmt.Println("Find Products error: " + err.Error())
		return nil, err
	}

	// Create a map to easily look up product names by their ID
	productMap := make(map[int]models.Products)
	for _, product := range products {
		productMap[product.ID] = product
	}

	// get total price and check if stock is available
	for _, item := range cartItems {
		productPrice := productMap[item.ProductID].Price
		productStock := productMap[item.ProductID].Stock
		productName := productMap[item.ProductID].Name
		totalPrice += float64(item.Quantity) * productPrice

		if productStock < item.Quantity {
			fmt.Println("product stock is less than stock price")

			return nil, errors.New(fmt.Sprintf("product stock for %s is not enough", productName))
		}
	}

	//  Create a new Order
	orderNumber, err := r.GetLatestOrderNumber(r.db)
	if err != nil {
		fmt.Println("GetLatestOrderNumber error: " + err.Error())
		return nil, err
	}
	order := models.Order{
		OrderNumber: *orderNumber,
		Username:    username,
		TotalPrice:  totalPrice,
	}

	err = r.db.Create(&order).Error
	if err != nil {
		fmt.Println("Create Order error: " + err.Error())
		return nil, fmt.Errorf("error creating order: %v", err)
	}

	//  Create OrderItems for the new Order
	var orderItems []models.OrderItem
	for _, cartItem := range cartItems {
		productPrice := productMap[cartItem.ProductID].Price
		product := productMap[cartItem.ProductID]

		// Reduce the product stock
		product.Stock -= cartItem.Quantity
		err = r.db.Save(product).Error
		if err != nil {
			fmt.Printf("error updating stock for product %s (ID: %d): %v\n", product.Name, product.ID, err)
			return nil, fmt.Errorf("error updating stock for product %s (ID: %d): %v", product.Name, product.ID, err)
		}

		orderItems = append(orderItems, models.OrderItem{
			OrderID:   order.ID,
			ProductID: cartItem.ProductID,
			Price:     productPrice,
			Quantity:  cartItem.Quantity,
		})
	}

	// Save all order items
	err = r.db.Create(&orderItems).Error
	if err != nil {
		fmt.Println("Create OrderItems error: " + err.Error())
		return nil, fmt.Errorf("error creating order items: %v", err)
	}

	// Mark the cart as completed
	cart.IsCompleted = true
	err = r.db.Save(&cart).Error
	if err != nil {
		fmt.Println("Update Cart to IsComplete error: " + err.Error())
		return nil, fmt.Errorf("error updating cart: %v", err)
	}

	// Return the created order with items
	order.OrderItems = orderItems

	return &order, nil
}

// GetLatestOrderNumber is a function that generate the next order number.
func (r *orderRepository) GetLatestOrderNumber(tx *gorm.DB) (*string, error) {
	// Get next value from sequence
	var (
		nextVal     int
		orderNumber string
	)
	err := tx.Raw("SELECT nextval('order_number_seq')").Scan(&nextVal).Error
	if err != nil {
		return nil, err
	}

	// Format the order number (e.g., ORD-0001)
	orderNumber = fmt.Sprintf("ORD-%04d", nextVal)
	return &orderNumber, nil
}
