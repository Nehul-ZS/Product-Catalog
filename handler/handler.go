package handler

import (
	"encoding/json"
	"exercises/Catalog/model"
	"exercises/Catalog/service"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type catHandler struct {
	catServ service.ServInterface
}

func New(catServ service.ServInterface) catHandler {
	return catHandler{
		catServ: catServ,
	}
}
func (c catHandler) Get(w http.ResponseWriter,r *http.Request){
	query:=r.URL.Query()
	if len(query)==0 {
		c.GetAll(w,r)
		return
	}
	c.GetById(w,r)
	return
}
func (c catHandler) GetById(w http.ResponseWriter, r *http.Request) {
	res := r.URL.Query()["id"][0]
	key, err := strconv.Atoi(res)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid ID"))
		return
	}
	prodDet, err := c.catServ.GetById(key)
	jsonRes, _ := json.Marshal(prodDet)
	//log.Println(err)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsonRes)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonRes)
}
func (c catHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	prodDet, err := c.catServ.GetAll()
	jsonRes, _ := json.Marshal(prodDet)
	//log.Println(err)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsonRes)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonRes)
}

func (c catHandler) Create(w http.ResponseWriter, r *http.Request) {
	var pD model.Prod
	body := r.Body
	err := json.NewDecoder(body).Decode(&pD)
	fmt.Println(err)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Product not added to database"))
		log.Println(err)
		return
	}
	res, err := c.catServ.Create(pD.Name, pD.BrandDetails.Brand)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Product not added to database"))
		log.Println(err)
		return
	}
	jsonRes, _ := json.Marshal(res)
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonRes)

}
func (c catHandler) Update(w http.ResponseWriter, r *http.Request) {
	var pD model.Prod
	key, err := strconv.Atoi(r.URL.Query()["id"][0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid ID"))
		return
	}
	body := r.Body
	err = json.NewDecoder(body).Decode(&pD)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Product not updated in database"))
		log.Println(err)
		return
	}
	res, err := c.catServ.Update(key,pD.Name, pD.BrandDetails.Brand)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Product not updated in database"))
		log.Println(err)
		return
	}
	jsonRes, _ := json.Marshal(res)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonRes)

}
func (c catHandler) Delete(w http.ResponseWriter, r *http.Request) {
	key, err := strconv.Atoi(r.URL.Query()["id"][0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid ID"))
		return
	}
	resErr := c.catServ.Delete(key)

	if resErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Product not deleted in database"))
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
