package handler

import (
	"bytes"
	"errors"
	"exercises/Catalog/model"
	"exercises/Catalog/service"
	"github.com/golang/mock/gomock"
	"io/ioutil"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestCatHandler_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	catServ := service.NewMockServInterface(ctrl)
	testcasesById := []struct {
		reqMethod string
		expIn     int
		reqPath   string
		expOut    string
		expRet    model.Prod
		expCode   int
		expErr    error
	}{
		{"GET", 1, "/product?id=1", `{"id":1,"name":"Shoes","brand":{"name":"Puma"}}`, model.Prod{1, "Shoes", model.Brand{3, "Puma"}}, 200, nil},
		{"GET", 5, "/product?id=5", `{"id":0,"name":"","brand":{"name":""}}`, model.Prod{}, 404, errors.New("Id not found")},
		{"GET", 3, "/product?id=3", `{"id":3,"name":"Cricket Shoes","brand":{"name":""}}`, model.Prod{3, "Cricket Shoes", model.Brand{5, ""}}, 404, errors.New("Brand Id not found")},
		{"GET", 3, "/product?id=abc", "Invalid ID", model.Prod{3, "Cricket Shoes", model.Brand{5, ""}}, 400, errors.New("Brand Id not found")},

	}
	for i, tc := range testcasesById {
		catServ.EXPECT().GetById(tc.expIn).Return(tc.expRet, tc.expErr)
		catHandle := New(catServ)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(tc.reqMethod, tc.reqPath, nil)
		catHandle.Get(w, r)
		res := w.Result()
		resBody, resErr := ioutil.ReadAll(res.Body)
		resCode := w.Code
		//fmt.Println(string(resBody))
		if !reflect.DeepEqual(resCode, tc.expCode) {
			t.Errorf("Test %v has failed, Expected status code: %v but got %v", i, tc.expCode, resCode)
		}
		if !reflect.DeepEqual(string(resBody), tc.expOut) {
			t.Errorf("Test %v has failed, Expected: %v but got %v", i, tc.expOut, string(resBody))
		}
		if !reflect.DeepEqual(resErr, nil) {
			t.Errorf("Test %v has failed, Expected nil error but got %v", i, resErr)
		}
	}
	testcases := []struct {
		reqMethod string
		reqPath   string
		expOut    string
		expRet    []model.Prod
		expCode   int
		expErr    error
	}{
		{"GET","/product", `[{"id":1,"name":"Shoes","brand":{"name":"Puma"}}]`, []model.Prod{{1, "Shoes", model.Brand{3, "Puma"}}}, 200, nil},
		{"GET","/product", `[]`, []model.Prod{}, 404, errors.New("Products not found")},

	}
	for i, tc := range testcases {
		catServ.EXPECT().GetAll().Return(tc.expRet, tc.expErr)
		catHandle := New(catServ)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(tc.reqMethod, tc.reqPath, nil)
		catHandle.Get(w, r)
		res := w.Result()
		resBody, resErr := ioutil.ReadAll(res.Body)
		resCode := w.Code
		//fmt.Println(string(resBody))
		if !reflect.DeepEqual(resCode, tc.expCode) {
			t.Errorf("Test %v has failed, Expected status code: %v but got %v", i, tc.expCode, resCode)
		}
		if !reflect.DeepEqual(string(resBody), tc.expOut) {
			t.Errorf("Test %v has failed, Expected: %v but got %v", i, tc.expOut, string(resBody))
		}
		if !reflect.DeepEqual(resErr, nil) {
			t.Errorf("Test %v has failed, Expected nil error but got %v", i, resErr)
		}
	}
}
func TestCatHandler_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	catServ := service.NewMockServInterface(ctrl)
	testcases := []struct {
		reqMethod  string
		expInName  string
		expInBName string
		reqPath    string
		reqBody    []byte
		expOut     string
		expRet     model.Prod
		expCode    int
		expServErr error
	}{
		{"POST", "Cricket Shoes", "Puma", "/product", []byte(`{"name":"Cricket Shoes","brand":{"name":"Puma"}}`), `{"id":1,"name":"Cricket Shoes","brand":{"name":"Puma"}}`, model.Prod{1, "Cricket Shoes", model.Brand{3, "Puma"}}, 201, nil},
		{"POST", "Shoes", "Nike", "/product", []byte(`{"name":"Shoes","brand":{"name":"Nike"}}`), "Product not added to database", model.Prod{0, "", model.Brand{0, ""}}, 400, errors.New("Product not created")},
		//{"POST","Shoes","Nike","/product",[]byte(`{"names":"Shoes","brands":{"name":"Nike"}}`),"Product not added to database",model.Prod{0,"",model.Brand{0,""}},400,errors.New("Product not created")},

	}
	for i, tc := range testcases {
		catServ.EXPECT().Create(tc.expInName, tc.expInBName).Return(tc.expRet, tc.expServErr)
		catHandle := New(catServ)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(tc.reqMethod, tc.reqPath, bytes.NewBuffer(tc.reqBody))
		catHandle.Create(w, r)
		res := w.Result()
		resBody, resErr := ioutil.ReadAll(res.Body)
		resCode := w.Code
		//fmt.Println(string(resBody))
		if !reflect.DeepEqual(resCode, tc.expCode) {
			t.Errorf("Test %v has failed, Expected status code: %v but got %v", i, tc.expCode, resCode)
		}
		if !reflect.DeepEqual(string(resBody), tc.expOut) {
			t.Errorf("Test %v has failed, Expected: %v but got %v", i, tc.expOut, string(resBody))
		}
		if !reflect.DeepEqual(resErr, nil) {
			t.Errorf("Test %v has failed, Expected %v error but got %v", i, nil, resErr)
		}
	}
}

func TestCatHandler_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	catServ := service.NewMockServInterface(ctrl)
	testcases := []struct {
		reqMethod  string
		expInName  string
		expInBName string
		reqPath    string
		reqBody    []byte
		expOut     string
		expRet     model.Prod
		expCode    int
		expServErr error
	}{
		{"PUT", "Cricket Shoes", "Puma", "/product?id=1", []byte(`{"name":"Cricket Shoes","brand":{"name":"Puma"}}`), `{"id":1,"name":"Cricket Shoes","brand":{"name":"Puma"}}`, model.Prod{1, "Cricket Shoes", model.Brand{3, "Puma"}}, 200, nil},
		{"PUT", "Shoes", "Nike", "/product?id=0", []byte(`{"name":"Shoes","brand":{"name":"Nike"}}`), "Product not updated in database", model.Prod{0, "", model.Brand{0, ""}}, 400, errors.New("Product not updated")},
		{"PUT","Shoes","Nike","/product?id=abc",[]byte(`{"names":"Shoes","brands":{"name":"Nike"}}`),"Invalid ID",model.Prod{0,"",model.Brand{0,""}},400,errors.New("Product not created")},

	}
	for i, tc := range testcases {
		catServ.EXPECT().Update(tc.expRet.Id,tc.expInName, tc.expInBName).Return(tc.expRet, tc.expServErr)
		catHandle := New(catServ)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(tc.reqMethod, tc.reqPath, bytes.NewBuffer(tc.reqBody))
		catHandle.Update(w, r)
		res := w.Result()
		resBody, resErr := ioutil.ReadAll(res.Body)
		resCode := w.Code
		//fmt.Println(string(resBody))
		if !reflect.DeepEqual(resCode, tc.expCode) {
			t.Errorf("Test %v has failed, Expected status code: %v but got %v", i, tc.expCode, resCode)
		}
		if !reflect.DeepEqual(string(resBody), tc.expOut) {
			t.Errorf("Test %v has failed, Expected: %v but got %v", i, tc.expOut, string(resBody))
		}
		if !reflect.DeepEqual(resErr, nil) {
			t.Errorf("Test %v has failed, Expected %v error but got %v", i, nil, resErr)
		}
	}
}

func TestCatHandler_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	catServ := service.NewMockServInterface(ctrl)
	testcases := []struct {
		reqMethod  string
		reqPath    string
		expIn int
		expOut string
		expRetErr   error
		expCode    int
	}{
		{"DELETE","/product?id=1",1,"",nil,204},
		{"DELETE","/product?id=abc",1,"Invalid ID",nil,400},
		{"DELETE","/product?id=5",5, "Product not deleted in database", errors.New("Product not deleted from database"),400},
	}
	for i, tc := range testcases {
		catServ.EXPECT().Delete(tc.expIn).Return(tc.expRetErr)
		catHandle := New(catServ)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(tc.reqMethod, tc.reqPath, bytes.NewBuffer(nil))
		catHandle.Delete(w, r)
		res := w.Result()
		resBody, resErr := ioutil.ReadAll(res.Body)
		resCode := w.Code
		//fmt.Println(string(resBody))
		if !reflect.DeepEqual(resCode, tc.expCode) {
			t.Errorf("Test %v has failed, Expected status code: %v but got %v", i, tc.expCode, resCode)
		}
		if !reflect.DeepEqual(string(resBody), tc.expOut) {
			t.Errorf("Test %v has failed, Expected: %v but got %v", i, tc.expOut, string(resBody))
		}
		if !reflect.DeepEqual(resErr, nil) {
			t.Errorf("Test %v has failed, Expected %v error but got %v", i, nil, resErr)
		}
	}
}