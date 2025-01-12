package resources

import (
	"context"
	"database/sql"

	"github.com/papidb/aqua/pkg/api"
	"github.com/uptrace/bun"
)

type ResourceService struct {
	db            *bun.DB
	ResourcesRepo *Repo
}

func NewService(db *bun.DB, resourcesRepo *Repo) *ResourceService {
	return &ResourceService{db, resourcesRepo}
}

func (svc *ResourceService) UpdateResource(ctx context.Context, resource_id string, dto UpdateResourceDTO) (*Resource, error) {
	// find resource
	tx, err := svc.db.BeginTx(ctx, &sql.TxOptions{})

	if err != nil {
		return nil, err
	}

	resource, err := svc.ResourcesRepo.WithDB(svc.db).Find(ctx, resource_id)

	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return nil, api.ErrResourceNotFound
		}
		return nil, err
	}

	resource.Name = dto.Name
	resource.Region = dto.Region

	err = svc.ResourcesRepo.WithDB(svc.db).Update(ctx, resource)

	if err != nil {
		return nil, err
	}

	return resource, tx.Commit()
}

func (svc *ResourceService) DeleteResource(ctx context.Context, resource_id string) (*Resource, error) {
	// find resource
	tx, err := svc.db.BeginTx(ctx, &sql.TxOptions{})

	if err != nil {
		return nil, err
	}

	resource, err := svc.ResourcesRepo.WithDB(tx).Find(ctx, resource_id)

	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return nil, api.ErrResourceNotFound
		}
		return nil, err
	}

	err = svc.ResourcesRepo.WithDB(tx).Delete(ctx, resource)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return resource, tx.Commit()
}
