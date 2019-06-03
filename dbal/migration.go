package dbal

import (
	"net/http"
	"sync"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/source"
	vfsdata "github.com/neermitt/migrate-vfsdata-source"
	"github.com/sirupsen/logrus"
)

type MigrationDriverFactory func(db DBDriver, migrationTable string) (database.Driver, error)

var (
	migrationDriverFactories = make(map[string]MigrationDriverFactory, 0)
	mdfmtx                   sync.Mutex
)

type DBDriver interface {
	DriverName() string
}

func MigrationSourceDriver(fs http.FileSystem, path string) source.Driver {
	driver, err := vfsdata.WithInstance(vfsdata.Resource(path, fs))
	if err != nil {
		logrus.WithError(err).Fatalf("Failed to create source driver for migrations")
		panic(err)
	}
	return driver
}

func ToMigrate(sourceDriver source.Driver, db DBDriver, migrationTable string) (*migrate.Migrate, error) {
	var dbDriver database.Driver

	if dbDriverFactory, err := GetMigrationDriverFactoryFor(db.DriverName()); err != nil {
		return nil, err
	} else {
		dbDriver, _ = dbDriverFactory(db, migrationTable)
	}

	return migrate.NewWithInstance("source", sourceDriver, db.DriverName(), dbDriver)
}

func MigrateUp(sourceDriver source.Driver, db DBDriver, migrationTable string) (int, error) {
	mig, err := ToMigrate(sourceDriver, db, migrationTable)
	if err != nil {
		return 0, err
	}

	err = mig.Up()
	version, _, _ := mig.Version()
	if err != nil {
		return int(version), err
	}
	return int(version), nil

}

func MigrateDown(sourceDriver source.Driver, db DBDriver, migrationTable string) error {
	mig, err := ToMigrate(sourceDriver, db, migrationTable)
	if err != nil {
		return err
	}

	return mig.Down()
}

// RegisterDriver registers a driver
func RegisterMigrationDriverFactory(driverName string, d MigrationDriverFactory) {
	mdfmtx.Lock()
	migrationDriverFactories[driverName] = d
	mdfmtx.Unlock()
}

// GetDriverFor returns a driver for the given DSN or ErrNoResponsibleDriverFound if no driver was found.
func GetMigrationDriverFactoryFor(driverName string) (MigrationDriverFactory, error) {
	if factory, hasDriver := migrationDriverFactories[driverName]; hasDriver {
		return factory, nil
	}
	return nil, ErrNoResponsibleDriverFound
}
