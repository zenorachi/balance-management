package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/zenorachi/balance-management/internal/entity"
	"github.com/zenorachi/balance-management/internal/repository"
)

type ProductService struct {
	repo repository.Product
}

func NewProduct(repo repository.Product) *ProductService {
	return &ProductService{repo: repo}
}

func (p *ProductService) Create(ctx context.Context, product entity.Product) (int, error) {
	if p.isProductExists(ctx, product.Name) {
		return 0, entity.ErrProductAlreadyExists
	}

	if product.Price <= 0 {
		return 0, entity.ErrPriceIsNegative
	}

	return p.repo.Create(ctx, product)
}

func (p *ProductService) GetByID(ctx context.Context, id int) (entity.Product, error) {
	product, err := p.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Product{}, entity.ErrProductDoesNotExist
		}
		return entity.Product{}, err
	}

	return product, nil
}

func (p *ProductService) GetAll(ctx context.Context) ([]entity.Product, error) {
	products, err := p.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	if len(products) == 0 {
		return nil, entity.ErrProductDoesNotExist
	}

	return products, nil
}

func (p *ProductService) isProductExists(ctx context.Context, name string) bool {
	_, err := p.repo.GetByName(ctx, name)
	return !errors.Is(err, sql.ErrNoRows)
}
