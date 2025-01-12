package resources

import (
	"context"

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

func (r *Repo) Find(ctx context.Context, id string) (*Resource, error) {
	resource := &Resource{}
	err := r.db.NewSelect().
		Model(resource).
		Where("id = ?", id).
		Scan(ctx)

	return resource, err
}
