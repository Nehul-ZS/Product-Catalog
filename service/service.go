package service

import (
	"errors"
	"exercises/Catalog/model"
	"exercises/Catalog/store/brand"
	"exercises/Catalog/store/product"
	"fmt"
	"log"
)

type catService struct {
	prodStore  product.StoreInterface
	brandStore brand.StoreInterface
}

//type category struct {
//	prodStore  store.Product
//	brandStore store.Brand
//}

func New(prodStore product.StoreInterface, brandStore brand.StoreInterface) ServInterface {
	return catService{
		prodStore:  prodStore,
		brandStore: brandStore,
	}
}

//Invalid Id
//type idCheck struct{
//	id int
//}
//func (id idCheck)Error() string{
//	return fmt.Sprintf("Invalid Id, Id got=%v",id.id)
//}
//Use return idCheck{id:1} for return error


func (C catService) GetById(id int) (model.Prod, error) {
	prodDet, err := C.prodStore.GetById(id)
	log.Println(err)
	if err != nil {
		log.Println(err)
		return model.Prod{}, errors.New("Id not found")
	}
	brandDet, err := C.brandStore.GetById(prodDet.BrandDetails.Id)
	log.Println(err)
	if err != nil {
		log.Println(err)
		return prodDet, errors.New("Brand Id not found")
	}
	prodDet.BrandDetails.Brand = brandDet.Brand
	return prodDet, nil
}
func (C catService) GetAll() ([]model.Prod, error) {
	prodDet, err := C.prodStore.GetAll()
	log.Println(err)
	if err != nil {
		log.Println(err)
		return []model.Prod{}, errors.New("Products not found")
	}
	for i,p:=range prodDet {
		brandDet, err := C.brandStore.GetById(p.BrandDetails.Id)
		log.Println(err)
		if err != nil {
			log.Println(err)
			return prodDet, errors.New("Brand Id not found")
		}
		prodDet[i].BrandDetails.Brand = brandDet.Brand
	}
	return prodDet, nil
}

func (C catService) Create(empName, bName string) (model.Prod, error) {
	bId, err := C.brandStore.CheckBrand(bName)
	if err != nil || bId == 0 {
		res, err := C.brandStore.Create(bName)
		if err != nil {
			return model.Prod{}, errors.New("Brand not created")
		}
		bId = res
	}
	res, err := C.prodStore.Create(empName, bId)
	if err != nil {
		return model.Prod{}, errors.New("Product not created")
	}
	prodDet, err := C.prodStore.GetById(res)
	if err != nil {
		return model.Prod{}, errors.New("Created Product not found")
	}
	brandDet, err := C.brandStore.GetById(prodDet.BrandDetails.Id)
	log.Println(err)
	if err != nil {
		log.Println(err)
		return prodDet, errors.New("Brand not found")
	}
	prodDet.BrandDetails.Brand = brandDet.Brand
	return prodDet, nil
}
func (C catService) Update(empId int, empName, bName string) (model.Prod, error) {
	bId, err := C.brandStore.CheckBrand(bName)
	if err != nil || bId == 0 {
		res, err := C.brandStore.Create(bName)
		if err != nil {
			return model.Prod{}, err
		}
		bId = res
	}
	res, err := C.prodStore.Update(empId,empName, bId)
	if err != nil {
		return model.Prod{}, errors.New("Product not updated")
	}
	fmt.Println(res)
	prodDet, err := C.prodStore.GetById(empId)
	if err != nil {
		return model.Prod{}, errors.New("Updated Product not found")
	}
	brandDet, err := C.brandStore.GetById(prodDet.BrandDetails.Id)
	log.Println(err)
	if err != nil {
		log.Println(err)
		return prodDet, errors.New("Updated Brand not found")
	}
	prodDet.BrandDetails.Brand = brandDet.Brand
	return prodDet, nil
}
func (C catService) Delete(empId int) error {
	res, err := C.prodStore.Delete(empId)
	if err != nil || res != 1 {
		return  errors.New("Product not deleted from database")
	}
	return nil
}
