package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

const username = "root"
const password = "cogitoergosum"
const dbName = "news"
const tableName = "links"
const column1 = "title"
const column2 = "href"

type newsEntry struct {
	//Entry row of SQL Table to be Created
	title    string
	hyperRef string
}

func createDBObject(username string, password string, dbName string) *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", username, password, dbName))
	if err != nil {
		fmt.Println("Can't creat DB object: ", err)
	}
	return db
}

func addNewsEntry(db *sql.DB, ne newsEntry) {
	duplicate := findNewsEntry(db, ne)
	if duplicate {
		fmt.Println("Duplicate Value")
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
			return true
		}
	}
	return false
}

func getAllRowsFromColumns(db *sql.DB, tableName string, column1 string, column2 string) *sql.Rows {
	rows, err := db.Query(fmt.Sprintf("SELECT %s, %s FROM %s", column1, column2, tableName))
	if err != nil {
		fmt.Println("Can't Select values: ", err)
	}
	return rows
}
