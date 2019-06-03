package postgres

import (
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/open-identity/utils/dbal"
)

const (
	DriverPostgresSQL = "postgres"
)

func init() {
	dbal.RegisterMigrationDriverFactory(DriverPostgresSQL, func(db dbal.DBDriver, migrationTable string) (database.Driver, error) {
		return postgres.WithInstance(db.(*sqlx.DB).DB, &postgres.Config{
			MigrationsTable: migrationTable,
		})
	})
}
