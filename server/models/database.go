package models

import (
	"database/sql"
	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
	"log"
)

type MyDb struct {
	gorp.DbMap
	sqlDb *sql.DB
}

func NewDb() *MyDb {
	db, err := sql.Open("postgres", "user=postgres password=postgres dbname=mila_dev sslmode=disable")
	checkErr(err, "sql.Open failed")
	dbmap := &MyDb{gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}, db}
	dbmap.AddTableWithName(User{}, "users").SetKeys(true, "Id")
	dbmap.AddTableWithName(Connection{}, "connections").SetKeys(false, "user1_id", "user2_id")
	dbmap.AddTableWithName(Post{}, "posts").SetKeys(true, "Id")
	dbmap.AddTableWithName(Comment{}, "comments").SetKeys(true, "Id")
	dbmap.AddTableWithName(Star{}, "stars").SetKeys(false, "post_id", "user_id")
	dbmap.AddTableWithName(Invite{}, "invites").SetKeys(false, "user1_id", "user2_id")
	dbmap.AddTableWithName(Activity{}, "activities").SetKeys(true, "Id")
	dbmap.AddTableWithName(Kid{}, "kids").SetKeys(true, "Id")
	dbmap.AddTableWithName(Feed{}, "feeds").SetKeys(false, "user_id", "post_id")
	return dbmap
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
