package customers

import (
	"time"

	"github.com/uptrace/bun"
)

type Customer struct {
	bun.BaseModel `bun:"table:customers,alias:c"`

	ID        string    `bun:",pk,type:uuid" json:"id"`
	Name      string    `bun:",notnull" json:"firstname"`
	Email     string    `json:"email"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `bun:",soft_delete,nullzero"`
}
