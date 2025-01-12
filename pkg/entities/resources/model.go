package resources

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/papidb/aqua/pkg/config"
	"github.com/uptrace/bun"
	"golang.org/x/exp/rand"
)

type Resource struct {
	bun.BaseModel `bun:"table:resources,alias:r"`

	ID     string `bun:",pk,type:uuid" json:"id"`
	Name   string `bun:",notnull" json:"name"`
	Region string `json:"region"`

	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:",nullzero" json:"updated_at"`
	DeletedAt time.Time `bun:",soft_delete,nullzero" json:"-"`
}

// SeedResources seeds the database with sample resources
func SeedResources(app *config.App, maxResources int) error {
	db := app.Database.DB

	rand.Seed(uint64(time.Now().UnixNano()))
	ctx := context.Background()

	var resources []Resource
	uniqueNames := make(map[string]struct{})
	totalAdjectives := len(adjectives)
	totalResourceNames := len(resourceNames)

	for i := 0; i < maxResources; i++ {
		// Ensure we always generate enough resources
		baseName := resourceNames[i%totalResourceNames]
		adjective := adjectives[(i/totalResourceNames)%totalAdjectives]
		name := fmt.Sprintf("%s-%s", adjective, baseName)

		// Ensure uniqueness for the generated names
		if _, exists := uniqueNames[name]; exists {
			// Skip adding duplicate names and move to the next iteration
			continue
		}

		region := regions[rand.Intn(len(regions))]
		resources = append(resources, Resource{
			ID:     uuid.New().String(),
			Name:   name,
			Region: region,
		})
		uniqueNames[name] = struct{}{}
	}

	// Insert resources into the database in bulk
	_, err := db.NewInsert().Model(&resources).On("CONFLICT (name) DO NOTHING").Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to seed database: %w", err)
	}

	log.Printf("Database seeding completed. %d resources added.\n", len(resources))
	return nil
}
