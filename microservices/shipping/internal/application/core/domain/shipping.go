package domain

type ShippingItem struct {
	ProductCode string
	Quantity    int32
}

type Shipping struct {
	OrderID      int32
	Items        []ShippingItem
	DeliveryDays int32
}

func (s *Shipping) CalculateDeliveryDays() int32 {
	// Calcular quantidade total de unidades
	totalQuantity := int32(0)
	for _, item := range s.Items {
		totalQuantity += item.Quantity
	}

	// Mínimo 1 dia + 1 dia a cada 5 unidades
	deliveryDays := int32(1)
	if totalQuantity > 5 {
		deliveryDays += (totalQuantity - 1) / 5
	}

	return deliveryDays
}
