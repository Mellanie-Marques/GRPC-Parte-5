package ports

import "github.com/Mellanie-Marques/microservices/order/internal/application/core/domain"

type ShippingPort interface {
	CalculateDelivery(order domain.Order) (int32, error)
}
