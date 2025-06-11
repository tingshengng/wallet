package migrations

import (
	"wallet/internal/models"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func NewMigrator(db *gorm.DB) *gormigrate.Gormigrate {
	return gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "20250611211100",
			Migrate: func(tx *gorm.DB) error {
				// Create tests table
				if err := tx.AutoMigrate(&models.Test{}); err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				if err := tx.Migrator().DropTable("tests"); err != nil {
					return err
				}
				return nil
			},
		},
	})
}
