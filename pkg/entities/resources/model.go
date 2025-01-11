package resources

import (
	"time"

	"github.com/uptrace/bun"
)

type Resource struct {
	bun.BaseModel `bun:"table:resources,alias:r"`

	ID     string `bun:",pk,type:uuid" json:"id"`
	Name   string `bun:",notnull" json:"firstname"`
	Region string `json:"region"`

	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `bun:",soft_delete,nullzero"`
}
