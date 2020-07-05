package controllers

import (
	"errors"
	"fmt"
	"github.com/emicklei/go-restful"
	"go-with-mongodb/app/models"
	"go-with-mongodb/app/repository"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

type ProductsControllers interface {
	PostProduct(req *restful.Request, resp *restful.Response)
	GetProducts(req *restful.Request, resp *restful.Response)
	GetProduct(req *restful.Request, resp *restful.Response)
	PutProduct(req *restful.Request, resp *restful.Response)
	DeleteProduct(req *restful.Request, resp *restful.Response)
}

type productsControllers struct {
	productsRepository repository.ProductsRepository
}

func NewProductsControllers(productsRepository repository.ProductsRepository) ProductsControllers {
	return &productsControllers{productsRepository: productsRepository}
}

func (c *productsControllers) PostProduct(req *restful.Request, resp *restful.Response) {
	var product *models.Product
	if err := req.ReadEntity(&product); err != nil {
		_ = resp.WriteError(http.StatusUnprocessableEntity, err)
	} else {

		err := c.productsRepository.Create(product)
		if err != nil {
			_ = resp.WriteError(http.StatusUnprocessableEntity, err)
		} else {

			location := fmt.Sprintf("%s/%s",
				req.Request.RequestURI,
				product.ID.Hex())

			resp.AddHeader("Location", location)

			_ = resp.WriteHeaderAndEntity(http.StatusCreated, product.ID)
		}
	}
}

func (c *productsControllers) GetProducts(req *restful.Request, resp *restful.Response) {
	products, err := c.productsRepository.GetAll()
	if err != nil {
		_ = resp.WriteError(http.StatusBadRequest, err)
	} else {
		_ = resp.WriteAsJson(products)
	}
}

func (c *productsControllers) GetProduct(req *restful.Request, resp *restful.Response) {
	id := req.PathParameter("id")

	if !bson.IsObjectIdHex(id) {
		_ = resp.WriteError(http.StatusBadRequest, errors.New("invalid product_id"))
		return
	}

	product, err := c.productsRepository.GetByID(id)
	if err != nil {
		_ = resp.WriteError(http.StatusBadRequest, err)
	} else {
		_ = resp.WriteAsJson(product)
	}
}

func (c *productsControllers) PutProduct(req *restful.Request, resp *restful.Response) {
	var product *models.Product
	if err := req.ReadEntity(&product); err != nil {
		_ = resp.WriteError(http.StatusUnprocessableEntity, err)
	} else {
		product.ID = bson.ObjectIdHex(req.PathParameter("id"))
		err := c.productsRepository.Update(product)
		if err != nil {
			_ = resp.WriteError(http.StatusUnprocessableEntity, err)
		} else {
			_ = resp.WriteHeaderAndEntity(http.StatusOK, product)
		}
	}
}

func (c *productsControllers) DeleteProduct(req *restful.Request, resp *restful.Response) {
	id := req.PathParameter("id")

	if !bson.IsObjectIdHex(id) {
		_ = resp.WriteError(http.StatusBadRequest, errors.New("invalid product_id"))
		return
	}

	err := c.productsRepository.Delete(id)
	if err != nil {
		_ = resp.WriteError(http.StatusBadRequest, err)
	} else {
		resp.AddHeader("Entity", id)
		_ = resp.WriteHeaderAndEntity(http.StatusNoContent, nil)
	}
}
