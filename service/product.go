package service

import (
	"net/http"
	"simple-productapi/models"
	"simple-productapi/repository"

	"strconv"

	"github.com/labstack/echo/v4"
)

type (
	IProduct interface{
		findById(id int) models.Product
		IndexOf(id int) int
	}

	ProductRepo struct {
		db []models.Product
		pid int
	}
)

var (
	repo = NewRepository()
)

func NewRepository() *ProductRepo{
	return &ProductRepo{
		db: repository.ProductsDB,
		pid: repository.ProductId,
	}
}

func HomeAPI(c echo.Context) error {
		return c.JSON(http.StatusOK, models.SuccessResponse{
			Code: http.StatusOK,
			Message: "API is Active",
			Data: "Home API",
		})
}

func CreateProduct(c echo.Context) error {
	dataQuery := c.QueryParam("data")
	if dataQuery == "bulk" {
		p := &[]models.Product{}
		if err := c.Bind(p); err != nil {
			return c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Code: http.StatusBadRequest,
				Message: "Invalid request body",
				Error: err.Error(),
			})
		}
		productRes := []models.ProductResponse{}
		for _, product := range *p {
			product.Id = repo.pid
			repo.db = append(repo.db, product)
			productRes = append(productRes, models.ProductResponse{
				Id: product.Id,
				Name: product.Name,
				Quantity: product.Quantity,
				Price: product.Price,
			})
			repo.pid++
		}
		return c.JSON(http.StatusCreated, models.SuccessResponse{
			Code: http.StatusCreated,
			Message: "Bulk products created successfully",
			Data: productRes,
		})
	} else {
		p := &models.Product{
			Id: repo.pid,
		}
		if err := c.Bind(p); err != nil {
			return c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Code: http.StatusBadRequest,
				Message: "Invalid request body",
				Error: err.Error(),
			})
		}
		repo.db = append(repo.db, *p)
		repo.pid++
		productRes := models.ProductResponse{
			Id: p.Id,
			Name: p.Name,
			Quantity: p.Quantity,
			Price: p.Price,
		}
		return c.JSON(http.StatusCreated, models.SuccessResponse{
			Code: http.StatusCreated,
			Message: "Product created successfully",
			Data: productRes,
		})
	}
}

func GetAllProducts(c echo.Context) error {
	if len(repo.db) <= 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Code: http.StatusNotFound,
				Message: "Get products failed.",
				Error: "Products is empty",
			})
	}
	productRes := []models.ProductResponse{}
	for _, product := range repo.db {
		appendData := models.ProductResponse{
			Id: product.Id,
			Name: product.Name,
			Quantity: product.Quantity,
			Price: product.Price,
		}
		productRes = append(productRes, appendData)
	}
	return c.JSON(http.StatusOK, models.SuccessResponse{
		Code: http.StatusOK,
		Message: "Products found",
		Data: productRes,
	})
}

func GetProductById(c echo.Context) error {
	intId, _ := strconv.Atoi(c.Param("id"))
	product := repo.FindById(intId)
	if product.Id != 0 {
		productRes := models.ProductResponse{
			Id: product.Id,
			Name: product.Name,
			Quantity: product.Quantity,
			Price: product.Price,
		}
		return c.JSON(http.StatusOK, models.SuccessResponse{
			Code: http.StatusOK,
			Message: "Product found",
			Data: productRes,
		})
	} else {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{
			Code: http.StatusNotFound,
			Message: "Get product failed",
			Error: "Product not found or does not exist",
		})
	}
}

func UpdateProduct(c echo.Context) error {
	intid,_ := strconv.Atoi(c.Param("id"))
	pd := &models.Product{}
	if err := c.Bind(pd); err != nil {
		return err
	}
	pIndex := repo.IndexOf(intid)
	if pIndex != -1 {
		pd.Id = repo.db[pIndex].Id
		repo.db[pIndex] = *pd
		productRes := models.ProductResponse{
			Id: repo.db[pIndex].Id,
			Name: repo.db[pIndex].Name,
			Quantity: repo.db[pIndex].Quantity,
			Price: repo.db[pIndex].Price,
		}
		return c.JSON(http.StatusOK, models.SuccessResponse{
			Code: http.StatusOK,
			Message: "Product updated successfully",
			Data: productRes,
		})
	}
	return c.JSON(http.StatusNotFound, models.ErrorResponse{
		Code: http.StatusNotFound,
		Message: "Update product failed",
		Error: "Product not found",
	})
}

func PatchProduct(c echo.Context) error {
	method := c.QueryParam("method")
	if method == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code: http.StatusBadRequest,
			Message: "Invalid request query",
			Error: "Must have 'method' query",
		})
	}
	p := &models.Product{}
	if err := c.Bind(p); err != nil {
			return err
		}
	intid,_ := strconv.Atoi(c.Param("id"))
	pIndex := repo.IndexOf(intid)
	if pIndex != -1 {
		switch method {
		case "add-stock":
			if p.Quantity == 0 {
        return c.JSON(http.StatusNotFound, models.ErrorResponse{
          Code: http.StatusNotFound,
          Message: "Patch update product failed",
          Error: "Invalid payload request",
        })
      }
			repo.db[pIndex].Quantity = repo.db[pIndex].Quantity + p.Quantity
			productRes := models.ProductResponse{
				Id: repo.db[pIndex].Id,
				Name: repo.db[pIndex].Name,
				Quantity: repo.db[pIndex].Quantity,
				Price: repo.db[pIndex].Price,
			}
			return c.JSON(http.StatusOK, models.SuccessResponse{
				Code: http.StatusOK,
				Message: "Product stock added successfully",
				Data: productRes,
			})
		case "reduce-stock":
			if p.Quantity == 0 {
        return c.JSON(http.StatusNotFound, models.ErrorResponse{
          Code: http.StatusNotFound,
          Message: "Patch update product failed",
          Error: "Invalid payload request",
        })
      }
			repo.db[pIndex].Quantity = repo.db[pIndex].Quantity - p.Quantity
			productRes := models.ProductResponse{
				Id: repo.db[pIndex].Id,
				Name: repo.db[pIndex].Name,
				Quantity: repo.db[pIndex].Quantity,
				Price: repo.db[pIndex].Price,
			}
			return c.JSON(http.StatusOK, models.SuccessResponse{
				Code: http.StatusOK,
				Message: "Product stock reduced successfully",
				Data: productRes,
			})
		case "update-price":
			if p.Price == 0 {
        return c.JSON(http.StatusNotFound, models.ErrorResponse{
          Code: http.StatusNotFound,
          Message: "Patch update product failed",
          Error: "Invalid payload request",
        })
      }
			repo.db[pIndex].Price = p.Price
			productRes := models.ProductResponse{
				Id: repo.db[pIndex].Id,
				Name: repo.db[pIndex].Name,
				Quantity: repo.db[pIndex].Quantity,
				Price: repo.db[pIndex].Price,
			}
			return c.JSON(http.StatusOK, models.SuccessResponse{
				Code: http.StatusOK,
				Message: "Product price updated successfully",
				Data: productRes,
			})
		}
	}
	return c.JSON(http.StatusNotFound, models.ErrorResponse{
		Code: http.StatusNotFound,
		Message: "Patch update product failed",
		Error: "Product not found",
	})
}
func DeleteProduct(c echo.Context) error {
	intid, _ := strconv.Atoi(c.Param("id"))
	pIndex := repo.IndexOf(intid)
	if pIndex != -1 {
		copy(repo.db[pIndex:], repo.db[pIndex+1:])
		repo.db = repo.db[:len(repo.db)-1]
		return c.JSON(http.StatusOK, models.SuccessResponse{
			Code: http.StatusOK,
			Message: "Product Deleted successfully",
			Data: "Deleted product ID:" + c.Param("id"),
		})
	} 
	return c.JSON(http.StatusNotFound, models.ErrorResponse{
		Code: http.StatusNotFound,
		Message: "Delete product failed",
		Error: "Product not found or does not exist",
	})
}

func (repo ProductRepo) FindById(id int) models.Product {
	for _, product := range repo.db {
		if product.Id == id {
			return product
		}
	}
	return models.Product{}
}

func (repo ProductRepo) IndexOf(id int) int {
	for k, v := range repo.db {
		if v.Id == id {
			return k
		}
	}
	return -1
}