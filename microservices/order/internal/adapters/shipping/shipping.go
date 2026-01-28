package shipping

import (
	"context"
	"log"
	"time"

	shippingpb "github.com/Mellanie-Marques/microservices-proto/golang/shipping"
	"github.com/Mellanie-Marques/microservices/order/internal/application/core/domain"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Adapter struct {
	client shippingpb.ShippingClient
}

func NewAdapter(shippingServiceUrl string) (*Adapter, error) {
	var opts []grpc.DialOption

	// Configurar interceptor de retry automático
	opts = append(opts,
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(
			grpc_retry.WithCodes(codes.Unavailable, codes.ResourceExhausted),
			grpc_retry.WithMax(5),
			grpc_retry.WithBackoff(grpc_retry.BackoffLinear(time.Second)),
		)),
	)

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(shippingServiceUrl, opts...)
	if err != nil {
		return nil, err
	}
	client := shippingpb.NewShippingClient(conn)
	return &Adapter{client: client}, nil
}

func (a *Adapter) CalculateDelivery(order domain.Order) (int32, error) {
	// Criar contexto com timeout de 2 segundos
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Converter itens de domínio para protobuf
	shippingItems := make([]*shippingpb.ShippingItem, len(order.OrderItems))
	for i, item := range order.OrderItems {
		shippingItems[i] = &shippingpb.ShippingItem{
			ProductCode: item.ProductCode,
			Quantity:    item.Quantity,
		}
	}

	requisicao := &shippingpb.CreateShippingRequest{
		OrderId: int32(order.ID),
		Items:   shippingItems,
	}

	resposta, err := a.client.Create(ctx, requisicao)

	// Verificar se o erro foi por timeout
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.DeadlineExceeded {
				log.Printf("TIMEOUT: Falha ao processar envio - deadline excedido. OrderID: %d", order.ID)
			}
		}
	}

	if err != nil {
		return 0, err
	}

	return resposta.DeliveryDays, nil
}
