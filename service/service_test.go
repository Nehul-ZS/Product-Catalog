package service

import (
	"errors"
	"exercises/Catalog/model"
	"exercises/Catalog/store/brand"
	"exercises/Catalog/store/product"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

type products struct {
	expIn   int
	pOut    model.Prod
	bOut    model.Brand
	expOut  model.Prod
	expPerr error
	expBerr error
	expErr  error
}

func TestCatService_GetById(t *testing.T) {
	ctrl := gomock.NewController(t)
	pS := product.NewMockStoreInterface(ctrl)
	bS := brand.NewMockStoreInterface(ctrl)
	testcases := []struct {
		expIn   int
		pOut    model.Prod
		bOut    model.Brand
		expOut  model.Prod
		expPerr error
		expBerr error
		expErr  error
	}{
		{1, model.Prod{1, "Shoes", model.Brand{3, ""}}, model.Brand{3, "Puma"}, model.Prod{
			1, "Shoes", model.Brand{3, "Puma"}}, nil, nil, error(nil)},
		{3, model.Prod{3, "Cricket Shoes", model.Brand{5, ""}}, model.Brand{}, model.Prod{
			3, "Cricket Shoes", model.Brand{5, ""}}, nil, errors.New("Brand Id not found"), errors.New("Brand Id not found")},
		{5, model.Prod{}, model.Brand{}, model.Prod{}, errors.New("Id not found"), errors.New("Id not found"), errors.New("Id not found")},
	}
	for i, tc := range testcases {
		pS.EXPECT().GetById(tc.expIn).Return(tc.pOut, tc.expPerr)
		if tc.expPerr == nil {
			bS.EXPECT().GetById(tc.pOut.BrandDetails.Id).Return(tc.bOut, tc.expBerr)
		}
		catServ := New(pS, bS)
		res, err := catServ.GetById(tc.expIn)
		if !reflect.DeepEqual(err, tc.expErr) {
			t.Errorf("For testcase %v, Expected %v error but got %v", i, tc.expErr, err)
		}

		if !reflect.DeepEqual(res, tc.expOut) {
			t.Errorf("For testcase %v,Expected %v but got %v", i, tc.expOut, res)
		}
	}
}
func TestCatService_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	pS := product.NewMockStoreInterface(ctrl)
	bS := brand.NewMockStoreInterface(ctrl)
	testcases := []struct {
		pOut    []model.Prod
		bOut    []model.Brand
		expOut  []model.Prod
		expPerr error
		expBerr error
		expErr  error
	}{
		{ []model.Prod{{1, "Shoes", model.Brand{3, ""}},{2, "Cricket Shoes", model.Brand{2, ""}}}, []model.Brand{{3, "Puma"},{2,"Nike"}}, []model.Prod{{
			1, "Shoes", model.Brand{3, "Puma"}},{2, "Cricket Shoes", model.Brand{2, "Nike"}}}, nil, nil, nil},
		{ []model.Prod{}, []model.Brand{}, []model.Prod{},  errors.New("Products not found"), errors.New("Brand Id not Found"), errors.New("Products not found")},
		{ []model.Prod{{1, "Shoes", model.Brand{3, ""}}}, []model.Brand{{}}, []model.Prod{{
			1, "Shoes", model.Brand{3, ""}}}, nil, errors.New("Brand Id not Found"),  errors.New("Brand Id not found")},

	}
	for i, tc := range testcases {
		pS.EXPECT().GetAll().Return(tc.pOut, tc.expPerr)
		if tc.expPerr == nil {
			for i,v:=range tc.pOut {
				bS.EXPECT().GetById(v.BrandDetails.Id).Return(tc.bOut[i], tc.expBerr)
			}
		}
		catServ := New(pS, bS)
		res, err := catServ.GetAll()
		if !reflect.DeepEqual(err, tc.expErr) {
			t.Errorf("For testcase %v, Expected %v error but got %v", i, tc.expErr, err)
		}

		if !reflect.DeepEqual(res, tc.expOut) {
			t.Errorf("For testcase %v,Expected %v but got %v", i, tc.expOut, res)
		}
	}
}
func TestCatService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	pS := product.NewMockStoreInterface(ctrl)
	bS := brand.NewMockStoreInterface(ctrl)
	testcases := []struct {
		expCheckBrand     model.Brand
		expCheckBrandErr  error
		expCreateBrand    model.Brand
		expBrandCreateErr error
		expProdCreate     model.Prod
		expProdCreateErr  error
		pOut              model.Prod
		bOut              model.Brand
		expOut            model.Prod
		expPerr           error
		expBerr           error
		expErr            error
	}{
		{expProdCreate: model.Prod{1, "Shoes", model.Brand{3, ""}}, expProdCreateErr: nil, expCheckBrand: model.Brand{3, "Puma"}, expCheckBrandErr: nil, expCreateBrand: model.Brand{3, "Puma"}, expBrandCreateErr: nil, pOut: model.Prod{1, "Shoes", model.Brand{3, ""}}, bOut: model.Brand{3, "Puma"}, expOut: model.Prod{
			1, "Shoes", model.Brand{3, "Puma"}}, expPerr: nil, expBerr: nil, expErr: nil},
		{expProdCreate: model.Prod{2, "Cricket Shoes", model.Brand{5, ""}}, expProdCreateErr: nil, expCheckBrand: model.Brand{0, "Asics"}, expCheckBrandErr: errors.New("Brand Name not Found"), expCreateBrand: model.Brand{5, "Asics"}, expBrandCreateErr: nil, pOut: model.Prod{2, "Cricket Shoes", model.Brand{5, "Asics"}}, bOut: model.Brand{}, expOut: model.Prod{
			2, "Cricket Shoes", model.Brand{5, "Asics"}}, expPerr: nil, expBerr: errors.New("Brand Id not found"), expErr: errors.New("Brand not found")},
		{expProdCreate: model.Prod{2, "Cricket Shoes", model.Brand{5, ""}}, expProdCreateErr: nil, expCheckBrand: model.Brand{0, "Asics"}, expCheckBrandErr: errors.New("Brand Name not Found"), expCreateBrand: model.Brand{5, "Asics"}, expBrandCreateErr: errors.New("Brand not created"), pOut: model.Prod{2, "Cricket Shoes", model.Brand{5, ""}}, bOut: model.Brand{}, expOut: model.Prod{}, expPerr: nil, expBerr: errors.New("Brand Id not found"),expErr: errors.New("Brand not created")},
		{expProdCreate: model.Prod{3, "Slippers", model.Brand{5, ""}}, expProdCreateErr: nil, expCheckBrand: model.Brand{0, "Crocs"}, expCheckBrandErr: errors.New("Brand Name not Found"), expCreateBrand: model.Brand{5, "Crocs"}, expBrandCreateErr: nil, pOut: model.Prod{}, bOut: model.Brand{}, expOut: model.Prod{}, expPerr: errors.New("Product not found"), expBerr: errors.New("Brand Id not found"), expErr: errors.New("Created Product not found")},
		{expProdCreate: model.Prod{0, "Sandals", model.Brand{3, ""}}, expProdCreateErr: errors.New("Employee not inserted"), expCheckBrand: model.Brand{3, "Mochi"}, expCheckBrandErr: nil, expCreateBrand: model.Brand{3, "Mochi"}, expBrandCreateErr: nil, pOut: model.Prod{}, bOut: model.Brand{}, expOut: model.Prod{}, expPerr: errors.New("Product not found"), expBerr: errors.New("Brand Id not found"), expErr: errors.New("Product not created")},
	}

	for i, tc := range testcases {
		bS.EXPECT().CheckBrand(tc.expCheckBrand.Brand).Return(tc.expCheckBrand.Id, tc.expCheckBrandErr)
		if tc.expCheckBrand.Id == 0 {
			bS.EXPECT().Create(tc.expCreateBrand.Brand).Return(tc.expCreateBrand.Id, tc.expBrandCreateErr)
		}
		pS.EXPECT().Create(tc.expProdCreate.Name, tc.expProdCreate.BrandDetails.Id).Return(tc.expProdCreate.Id, tc.expProdCreateErr)
		if tc.expProdCreateErr == nil {
			pS.EXPECT().GetById(tc.expProdCreate.Id).Return(tc.pOut, tc.expPerr)
			if tc.expPerr == nil {
				bS.EXPECT().GetById(tc.expProdCreate.BrandDetails.Id).Return(tc.bOut, tc.expBerr)
			}
		}
		catServ := New(pS, bS)
		res, err := catServ.Create(tc.expProdCreate.Name, tc.expCheckBrand.Brand)
		if !reflect.DeepEqual(err, tc.expErr) {
			t.Errorf("For testcase %v, Expected %v error but got %v", i, tc.expErr, err)
		}

		if !reflect.DeepEqual(res, tc.expOut) {
			t.Errorf("For testcase %v,Expected %v but got %v", i, tc.expOut, res)
		}
	}
}
func TestCatService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	pS := product.NewMockStoreInterface(ctrl)
	bS := brand.NewMockStoreInterface(ctrl)
	testcases := []struct {
		expCheckBrand     model.Brand
		expCheckBrandErr  error
		expCreateBrand    model.Brand
		expBrandCreateErr error
		expProdCreate     model.Prod
		expProdCreateErr  error
		pOut              model.Prod
		bOut              model.Brand
		expOut            model.Prod
		expPerr           error
		expBerr           error
		expErr            error
	}{
		{expProdCreate: model.Prod{1, "Shoes", model.Brand{3, ""}}, expProdCreateErr: nil, expCheckBrand: model.Brand{3, "Puma"}, expCheckBrandErr: nil, expCreateBrand: model.Brand{3, "Puma"}, expBrandCreateErr: nil, pOut: model.Prod{1, "Shoes", model.Brand{3, ""}}, bOut: model.Brand{3, "Puma"}, expOut: model.Prod{
			1, "Shoes", model.Brand{3, "Puma"}}, expPerr: nil, expBerr: nil, expErr: nil},
		{expProdCreate: model.Prod{2, "Cricket Shoes", model.Brand{5, ""}}, expProdCreateErr: nil, expCheckBrand: model.Brand{0, "Asics"}, expCheckBrandErr: errors.New("Brand Name not Found"), expCreateBrand: model.Brand{5, "Asics"}, expBrandCreateErr: nil, pOut: model.Prod{2, "Cricket Shoes", model.Brand{5, "Asics"}}, bOut: model.Brand{}, expOut: model.Prod{
			2, "Cricket Shoes", model.Brand{5, "Asics"}}, expPerr: nil, expBerr: errors.New("Brand Id not found"), expErr: errors.New("Updated Brand not found")},
		{expProdCreate: model.Prod{2, "Cricket Shoes", model.Brand{5, ""}}, expProdCreateErr: nil, expCheckBrand: model.Brand{0, "Asics"}, expCheckBrandErr: errors.New("Brand Name not Found"), expCreateBrand: model.Brand{5, "Asics"}, expBrandCreateErr: errors.New("Brand not created"), pOut: model.Prod{2, "Cricket Shoes", model.Brand{5, ""}}, bOut: model.Brand{}, expOut: model.Prod{}, expPerr: nil, expBerr: errors.New("Brand Id not found"),expErr: errors.New("Brand not created")},
		{expProdCreate: model.Prod{3, "Slippers", model.Brand{5, ""}}, expProdCreateErr: nil, expCheckBrand: model.Brand{0, "Crocs"}, expCheckBrandErr: errors.New("Brand Name not Found"), expCreateBrand: model.Brand{5, "Crocs"}, expBrandCreateErr: nil, pOut: model.Prod{}, bOut: model.Brand{}, expOut: model.Prod{}, expPerr: errors.New("Product not found"), expBerr: errors.New("Brand Id not found"), expErr: errors.New("Updated Product not found")},
		{expProdCreate: model.Prod{0, "Sandals", model.Brand{3, ""}}, expProdCreateErr: errors.New("Employee not inserted"), expCheckBrand: model.Brand{3, "Mochi"}, expCheckBrandErr: nil, expCreateBrand: model.Brand{3, "Mochi"}, expBrandCreateErr: nil, pOut: model.Prod{}, bOut: model.Brand{}, expOut: model.Prod{}, expPerr: errors.New("Product not found"), expBerr: errors.New("Brand Id not found"), expErr: errors.New("Product not updated")},
	}

	for i, tc := range testcases {
		bS.EXPECT().CheckBrand(tc.expCheckBrand.Brand).Return(tc.expCheckBrand.Id, tc.expCheckBrandErr)
		if tc.expCheckBrand.Id == 0 {
			bS.EXPECT().Create(tc.expCreateBrand.Brand).Return(tc.expCreateBrand.Id, tc.expBrandCreateErr)
		}
		pS.EXPECT().Update(tc.expProdCreate.Id,tc.expProdCreate.Name, tc.expProdCreate.BrandDetails.Id).Return(tc.expProdCreate.Id, tc.expProdCreateErr)
		if tc.expProdCreateErr == nil {
			pS.EXPECT().GetById(tc.expProdCreate.Id).Return(tc.pOut, tc.expPerr)
			if tc.expPerr == nil {
				bS.EXPECT().GetById(tc.expProdCreate.BrandDetails.Id).Return(tc.bOut, tc.expBerr)
			}
		}
		catServ := New(pS, bS)
		res, err := catServ.Update(tc.expProdCreate.Id,tc.expProdCreate.Name, tc.expCheckBrand.Brand)
		if !reflect.DeepEqual(err, tc.expErr) {
			t.Errorf("For testcase %v, Expected %v error but got %v", i, tc.expErr, err)
		}

		if !reflect.DeepEqual(res, tc.expOut) {
			t.Errorf("For testcase %v,Expected %v but got %v", i, tc.expOut, res)
		}
	}
}
func TestCatService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	pS := product.NewMockStoreInterface(ctrl)
	bS := brand.NewMockStoreInterface(ctrl)
	testcases := []struct {
		expIn int
		expErr            error
		expRet int
		expRetErr error
	}{
		{1,nil,1,nil},
		{10,errors.New("Product not deleted from database"),0,errors.New("Product not deleted from database")},
		{5,errors.New("Product not deleted from database"),0,nil},
	}

	for i, tc := range testcases {
		pS.EXPECT().Delete(tc.expIn).Return(tc.expRet,tc.expRetErr)
		catServ := New(pS, bS)
		resErr := catServ.Delete(tc.expIn)
		if !reflect.DeepEqual(resErr, tc.expErr) {
			t.Errorf("For testcase %v, Expected %v error but got %v", i, tc.expErr, resErr)
		}
		}
}
