package db

import (
	"fmt"
	"forecast/config"
	"github.com/go-pg/migrations/v8"
	"github.com/go-pg/pg/v10"
	"log"
)

func StartDB(config *config.Configuration) (*pg.DB, error) {
	var (
		opts *pg.Options
		err  error
	)

	addr := fmt.Sprintf("%s:%d", config.Db.Host, config.Db.Port)
	opts = &pg.Options{
		//default port
		//depends on the db service from docker compose
		Addr:     addr,
		User:     config.Db.User,
		Password: config.Db.Password,
		Database: config.Db.Database,
	}

	//connect db
	db := pg.Connect(opts)
	//run migrations
	collection := migrations.NewCollection()
	collection.SetTableName("forecast_migrations")
	err = collection.DiscoverSQLMigrations("migrations")
	if err != nil {
		return nil, err
	}

	//start the migrations
	_, _, err = collection.Run(db, "init")
	if err != nil {
		return nil, err
	}

	oldVersion, newVersion, err := collection.Run(db, "up")
	if err != nil {
		return nil, err
	}
	if newVersion != oldVersion {
		log.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
	} else {
		log.Printf("version is %d\n", oldVersion)
	}

	//return the db connection
	return db, err
}
