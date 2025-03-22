package grpc

import (
	"cart/configs"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Init initializes and starts the gRPC server
func Init() {
	config := configs.GetEnvConfig()
	address := fmt.Sprintf(":%s", config.GRPCPort)

	// Create TCP listener
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create gRPC server
	srv := grpc.NewServer()

	// We need to properly generate the protobuf code first
	// Using the server in a simplified form to avoid errors with undefined types

	// Register reflection service for debugging
	reflection.Register(srv)

	log.Printf("Starting gRPC server on %s", address)

	// Start serving
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

/*
// cartServer is the gRPC server implementation
type cartServer struct {
	cartService *services.CartService
	UnimplementedCartServiceServer
}

// AddItem implements the AddItem RPC method
func (s *cartServer) AddItem(ctx context.Context, req *AddItemRequest) (*Empty, error) {
	item := services.CartItem{
		ProductID: req.Item.ProductId,
		Quantity:  req.Item.Quantity,
	}

	err := s.cartService.AddItem(req.UserId, item)
	if err != nil {
		return nil, err
	}

	return &Empty{}, nil
}

// GetCart implements the GetCart RPC method
func (s *cartServer) GetCart(ctx context.Context, req *GetCartRequest) (*Cart, error) {
	cart, err := s.cartService.GetCart(req.UserId)
	if err != nil {
		return nil, err
	}

	// Convert the domain cart to the protobuf cart
	result := &Cart{
		UserId: cart.UserID,
		Items:  make([]*CartItem, 0, len(cart.Items)),
	}

	for _, item := range cart.Items {
		result.Items = append(result.Items, &CartItem{
			ProductId: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	return result, nil
}

// EmptyCart implements the EmptyCart RPC method
func (s *cartServer) EmptyCart(ctx context.Context, req *EmptyCartRequest) (*Empty, error) {
	err := s.cartService.EmptyCart(req.UserId)
	if err != nil {
		return nil, err
	}

	return &Empty{}, nil
}

// RemoveItem implements the RemoveItem RPC method
func (s *cartServer) RemoveItem(ctx context.Context, req *RemoveItemRequest) (*Empty, error) {
	err := s.cartService.RemoveItem(req.UserId, req.ProductId)
	if err != nil {
		return nil, err
	}

	return &Empty{}, nil
}

// UpdateItemQuantity implements the UpdateItemQuantity RPC method
func (s *cartServer) UpdateItemQuantity(ctx context.Context, req *UpdateItemQuantityRequest) (*Empty, error) {
	err := s.cartService.UpdateItemQuantity(req.UserId, req.ProductId, req.Quantity)
	if err != nil {
		return nil, err
	}

	return &Empty{}, nil
}

// MergeCart implements the MergeCart RPC method
func (s *cartServer) MergeCart(ctx context.Context, req *MergeCartRequest) (*Cart, error) {
	cart, err := s.cartService.MergeCart(req.UserId, req.GuestUserId)
	if err != nil {
		return nil, err
	}

	// Convert the domain cart to the protobuf cart
	result := &Cart{
		UserId: cart.UserID,
		Items:  make([]*CartItem, 0, len(cart.Items)),
	}

	for _, item := range cart.Items {
		result.Items = append(result.Items, &CartItem{
			ProductId: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	return result, nil
}

// GetCartTTL implements the GetCartTTL RPC method
func (s *cartServer) GetCartTTL(ctx context.Context, req *GetCartRequest) (*CartTTL, error) {
	ttl, err := s.cartService.GetCartTTL(req.UserId)
	if err != nil {
		return nil, err
	}

	return &CartTTL{
		UserId:    req.UserId,
		TtlSeconds: ttl,
	}, nil
}

// ExtendCartTTL implements the ExtendCartTTL RPC method
func (s *cartServer) ExtendCartTTL(ctx context.Context, req *ExtendCartTTLRequest) (*CartTTL, error) {
	ttl, err := s.cartService.ExtendCartTTL(req.UserId, req.TtlSeconds)
	if err != nil {
		return nil, err
	}

	return &CartTTL{
		UserId:    req.UserId,
		TtlSeconds: ttl,
	}, nil
}
*/
