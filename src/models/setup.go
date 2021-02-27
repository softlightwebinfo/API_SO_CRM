package models

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"so-crm/src/libs/ORM"
	"strconv"
)

var DB *ORM.ORM

func ConnectDataBase(configFile Config) {
	db := ORM.Open(
		configFile.Database.Host,
		configFile.Database.User,
		configFile.Database.Pass,
		configFile.Database.DB,
		strconv.Itoa(configFile.Database.Port),
	)

	db.LogMode(true)

	DB = db
}
