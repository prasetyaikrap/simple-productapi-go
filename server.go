package main

import (
	"simple-productapi-golang/service"

	"github.com/labstack/echo/v4"
)

func main() {
	//init echo instances
	e := echo.New()

	// set Routing
	// Home API
	e.GET("/", service.HomeAPI)
	//Get all products
	e.GET("/products", service.GetAllProducts)
	//Add new product to products collection
	e.POST("/products", service.CreateProduct)
	//Get product by ID
	e.GET("/products/:id", service.GetProductById)
	//Update product by ID
	e.PUT("/products/:id", service.UpdateProduct)
	//Update product by ID on selected key
	e.PATCH("/products/:id", service.PatchProduct)
	//Delete product by ID
	e.DELETE("/products/:id", service.DeleteProduct)

	//Start Server
	e.Logger.Fatal(e.Start(":1323"))
}