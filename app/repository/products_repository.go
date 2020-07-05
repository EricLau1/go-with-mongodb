package repository

import (
	"errors"
	"go-with-mongodb/app/models"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type ProductsRepository interface {
	Create(product *models.Product) error
	GetAll() (products []*models.Product, err error)
	GetByID(id string) (product *models.Product, err error)
	Update(product *models.Product) error
	Delete(id string) error
}

type productsRepository struct {
	db *mgo.Database
}

func NewProductsRepository(db *mgo.Database) ProductsRepository {
	return &productsRepository{db}
}

func (p *productsRepository) Create(product *models.Product) error {
	product.ID = bson.NewObjectId()
	product.CreatedAt = time.Now()
	product.UpdatedAt = product.CreatedAt
	return p.db.C(PRODUCTS).Insert(product)
}

func (p *productsRepository) GetAll() (products []*models.Product, err error) {
	err = p.db.C(PRODUCTS).Find(bson.M{}).All(&products)
	return
}

func (p *productsRepository) GetByID(id string) (product *models.Product, err error) {
	err = p.db.C(PRODUCTS).FindId(bson.ObjectIdHex(id)).One(&product)
	if product == nil {
		return nil, errors.New("product not found")
	}
	return
}

func (p *productsRepository) Update(product *models.Product) error {
	found, err := p.GetByID(product.ID.Hex())
	if err != nil {
		return err
	}
	product.CreatedAt = found.CreatedAt
	product.UpdatedAt = time.Now()
	return p.db.C(PRODUCTS).UpdateId(product.ID, product)
}

func (p *productsRepository) Delete(id string) error {
	return p.db.C(PRODUCTS).RemoveId(bson.ObjectIdHex(id))
}
