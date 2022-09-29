package models

type (
	Product struct{
		Id int `json:"id"`
		Name string `json:"name"`
		Quantity int `json:"quantity"`
		Price int `json:"price"`
	}

	ProductResponse struct{
		Id int `json:"id"`
		Name string `json:"name"`
		Quantity int `json:"quantity"`
		Price int `json:"price"`
	}
)