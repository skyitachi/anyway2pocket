package common

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// TODO: use toml storage config
const (
	DB_NAME = "pocket"
	DB_USER = "skyitachi"
	// TODO: secret information
	DB_HOST = "localhost"
)

var dbClient *sql.DB

func createTable() {
	tableStr := `
		create table if not exists urlinfo (
			id serial not null,
			url character varying(100) not null,
			status smallint not null,
			created date,
			finished date,
			constraint urlinfo_pkey primary key (id)
		)
	`
	_, err := dbClient.Exec(tableStr)
	if err != nil {
		log.Fatal("[Database] create error: " + err.Error())
	}
}

func init() {
	var err error
	dbInfo := fmt.Sprintf("user=%s dbname=%s host=%s sslmode=disable", DB_USER, DB_NAME, DB_HOST)
	dbClient, err = sql.Open("postgres", dbInfo)
	createTable()
	if err != nil {
		log.Fatal("[Database] open error: " + err.Error())
	}
}

// URLExists check the url exists
func URLExists(url string) bool {
	var ret int
	query := fmt.Sprintf("select id from urlinfo where url='%s'", url)
	err := dbClient.QueryRow(query).Scan(&ret)
	switch {
	case err == sql.ErrNoRows:
		return false
	case err != nil:
		log.Fatalf("[URLExists]: query %s error: %s", url, err.Error())
	}
	return true
}
