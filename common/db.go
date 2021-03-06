package common

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// PocketDBClient db client
type PocketDBClient struct {
	DBName   string
	DBUser   string
	DBHost   string
	dbClient *sql.DB
}

// Init init the database
func (client *PocketDBClient) Init() {
	logger := GetLogger()
	var err error
	dbInfo := fmt.Sprintf("user=%s dbname=%s host=%s sslmode=disable", client.DBUser, client.DBName, client.DBHost)
	client.dbClient, err = sql.Open("postgres", dbInfo)
	client.createTable()
	if err != nil {
		logger.Fatal("[Database] open error: " + err.Error())
	}
	logger.Info("[PocketDBClient]: Init successfully")
}

func (client *PocketDBClient) createTable() {
	logger := GetLogger()
	tableStr := `
		create table if not exists urlinfo (
			id serial not null,
			url character varying(100) not null,
			status smallint not null,
			created timestamp with time zone default NOW(),
			updated timestamp with time zone default NOW(),
			finished timestamp with time zone,
			constraint urlinfo_pkey primary key (id)
		)
	`
	_, err := client.dbClient.Exec(tableStr)
	if err != nil {
		logger.Fatal("[Database] create error: " + err.Error())
	}
}

// URLExists check the url exists
func (client *PocketDBClient) URLExists(url string) bool {
	logger := GetLogger()
	var ret int
	query := fmt.Sprintf("select id from urlinfo where url='%s'", url)
	err := client.dbClient.QueryRow(query).Scan(&ret)
	switch {
	case err == sql.ErrNoRows:
		return false
	case err != nil:
		logger.Fatalf("[URLExists]: query %s error: %s", url, err.Error())
	}
	return true
}

// AddURL add url to the datebase
func (client *PocketDBClient) AddURL(url string, status int) error {
	logger := GetLogger()
	query := fmt.Sprintf("insert into urlinfo (url, status) values('%s', %d)", url, status)
	_, err := client.dbClient.Exec(query)
	if err != nil {
		logger.Error("[PocketDBClient]: AddURL error " + err.Error())
		return err
	}
	return nil
}

// UpdateURL if not exists will use AddURL
func (client *PocketDBClient) UpdateURL(url string) {
	logger := GetLogger()
	cTimestamp := GetCurrentTimestamp()
	query := fmt.Sprintf("update urlinfo set updated = to_timestamp(%d) where url = '%s'", cTimestamp, url)
	result, err := client.dbClient.Exec(query)
	if err != nil {
		logger.Infof("[UpdateURL]: update %s error: %s", url, err.Error())
	}
	c, _ := result.RowsAffected()
	if c == 0 {
		logger.Info("[UpdateURL]: in the update url: " + url)
		client.AddURL(url, URLStatusCreated)
	}
}

// GetDateByURL get url updated date
func (client *PocketDBClient) GetDateByURL(url string) (time.Time, error) {
	var t time.Time
	query := fmt.Sprintf("select updated from urlinfo where url = '%s'", url)
	err := client.dbClient.QueryRow(query).Scan(&t)
	if err != nil {
		return time.Now(), err
	}
	return t, nil
}
