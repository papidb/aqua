package customers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/papidb/aqua/pkg/api"
	"github.com/papidb/aqua/pkg/entities/resources"
	"github.com/papidb/aqua/pkg/internal"
	"github.com/uptrace/bun"
)

type CustomerService struct {
	db            *bun.DB
	CustomersRepo *Repo
	ResourcesRepo *resources.Repo
}

func NewService(db *bun.DB, customersRepo *Repo, resourcesRepo *resources.Repo) *CustomerService {
	return &CustomerService{db, customersRepo, resourcesRepo}
}

func (svc *CustomerService) CreateCustomer(ctx context.Context, dto CreateCustomerDTO) (*Customer, error) {
	tx, err := svc.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	customer := Customer{
		ID:        internal.GenerateUUID(),
		Name:      dto.Name,
		Email:     dto.Email,
		CreatedAt: time.Now(),
	}

	if err := svc.CustomersRepo.WithDB(tx).Create(ctx, &customer); err != nil {
		tx.Rollback()
		if ok := errors.As(err, &ErrExistingEmailOrName{}); ok {
			return nil, ErrExistingEmailOrName{}
		}
		return nil, fmt.Errorf("unable to create customer: %s", err.Error())
	}

	return &customer, tx.Commit()
}

func (svc *CustomerService) AddResourceToCustomer(ctx context.Context, customer_id string, dto AddResourceToCustomerDTO) (*Customer, *resources.Resource, error) {
	tx, err := svc.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, nil, err
	}

	// search for resource
	resource, err := svc.ResourcesRepo.WithDB(tx).Find(ctx, dto.ResourceID)

	if (err != nil) || (resource == nil) {
		tx.Rollback()
		return nil, nil, api.ErrResourceNotFound
	}

	customer, err := svc.CustomersRepo.WithDB(tx).Find(ctx, customer_id)

	if (err != nil) || (customer == nil) {
		tx.Rollback()
		return nil, nil, api.ErrCustomerNotFound
	}

	// link customer to resource

	customerResource := &CustomerResource{
		CustomerID: customer.ID,
		ResourceID: resource.ID,
	}

	if err := svc.CustomersRepo.WithDB(tx).CreateCustomerResource(ctx, customerResource); err != nil {
		tx.Rollback()
		if ok := errors.As(err, &ErrExistingCustomerResource{}); ok {
			return nil, nil, ErrExistingCustomerResource{}
		}
		return nil, nil, fmt.Errorf("unable to create customer: %s", err.Error())
	}

	return customer, resource, tx.Commit()
}
