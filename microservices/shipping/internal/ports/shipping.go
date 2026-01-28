package ports

import "github.com/Mellanie-Marques/microservices/shipping/internal/application/core/domain"

type ShippingPort interface {
	Create(shipping *domain.Shipping) (int32, error)
}
