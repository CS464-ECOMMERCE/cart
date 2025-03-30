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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
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
	if req.Item.Id == 0 || req.Item.Quantity <= 0 {
		return nil, status.Error(codes.InvalidArgument, "Id or Quantity is invalid")
	}
	item := services.CartItem{
		Id:       req.Item.Id,
		Quantity: req.Item.Quantity,
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
		return nil, status.Error(codes.InvalidArgument, "session ID is required")
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
			Id:       item.Id,
			Quantity: item.Quantity,
		})
	}

	return result, nil
}

// EmptyCart implements the EmptyCart RPC method
func (s *cartServer) EmptyCart(ctx context.Context, req *pb.EmptyCartRequest) (*pb.Empty, error) {
	err := s.cartService.EmptyCart(req.SessionId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Something went wrong. Unable to empty cart. %v", err.Error())
	}

	return &pb.Empty{}, nil
}

// RemoveItem implements the RemoveItem RPC method
func (s *cartServer) RemoveItem(ctx context.Context, req *pb.RemoveItemRequest) (*pb.Empty, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "Id is invalid")
	}
	err := s.cartService.RemoveItem(req.SessionId, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Something went wrong. Unable to remove item. %v", err.Error())
	}

	return &pb.Empty{}, nil
}

// UpdateItemQuantity implements the UpdateItemQuantity RPC method
func (s *cartServer) UpdateItemQuantity(ctx context.Context, req *pb.UpdateItemQuantityRequest) (*pb.Empty, error) {
	if req.Id == 0 || req.Quantity <= 0 {
		return nil, status.Error(codes.InvalidArgument, "Id or Quantity is invalid")
	}
	err := s.cartService.UpdateItemQuantity(req.SessionId, req.Id, req.Quantity)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Something went wrong. Unable to update quantity. %v", err.Error())
	}

	return &pb.Empty{}, nil
}
