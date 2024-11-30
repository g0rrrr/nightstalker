package models

import (
	"database/sql"
	"fmt"
	"github.com/coopernurse/gorp"
)

var db_map *gorp.DbMap

func GetDbSession() *gorp.DbMap {
	if db_map != nil {
		return db_map
	}

    db_username := "nightstalker"
    db_password := "gor"
    db_database := "nightstalker"
    db_hostname := "localhost"
    db_port     := "5432"

	db, err := sql.Open("postgres",
	"user="+db_username+
			" password="+db_password+
			" dbname="+db_database+
			" host="+db_hostname+
			" port="+db_port+
			" sslmode=disable")

	if err != nil {
        fmt.Printf("[NIGHSTALKER ERROR] CANNOT OPEN THE DATABASE: %s\n", err.Error())
		return nil
	}

	db_map = &gorp.DbMap{
		Db:      db,
		Dialect: gorp.PostgresDialect{},
	}

	// TODO: Do we need this every time?
	db_map.AddTableWithName(User{}, "users").SetKeys(true, "Id")
	db_map.AddTableWithName(Board{}, "boards").SetKeys(true, "Id")
	db_map.AddTableWithName(Post{}, "posts").SetKeys(true, "Id")
	db_map.AddTableWithName(View{}, "views").SetKeys(false, "Id")
	db_map.AddTableWithName(Setting{}, "settings").SetKeys(true, "Key")

	return db_map
}
