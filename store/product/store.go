package product

import (
	"database/sql"
	"exercises/Catalog/model"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type prodStore struct {
	DB *sql.DB
}

func New(db *sql.DB) StoreInterface {
	return &prodStore{
		DB: db,
	}
}
func (pS prodStore) GetById(key int) (model.Prod, error) {
	var p model.Prod
	var res *sql.Rows
	var err error
	//var prodList model.Prod
	//tx,err:=pS.DB.Begin()
	//defer pS.DB.Close()
	//if err != nil {
	//	log.Println(err)
	//	return model.Prod{},err
	//}
	res, err = pS.DB.Query("select P.id,P.name,P.brand from product as P where id=?", key)
	if err != nil {
		log.Println(err)
		//return model.Prod{}, errors.New("Product Id Not Found")
		return model.Prod{}, err
	}
	defer res.Close()
	for res.Next() {
		err := res.Scan(&p.Id, &p.Name, &p.BrandDetails.Id)
		if err != nil {
			log.Println(err)
			return model.Prod{}, err
		}
	}
	//prodList.Id=p.Id
	//prodList.Name=p.Name
	//prodList.BrandDetails.Id=p.BrandDetails.Id
	//fmt.Printf("%q",p.BrandDetails.Brand)
	return p, nil
}
func (pS prodStore) GetAll() ([]model.Prod, error) {
	var p model.Prod
	var res *sql.Rows
	var err error
	prodDetails := make([]model.Prod, 0)
	res, err = pS.DB.Query("select P.id,P.name,P.brand from product as P")
	if err != nil {
		log.Println(err)
		return []model.Prod{}, err
		//return []model.Prod{}, errors.New("Products Not Found")
	}
	defer res.Close()
	for res.Next() {
		err := res.Scan(&p.Id, &p.Name, &p.BrandDetails.Id)
		if err != nil {
			log.Println(err)
			return []model.Prod{}, err
			//return []model.Prod{}, errors.New("Products Not Found")
		}
		prodDetails = append(prodDetails, p)
	}
	return prodDetails, nil
}
func (pS prodStore) Create(name string, brandId int) (int, error) {
	//dataChanges:=make([]int,0)
	res, err := pS.DB.Exec("insert into product(name,brand) values(?,?)", name, brandId)
	if err != nil {
		log.Println(err)
		return 0, err
		//return 0, errors.New("Product not inserted into database")
	}
	lastId, err := res.LastInsertId()
	if err != nil || lastId <= 0 {
		log.Println(err)
		return 0, err
	}
	return int(lastId), nil
}
func (pS prodStore) Update(key int,name string, brandId int) (int, error) {
	//dataChanges:=make([]int,0)
	res, err := pS.DB.Exec("update product set name=?,brand=? where id=?", name, brandId,key)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	lastId, err := res.LastInsertId()
	if err != nil || lastId <= 0 {
		log.Println(err)
		return 0, err
	}
	return int(lastId), nil
}
func (pS prodStore) Delete(key int) (int, error) {
	//dataChanges:=make([]int,0)
	res, err := pS.DB.Exec("delete from product where id=?", key)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	rows, err := res.RowsAffected()
	if err != nil || rows != 1 {
		log.Println(err)
		return 0, err
	}
	return int(rows), nil
}