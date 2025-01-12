package customers

import (
	"context"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/uptrace/bun"
)

type Repo struct {
	db bun.IDB
}

func NewRepo(db bun.IDB) *Repo {
	return &Repo{db: db}
}

// WithDB returns a new instance of the repository but with the db set to the provided db connection.
// This helps with separation of concerns even with a shared transaction
func (r *Repo) WithDB(db bun.IDB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Create(ctx context.Context, customer *Customer) error {
	_, err := r.db.NewInsert().
		Model(customer).
		ExcludeColumn("created_at").
		Returning("*").
		Exec(ctx)

	if err, ok := err.(*pgconn.PgError); ok && err.Code == pgerrcode.UniqueViolation {
		return ErrExistingEmailOrName{}
	}
	return err

}

func (r *Repo) CreateCustomerResource(ctx context.Context, customerResource *CustomerResource) error {
	_, err := r.db.NewInsert().
		Model(customerResource).
		ExcludeColumn("created_at").
		Returning("*").
		Exec(ctx)

	if err, ok := err.(*pgconn.PgError); ok && err.Code == pgerrcode.UniqueViolation {
		return ErrExistingCustomerResource{}
	}
	return err

}

func (r *Repo) Find(ctx context.Context, id string) (*Customer, error) {
	customer := &Customer{}
	err := r.db.NewSelect().
		Model(customer).
		Where("id = ?", id).
		Scan(ctx)

	return customer, err
}

func (r *Repo) DeleteCustomerResource(ctx context.Context, resource_id string) error {
	_, err := r.db.NewDelete().
		Model((*CustomerResource)(nil)).
		Where("resource_id = ?", resource_id).
		Exec(ctx)
	return err
}
