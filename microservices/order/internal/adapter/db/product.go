package db

import (
	"fmt"

	"gorm.io/gorm"
)

type Product struct {
	ID          uint
	ProductCode string `gorm:"uniqueIndex"`
	Name        string
	Quantity    int
}

func (a *Adapter) ValidateProduct(productCode string) (bool, error) {
	var product Product
	if err := a.db.Where("product_code = ?", productCode).First(&product).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (a *Adapter) GetProduct(productCode string) (*Product, error) {
	var product Product
	if err := a.db.Where("product_code = ?", productCode).First(&product).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("produto '%s' não encontrado", productCode)
		}
		return nil, err
	}
	return &product, nil
}
