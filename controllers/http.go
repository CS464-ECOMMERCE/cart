package controllers

import (
	"cart/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CartController handles HTTP requests for the cart service
type CartController struct {
	cartService *services.CartService
}

// NewCartController creates a new cart controller
func NewCartController() *CartController {
	return &CartController{
		cartService: services.NewCartService(),
	}
}

// SetupRoutes sets up the routes for the cart controller
func (c *CartController) SetupRoutes(router *gin.Engine) {
	cartGroup := router.Group("/cart")
	{
		cartGroup.GET("/:user_id", c.GetCart)
		cartGroup.POST("/:user_id/items", c.AddItem)
		cartGroup.DELETE("/:user_id/items/:product_id", c.RemoveItem)
		cartGroup.PUT("/:user_id/items/:product_id", c.UpdateItemQuantity)
		cartGroup.DELETE("/:user_id", c.EmptyCart)
		cartGroup.POST("/:user_id/merge/:guest_user_id", c.MergeCart)
		cartGroup.GET("/:user_id/ttl", c.GetCartTTL)
		cartGroup.PUT("/:user_id/ttl", c.ExtendCartTTL)
	}
}

// GetCart handles GET /cart/:user_id
func (c *CartController) GetCart(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("user_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	cart, err := c.cartService.GetCart(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, cart)
}

// AddItem handles POST /cart/:user_id/items
func (c *CartController) AddItem(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("user_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var item services.CartItem
	if err := ctx.ShouldBindJSON(&item); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.cartService.AddItem(userID, item)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the updated cart
	cart, err := c.cartService.GetCart(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, cart)
}

// RemoveItem handles DELETE /cart/:user_id/items/:product_id
func (c *CartController) RemoveItem(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("user_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	productID, err := strconv.ParseUint(ctx.Param("product_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	err = c.cartService.RemoveItem(userID, productID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the updated cart
	cart, err := c.cartService.GetCart(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, cart)
}

// UpdateItemQuantity handles PUT /cart/:user_id/items/:product_id
func (c *CartController) UpdateItemQuantity(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("user_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	productID, err := strconv.ParseUint(ctx.Param("product_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var request struct {
		Quantity uint64 `json:"quantity" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.cartService.UpdateItemQuantity(userID, productID, request.Quantity)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the updated cart
	cart, err := c.cartService.GetCart(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, cart)
}

// EmptyCart handles DELETE /cart/:user_id
func (c *CartController) EmptyCart(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("user_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = c.cartService.EmptyCart(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Cart emptied successfully"})
}

// MergeCart handles POST /cart/:user_id/merge/:guest_user_id
func (c *CartController) MergeCart(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("user_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	guestUserID, err := strconv.ParseUint(ctx.Param("guest_user_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid guest user ID"})
		return
	}

	cart, err := c.cartService.MergeCart(userID, guestUserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, cart)
}

// GetCartTTL handles GET /cart/:user_id/ttl
func (c *CartController) GetCartTTL(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("user_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	ttl, err := c.cartService.GetCartTTL(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"ttl_seconds": ttl})
}

// ExtendCartTTL handles PUT /cart/:user_id/ttl
func (c *CartController) ExtendCartTTL(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("user_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var request struct {
		TTLSeconds int64 `json:"ttl_seconds" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ttl, err := c.cartService.ExtendCartTTL(userID, request.TTLSeconds)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"ttl_seconds": ttl})
}
