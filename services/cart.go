package services

import (
	"encoding/json"
	"fmt"
	"time"
)

// CartItem represents an item in the cart
type CartItem struct {
	ProductID uint64 `json:"product_id"`
	Quantity  uint64 `json:"quantity"`
}

// Cart represents a user's shopping cart
type Cart struct {
	UserID uint64     `json:"user_id"`
	Items  []CartItem `json:"items"`
}

// CartService provides operations for manipulating carts
type CartService struct {
	redis *RedisClient
}

// NewCartService creates a new cart service instance
func NewCartService() *CartService {
	return &CartService{
		redis: GetRedisClient(),
	}
}

// GetCart retrieves a user's cart
func (s *CartService) GetCart(userID uint64) (*Cart, error) {
	key := s.redis.getCartKey(userID)
	data, err := s.redis.client.Get(s.redis.ctx, key).Bytes()
	if err != nil {
		if err.Error() == "redis: nil" {
			// Return empty cart if not found
			return &Cart{
				UserID: userID,
				Items:  []CartItem{},
			}, nil
		}
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	// Extend TTL on cart access
	s.redis.ExtendCartTTL(userID)

	var cart Cart
	if err := json.Unmarshal(data, &cart); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cart: %w", err)
	}

	return &cart, nil
}

// AddItem adds an item to a user's cart
func (s *CartService) AddItem(userID uint64, item CartItem) error {
	cart, err := s.GetCart(userID)
	if err != nil {
		return err
	}

	// Check if item already exists in cart
	for i, cartItem := range cart.Items {
		if cartItem.ProductID == item.ProductID {
			// Update quantity
			cart.Items[i].Quantity += item.Quantity
			return s.saveCart(cart)
		}
	}

	// Add new item
	cart.Items = append(cart.Items, item)
	return s.saveCart(cart)
}

// RemoveItem removes an item from a user's cart
func (s *CartService) RemoveItem(userID uint64, productID uint64) error {
	cart, err := s.GetCart(userID)
	if err != nil {
		return err
	}

	// Find and remove item
	for i, item := range cart.Items {
		if item.ProductID == productID {
			// Remove item by replacing with last element and truncating
			cart.Items[i] = cart.Items[len(cart.Items)-1]
			cart.Items = cart.Items[:len(cart.Items)-1]
			return s.saveCart(cart)
		}
	}

	return nil // Item not found, no action needed
}

// UpdateItemQuantity updates the quantity of an item in a user's cart
func (s *CartService) UpdateItemQuantity(userID uint64, productID uint64, quantity uint64) error {
	cart, err := s.GetCart(userID)
	if err != nil {
		return err
	}

	// Find and update item
	for i, item := range cart.Items {
		if item.ProductID == productID {
			// Remove item if quantity is 0
			if quantity == 0 {
				return s.RemoveItem(userID, productID)
			}

			// Update quantity
			cart.Items[i].Quantity = quantity
			return s.saveCart(cart)
		}
	}

	// If item not found and quantity > 0, add it
	if quantity > 0 {
		return s.AddItem(userID, CartItem{
			ProductID: productID,
			Quantity:  quantity,
		})
	}

	return nil // Item not found and quantity is 0, no action needed
}

// EmptyCart removes all items from a user's cart
func (s *CartService) EmptyCart(userID uint64) error {
	return s.redis.DeleteCart(userID)
}

// MergeCart merges a guest cart into a user's cart
func (s *CartService) MergeCart(userID uint64, guestUserID uint64) (*Cart, error) {
	// Get both carts
	userCart, err := s.GetCart(userID)
	if err != nil {
		return nil, err
	}

	guestCart, err := s.GetCart(guestUserID)
	if err != nil {
		return nil, err
	}

	// If guest cart is empty, nothing to merge
	if len(guestCart.Items) == 0 {
		return userCart, nil
	}

	// Merge items
	for _, guestItem := range guestCart.Items {
		found := false

		// Look for matching item in user cart
		for i, userItem := range userCart.Items {
			if userItem.ProductID == guestItem.ProductID {
				// Update quantity
				userCart.Items[i].Quantity += guestItem.Quantity
				found = true
				break
			}
		}

		// If not found, add to user cart
		if !found {
			userCart.Items = append(userCart.Items, guestItem)
		}
	}

	// Save user cart
	if err := s.saveCart(userCart); err != nil {
		return nil, err
	}

	// Empty guest cart
	if err := s.EmptyCart(guestUserID); err != nil {
		return nil, err
	}

	return userCart, nil
}

// GetCartTTL gets the TTL for a user's cart
func (s *CartService) GetCartTTL(userID uint64) (int64, error) {
	ttl, err := s.redis.GetCartTTL(userID)
	if err != nil {
		return 0, err
	}
	return int64(ttl.Seconds()), nil
}

// ExtendCartTTL extends the TTL for a user's cart
func (s *CartService) ExtendCartTTL(userID uint64, ttlSeconds int64) (int64, error) {
	if err := s.redis.SetCartTTL(userID, time.Duration(ttlSeconds)*time.Second); err != nil {
		return 0, err
	}
	return ttlSeconds, nil
}

// saveCart saves a cart to Redis
func (s *CartService) saveCart(cart *Cart) error {
	data, err := json.Marshal(cart)
	if err != nil {
		return fmt.Errorf("failed to marshal cart: %w", err)
	}

	key := s.redis.getCartKey(cart.UserID)
	if err := s.redis.client.Set(s.redis.ctx, key, data, s.redis.config.RedisDefaultTTL).Err(); err != nil {
		return fmt.Errorf("failed to save cart: %w", err)
	}

	return nil
}
