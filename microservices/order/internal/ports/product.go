package ports

type ProductPort interface {
	ValidateProduct(productCode string) (bool, error)
}
