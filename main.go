package main

import (
	"database/sql"
	"flag"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	var user, password, dbName, conn string
	flag.StringVar(&user, "user", "root", "The DB username (required)")
	flag.StringVar(&password, "password", "@", "The DB password")
	flag.StringVar(&dbName, "db", "", "The DB name")
	flag.StringVar(&conn, "s", "", "The connection string")
	flag.Parse()

	if dbName == "" {
		flag.Usage()
		return
	}

	if conn == "" {

		InitConnection(user + ":" + password + "/" + dbName)

	} else {

		InitConnection(conn + "/" + dbName)
	}

	db := DB()
	OutputQueries(db, dbName)
}

//OutputQueries ....
func OutputQueries(db *sql.DB, dbName string) {

	tables := GetTableNames(db)

	for i := 0; i < len(tables); i++ {

		table := fmt.Sprintf("%s.%s", dbName, tables[i])
		pk := GetPrimaryKey(db, table)
		cols := GetColumnNames(db, table)
		fmt.Println("\n------ " + table + " ------")
		fmt.Println(CreateInsertQuery(table, cols))
		fmt.Println("")
		fmt.Println(CreateSelectAllQuery(table, cols))
		fmt.Println("")
		fmt.Println(CreateSelectOneQuery(table, pk, cols))
		fmt.Println("")
		fmt.Println(CreateUpdateQuery(table, pk, cols))
		fmt.Println("")
		fmt.Println(CreateDeleteQuery(table, pk))
		fmt.Println("")
	}
}

//GetTableNames ...
func GetTableNames(db *sql.DB) []string {

	var out []string

	rows, err := db.Query("SHOW TABLES")

	if err != nil {
		fmt.Println(err)
	}

	for rows.Next() {
		var tbl string
		rows.Scan(&tbl)
		out = append(out, tbl)
	}

	//fmt.Printf("%s", out)

	return out
}

//GetColumnNames ...
func GetColumnNames(db *sql.DB, table string) []string {
	row, err := db.Query("SELECT * FROM " + table + " LIMIT 0,1")

	if err != nil {
		fmt.Println(err)
	}

	cols, err := row.Columns()

	if err != nil {
		fmt.Println(err)
	}

	//fmt.Printf("%v", cols)

	return cols
}

//GetPrimaryKey ...
func GetPrimaryKey(db *sql.DB, table string) string {

	tpl := `SELECT k.COLUMN_NAME ` +
		`FROM information_schema.table_constraints t ` +
		`LEFT JOIN information_schema.key_column_usage k ` +
		`USING(constraint_name,table_schema,table_name) ` +
		`WHERE t.constraint_type='PRIMARY KEY' ` +
		`AND t.table_schema=DATABASE() ` +
		`AND CONCAT(t.table_schema, '.', t.table_name) ='%s';`

	var out []string

	rows, err := db.Query(fmt.Sprintf(tpl, table))

	if err != nil {
		fmt.Println(err)
	}

	for rows.Next() {
		var tbl string
		rows.Scan(&tbl)
		out = append(out, tbl)
	}

	//fmt.Printf("%s", out)

	return out[0]
}

//CreateInsertQuery ...
func CreateInsertQuery(table string, cols []string) string {

	tpl := "INSERT INTO %s (%s) VALUES (%s);"
	var vals []string
	numCols := len(cols)

	for i := 0; i < numCols; i++ {
		vals = append(vals, "?")
	}

	return fmt.Sprintf(tpl, table,
		strings.Join(cols, ", "),
		strings.Join(vals, ", "))
}

//CreateSelectAllQuery ...
func CreateSelectAllQuery(table string, cols []string) string {

	tpl := "SELECT %s FROM %s LIMIT ? OFFSET ?;"

	return fmt.Sprintf(tpl, strings.Join(cols, ", "), table)
}

//CreateSelectOneQuery ...
func CreateSelectOneQuery(table string, primaryKey string, cols []string) string {

	if primaryKey == "" {
		return ""
	}

	tpl := "SELECT %s FROM %s WHERE %s = ? LIMIT 1;"

	return fmt.Sprintf(tpl, strings.Join(cols, ", "), table, primaryKey)
}

//CreateUpdateQuery ...
func CreateUpdateQuery(table string, primaryKey string, cols []string) string {

	if primaryKey == "" {
		return ""
	}

	tpl := "UPDATE %s SET %s WHERE %s = ? LIMIT 1;"

	var vals []string
	numCols := len(cols)

	for i := 0; i < numCols; i++ {
		if cols[i] != primaryKey {
			vals = append(vals, cols[i]+" = ?")
		}

	}
	return fmt.Sprintf(tpl, table, strings.Join(vals, ", "), primaryKey)
}

//CreateDeleteQuery ...
func CreateDeleteQuery(table string, primaryKey string) string {

	if primaryKey == "" {
		return ""
	}

	tpl := "DELETE FROM %s WHERE %s = ? LIMIT 1;"

	return fmt.Sprintf(tpl, table, primaryKey)
}
