package main

import (
	"database/sql"
	"log"
	"time"
)

var connection *sql.DB

// Database is the common interface for database operations that can be used with
// types from schema 'gmail_notifier'.
//
// This should work with database/sql.DB and database/sql.Tx.
type Database interface {
	Exec(string, ...interface{}) (sql.Result, error)
	Query(string, ...interface{}) (*sql.Rows, error)
	QueryRow(string, ...interface{}) *sql.Row
}

//RowTimeStamps ...
type RowTimeStamps struct {
	Created int64 `json:"created,omitempty"`
	Updated int64 `json:"updated,omitempty"`
}

//createdNow ....
func (rt *RowTimeStamps) createdNow() {

	t := time.Now().Unix()
	rt.Created = t
	rt.Updated = t
}

//updatedNow ....
func (rt *RowTimeStamps) updatedNow() {

	t := time.Now().Unix()
	rt.Updated = t
}

//RecordStatus ...
type RecordStatus struct {
	exists, deleted bool
}

//Exists ...
func (rs *RecordStatus) Exists() bool {

	return rs.exists
}

//Deleted ...
func (rs *RecordStatus) Deleted() bool {

	return rs.deleted
}

//DBLog ...
var DBLog = func(string, ...interface{}) {}

//InitConnection ...
func InitConnection(ds string) {

	var err error

	//sql.ErrNoRows = errors.New("No results found")

	connection, err = sql.Open("mysql", ds)

	if err != nil {
		log.Panic(err)
	}

	if err = connection.Ping(); err != nil {
		log.Panic(err)
	}
}

//DB ...
func DB() *sql.DB {

	if connection == nil {
		log.Fatal("DB connection has not been initialized")
	}

	return connection
}
