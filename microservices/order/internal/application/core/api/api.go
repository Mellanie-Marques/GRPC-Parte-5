package api

import (
	"github.com/Mellanie-Marques/microservices/order/internal/application/core/domain"
	"github.com/Mellanie-Marques/microservices/order/internal/ports"
)

type Application struct {
	db       ports.DBPort
	payment  ports.PaymentPort
	product  ports.ProductPort
	shipping ports.ShippingPort
}

func NewApplication(db ports.DBPort, payment ports.PaymentPort, product ports.ProductPort, shipping ports.ShippingPort) *Application {
	return &Application{
		db:       db,
		payment:  payment,
		product:  product,
		shipping: shipping,
	}
}

func (a Application) PlaceOrder(order domain.Order) (domain.Order, error) {
	// Validar se todos os produtos existem
	for _, item := range order.OrderItems {
		exists, err := a.product.ValidateProduct(item.ProductCode)
		if err != nil {
			return domain.Order{}, err
		}
		if !exists {
			return domain.Order{}, domain.ErrProductNotFound
		}
	}

	// Salvar order
	err := a.db.Save(&order)
	if err != nil {
		return domain.Order{}, err
	}

	// Processar pagamento
	paymentErr := a.payment.Charge(&order)
	if paymentErr != nil {
		return domain.Order{}, paymentErr
	}

	// Solicitar entrega apenas se pagamento foi bem-sucedido
	deliveryDays, shippingErr := a.shipping.CalculateDelivery(order)
	if shippingErr != nil {
		return domain.Order{}, shippingErr
	}

	order.DeliveryDays = deliveryDays

	return order, nil
}
