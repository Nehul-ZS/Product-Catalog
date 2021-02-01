package model

type Prod struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	BrandDetails Brand  `json:"brand"`
}
