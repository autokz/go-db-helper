package db

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"

	_ "github.com/jackc/pgx/v4/stdlib"
)

const pingIntervalDefault = 5 * time.Second
const initTimeoutDefault = 60 * time.Second

type PostgreConfigParams interface {
	GetHost() string
	GetPort() string
	GetDbName() string
	GetUser() string
	GetPassword() string
}

type Db struct {
	dbConn       *sqlx.DB
	dbChan       chan *sqlx.DB
	done         chan bool
	params       PostgreConfigParams
	pingInterval time.Duration
	initTimeout  time.Duration
}

// GetDb if initTimeout/pingInterval <= 0, then will be used default time Duration for this params.
func GetDb(params PostgreConfigParams, initTimeout, pingInterval time.Duration) *Db {
	db := &Db{
		params:       params,
		pingInterval: pingIntervalDefault,
		initTimeout:  initTimeoutDefault,
		dbChan:       make(chan *sqlx.DB, 1),
		done:         make(chan bool, 1),
	}

	if initTimeout > 0 {
		db.initTimeout = initTimeout
	}

	if pingInterval > 0 {
		db.pingInterval = pingInterval
	}

	db.dbConn = db.getDbWithTicker()
	return db
}

func (db Db) GetConn() *sqlx.DB {
	if db.dbConn == nil {
		panic("Db connection not initialized. Call GetDb(params)")
	}
	return db.dbConn
}

func (db Db) getDbWithTicker() *sqlx.DB {
	go db.fatalAfterTime(db.initTimeout)

	go db.pingDb(db.pingInterval)

	dbConn := <-db.dbChan
	close(db.done)
	close(db.dbChan)

	return dbConn
}

func (db Db) fatalAfterTime(duration time.Duration) {
	tickerFail := time.NewTicker(duration)
	for {
		select {
		case <-db.done:
			return
		case <-tickerFail.C:
			panic("Db init timeout exceeded")
		}
	}
}

func (db *Db) pingDb(duration time.Duration) {
	var err error
	dbConn, err := db.getDbConn()
	if err == nil {
		db.dbChan <- dbConn
		return
	}
	log.Printf("Get dbChan error=%v\n", err)

	ticker := time.NewTicker(duration)

	for {
		<-ticker.C
		dbConn, err = db.getDbConn()
		if err != nil {
			log.Printf("Get dbChan error=%v\n", err)
			continue
		}
		db.dbChan <- dbConn
		return
	}
}

func (db Db) getDbConn() (*sqlx.DB, error) {
	params := db.params
	var err error
	dsn := "postgres://" +
		params.GetUser() + ":" + params.GetPassword() + "@" + params.GetHost() + ":" + params.GetPort() + "/" +
		params.GetDbName() + "?sslmode=disable"
	dbConn, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	_, err = dbConn.Exec("SELECT 1;")
	if err != nil {
		return nil, err
	}

	return dbConn, nil
}
