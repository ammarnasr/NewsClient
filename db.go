package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const ( //Insert your own Credentials
	username        = "root"
	password        = "cogitoergosum"
	hostname        = "127.0.0.1:3306"
	dbname          = "newsa"
	tableName       = "linksa"
	column1         = "title"
	column2         = "href"
	maxIdleConns    = 20
	maxOpenConns    = 20
	connMaxLifetime = 5 * time.Minute
)

type newsEntry struct {
	//Entry row of SQL Table to be Created
	title    string
	hyperRef string
}

func addNewsEntry(db *sql.DB, ne newsEntry) {
	duplicate := findNewsEntry(db, ne)
	if duplicate {
		return
	}
	insertValueToColumnOfTable(db, tableName, column1, column2, ne.title, ne.hyperRef)
	fmt.Printf("Just Added %+v Successfully \n", ne)
}

func findNewsEntry(db *sql.DB, ne newsEntry) bool {
	return findValueFromTable(db, tableName, column1, column2, ne.title, ne.hyperRef)
}

func insertValueToColumnOfTable(db *sql.DB, tableName string, column1 string, column2 string, value1 string, value2 string) {
	insert, err := db.Query(fmt.Sprintf("INSERT INTO %s(%s, %s) VALUES('%s', '%s')", tableName, column1, column2, value1, value2))
	if err != nil {
		fmt.Println("Can't insert value: ", err)
	}
	defer insert.Close()
}

func findValueFromTable(db *sql.DB, tableName string, column1 string, column2 string, value1 string, value2 string) bool {
	rows := getAllRowsFromColumns(db, tableName, column1, column2)
	for rows.Next() {
		ne := newsEntry{}
		err := rows.Scan(&ne.title, &ne.hyperRef)
		if err != nil {
			fmt.Println("Can't Scan values: ", err)
		}
		if ne.title == value1 && ne.hyperRef == value2 {
			defer rows.Close()
			return true
		}
	}
	defer rows.Close()
	return false
}

func getAllRowsFromColumns(db *sql.DB, tableName string, column1 string, column2 string) *sql.Rows {
	rows, err := db.Query(fmt.Sprintf("SELECT %s, %s FROM %s", column1, column2, tableName))
	if err != nil {
		fmt.Println("Can't Select values: ", err)
	}
	return rows
}

func displayTable(db *sql.DB, tableName string) {
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s", tableName))
	if err != nil {
		fmt.Println("Can't Select values: ", err)
	}
	for rows.Next() {
		ne := newsEntry{}
		err := rows.Scan(&ne.title, &ne.hyperRef)
		if err != nil {
			fmt.Println("Can't Scan values: ", err)
		}
		fmt.Printf("%s \t %s \n", ne.title, ne.hyperRef)
	}
	defer rows.Close()
}

func openCreateDB(username string, password string, hostname string, dbname string) *sql.DB {
	createNewDB(username, password, hostname, dbname)
	return openExistingDB(username, password, hostname, dbname)
}

func dataSourceName(username string, password string, hostname string, dbname string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbname)
}

func createNewDB(username string, password string, hostname string, dbname string) {
	db, err := sql.Open("mysql", dataSourceName(username, password, hostname, ""))
	if err != nil {
		fmt.Printf("Error when opening DB: %s\n", err)
		return
	}

	cntx, cancelFunction := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunction()

	res, err := db.ExecContext(cntx, "CREATE DATABASE IF NOT EXISTS "+dbname)
	if err != nil {
		fmt.Printf("Error when here creating DB: %s\n", err)
		return
	}

	nRows, err := res.RowsAffected()
	if err != nil {
		fmt.Printf("Error when fetching Rows: %s\n", err)
		return
	}

	fmt.Printf("Rows affected: %d \n", nRows)
	defer db.Close()
}

func openExistingDB(username string, password string, hostname string, dbname string) *sql.DB {
	db, err := sql.Open("mysql", dataSourceName(username, password, hostname, dbname))
	if err != nil {
		fmt.Printf("Error when opening DB: %s\n", err)
		os.Exit(1)
	}
	pingConnection(db)
	configureConnection(db)
	return db
}

func pingConnection(db *sql.DB) {
	cntx, cancelFunction := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunction()

	err := db.PingContext(cntx)

	if err != nil {
		fmt.Printf("Error when Pinging DB: %s", err)
		os.Exit(1)
	}

	fmt.Println("Connected to DB Successfully")
}

func configureConnection(db *sql.DB) {
	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(connMaxLifetime)
}

func openCreateTable(db *sql.DB, tablename string, column1 string, column2 string) {
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s(%s text, %s text)", tablename, column1, column2)
	cntx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	res, err := db.ExecContext(cntx, query)
	if err != nil {
		fmt.Printf("Error %s when creating table", err)
		os.Exit(1)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		fmt.Printf("Error %s when fetching rows", err)
		os.Exit(1)
	}
	fmt.Printf("Rows affected when creating table: %d", rows)
}
