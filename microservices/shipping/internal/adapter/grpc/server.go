package grpc

import (
	"context"
	"log"
	"net"

	pb "github.com/Mellanie-Marques/microservices-proto/golang/shipping"
	"github.com/Mellanie-Marques/microservices/shipping/internal/application/core/api"
	"github.com/Mellanie-Marques/microservices/shipping/internal/application/core/domain"
	"google.golang.org/grpc"
)

type Adapter struct {
	app  *api.Application
	port string
	pb.UnimplementedShippingServer
}

func NewAdapter(app *api.Application, port string) *Adapter {
	return &Adapter{
		app:  app,
		port: port,
	}
}

func (a *Adapter) Create(ctx context.Context, req *pb.CreateShippingRequest) (*pb.CreateShippingResponse, error) {
	// Converter protobuf para domínio
	items := make([]domain.ShippingItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = domain.ShippingItem{
			ProductCode: item.ProductCode,
			Quantity:    item.Quantity,
		}
	}

	shipping := &domain.Shipping{
		OrderID: req.OrderId,
		Items:   items,
	}

	// Calcular prazo de entrega
	deliveryDays, err := a.app.CreateShipping(shipping)
	if err != nil {
		return nil, err
	}

	return &pb.CreateShippingResponse{
		DeliveryDays: deliveryDays,
	}, nil
}

func (a *Adapter) Run() {
	listener, err := net.Listen("tcp", ":"+a.port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterShippingServer(grpcServer, a)

	log.Printf("gRPC server starting and listening on port %s", a.port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
