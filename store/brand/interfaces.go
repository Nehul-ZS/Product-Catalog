package brand

import "exercises/Catalog/model"

type StoreInterface interface {
	GetById(int) (model.Brand, error)
	GetAll()([]model.Brand,error)
	Create(string) (int, error)
	CheckBrand(string) (int, error)
}
