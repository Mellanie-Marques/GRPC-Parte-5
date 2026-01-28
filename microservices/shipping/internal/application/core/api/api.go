package api

import "github.com/Mellanie-Marques/microservices/shipping/internal/application/core/domain"

type Application struct {
}

func NewApplication() *Application {
	return &Application{}
}

func (a *Application) CreateShipping(shipping *domain.Shipping) (int32, error) {
	// Calcular prazo de entrega
	deliveryDays := shipping.CalculateDeliveryDays()
	return deliveryDays, nil
}
