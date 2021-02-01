package main

import (
	"database/sql"
	"exercises/Catalog/handler"
	"exercises/Catalog/service"
	"exercises/Catalog/store/brand"
	"exercises/Catalog/store/product"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	db, err := sql.Open("mysql", "nehul:9618181838@tcp(127.0.0.1)/catalogue")
	defer db.Close()
	if err != nil {
		log.Println(err)
	}
	pS := product.New(db)
	bS := brand.New(db)
	catServ := service.New(pS, bS)
	catHandle := handler.New(catServ)
	myRouter := mux.NewRouter()
	myRouter.HandleFunc("/product", catHandle.Get).Methods("GET")
	myRouter.HandleFunc("/product", catHandle.Create).Methods("POST")
	myRouter.HandleFunc("/product", catHandle.Update).Methods("PUT")
	myRouter.HandleFunc("/product", catHandle.Delete).Methods("DELETE")
	err = http.ListenAndServe(":8080", myRouter)
	if err != nil {
		log.Println(err)
	}
}
