package services

import (
	pb "cart/proto"
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
)

type ProductService struct {
	client pb.ProductServiceClient
}

func NewProductService(conn *grpc.ClientConn) *ProductService {
	return &ProductService{
		client: pb.NewProductServiceClient(conn),
	}
}

// ValidateInventory checks if product exists and there is sufficient inventory for the requested quantity
func (pc *ProductService) ValidateInventory(productID, requestedQuantity uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := pc.client.ValidateProductInventory(ctx, &pb.ValidateProductInventoryRequest{
		ProductId: productID,
		Quantity:  requestedQuantity,
	})

	if !resp.GetValid() || err != nil {
		return fmt.Errorf("failed to validate inventory: %w", err)
	}

	return nil
}
