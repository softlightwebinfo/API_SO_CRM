package ORM

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"so-crm/src/libs/ORM/model"
)

type ORM struct {
	db  *sql.DB
	log bool
}

func (o *ORM) LogMode(isLog bool) {
	o.log = isLog
}

func (o *ORM) Query() *model.ORMQuery {
	return &model.ORMQuery{DB: o.db, IsLog: o.log}
}

func Open(host string, user string, pass string, dbName string, port string) *ORM {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		user,
		pass,
		dbName,
	)
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	//defer db.Close()

	err = db.Ping()

	if err != nil {
		panic(err)
	}
	return &ORM{
		db: db,
	}
}
