package product

import "exercises/Catalog/model"

type StoreInterface interface {
	GetById(int) (model.Prod, error)
	GetAll()([]model.Prod,error)
	Create(string, int) (int, error)
	Update(int,string,int)(int,error)
	Delete(int)(int,error)
}
