package customers

import (
	"time"

	"github.com/uptrace/bun"
)

type Customer struct {
	bun.BaseModel `bun:"table:customers,alias:c"`

	ID    string `bun:",pk,type:uuid" json:"id"`
	Name  string `bun:",notnull" json:"firstname"`
	Email string `json:"email"`

	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `bun:",soft_delete,nullzero" json:"-"`
}

type CustomerResource struct {
	bun.BaseModel `bun:"table:customer_resources,alias:cr"`
	CustomerID    string    `bun:",pk,type:uuid" json:"customer_id"`
	ResourceID    string    `bun:",pk,type:uuid" json:"resource_id"`
	CreatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`
}
