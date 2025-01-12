package customers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/papidb/aqua/pkg/internal"
	"github.com/uptrace/bun"
)

type CustomerService struct {
	db            *bun.DB
	CustomersRepo *Repo
}

func NewService(db *bun.DB, customersRepo *Repo) *CustomerService {
	return &CustomerService{db, customersRepo}
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
