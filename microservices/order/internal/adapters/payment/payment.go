package payment

import (
	"context"
	"log"
	"time"

	pb "github.com/Mellanie-Marques/microservices-proto/golang/payment/payment"
	"github.com/Mellanie-Marques/microservices/order/internal/application/core/domain"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Adapter struct {
	client pb.PaymentClient // from the generated protobuf code
}

func NewAdapter(paymentServiceUrl string) (*Adapter, error) {
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
	conn, err := grpc.Dial(paymentServiceUrl, opts...)
	if err != nil {
		return nil, err
	}
	client := pb.NewPaymentClient(conn)
	return &Adapter{client: client}, nil
}

func (a *Adapter) Charge(order *domain.Order) error {
	// Criar contexto com timeout de 2 segundos
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := a.client.Create(ctx, &pb.CreatePaymentRequest{
		UserId:     order.CustomerID,
		OrderId:    order.ID,
		TotalPrice: order.TotalPrice(),
	})

	// Verificar se o erro foi por timeout
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.DeadlineExceeded {
				log.Printf("TIMEOUT: Falha ao processar pagamento - deadline excedido. OrderID: %d", order.ID)
			}
		}
	}

	return err
}
