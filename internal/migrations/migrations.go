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
				if err := tx.AutoMigrate(&models.User{}); err != nil {
					return err
				}

				if err := tx.AutoMigrate(&models.UserToken{}); err != nil {
					return err
				}

				if err := tx.AutoMigrate(&models.Wallet{}); err != nil {
					return err
				}

				if err := tx.AutoMigrate(&models.Transaction{}); err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				if err := tx.Migrator().DropTable("transactions"); err != nil {
					return err
				}
				if err := tx.Migrator().DropTable("wallets"); err != nil {
					return err
				}
				if err := tx.Migrator().DropTable("user_tokens"); err != nil {
					return err
				}
				if err := tx.Migrator().DropTable("users"); err != nil {
					return err
				}
				return nil
			},
		},
	})
}
