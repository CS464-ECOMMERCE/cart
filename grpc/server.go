package grpc

import (
	"cart/configs"
	pb "cart/proto"
	"cart/services"
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// Init initializes and starts the gRPC server
func Init() {
	config := configs.GetEnvConfig()
	address := fmt.Sprintf(":%s", config.GRPCPort)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthServer)
	healthServer.SetServingStatus("ProductService", grpc_health_v1.HealthCheckResponse_SERVING)
	pb.RegisterCartServiceServer(s, newCartServer())

	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// cartServer is the gRPC server implementation
type cartServer struct {
	cartService *services.CartService
	pb.UnimplementedCartServiceServer
}

func newCartServer() *cartServer {
	return &cartServer{
		cartService: services.NewCartService(),
	}
}

// AddItem implements the AddItem RPC method
func (s *cartServer) AddItem(ctx context.Context, req *pb.AddItemRequest) (*pb.Empty, error) {
	item := services.CartItem{
		ProductID: req.Item.ProductId,
		Quantity:  req.Item.Quantity,
	}

	err := s.cartService.AddItem(req.SessionId, item)
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

// GetCart implements the GetCart RPC method
func (s *cartServer) GetCart(ctx context.Context, req *pb.GetCartRequest) (*pb.Cart, error) {
	if req.SessionId == "" {
		return nil, fmt.Errorf("session ID is required")
	}
	cart, err := s.cartService.GetCart(req.SessionId)
	if err != nil {
		return nil, err
	}

	// Convert the domain cart to the protobuf cart
	result := &pb.Cart{
		SessionId: cart.SessionId,
		Items:     make([]*pb.CartItem, 0, len(cart.Items)),
	}

	for _, item := range cart.Items {
		result.Items = append(result.Items, &pb.CartItem{
			ProductId: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	return result, nil
}

// EmptyCart implements the EmptyCart RPC method
func (s *cartServer) EmptyCart(ctx context.Context, req *pb.EmptyCartRequest) (*pb.Empty, error) {
	err := s.cartService.EmptyCart(req.SessionId)
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

// RemoveItem implements the RemoveItem RPC method
func (s *cartServer) RemoveItem(ctx context.Context, req *pb.RemoveItemRequest) (*pb.Empty, error) {
	err := s.cartService.RemoveItem(req.SessionId, req.ProductId)
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

// UpdateItemQuantity implements the UpdateItemQuantity RPC method
func (s *cartServer) UpdateItemQuantity(ctx context.Context, req *pb.UpdateItemQuantityRequest) (*pb.Empty, error) {
	err := s.cartService.UpdateItemQuantity(req.SessionId, req.ProductId, req.Quantity)
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}
