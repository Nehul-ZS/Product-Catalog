package store

import "exercises/Catalog/model"

type Brand interface {
	GetById(int) (model.Prod, error)
}

type Product interface {
	GetById(int) (model.Prod, error)
}
