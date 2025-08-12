package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DbConnection *sql.DB

type Person struct {
	Name string
	Age  int
}

func main() {
	DbConnection, _ = sql.Open("sqlite3", "./example.sql")
	defer DbConnection.Close()
	cmd := `CREATE TABLE IF NOT EXISTS person(
		name STRING,
		age INT
	)`
	_, err := DbConnection.Exec(cmd)
	if err != nil {
		log.Fatalln(err)
	}

	tableName := "person; INSERT INTO person (name, age) VALUES ('Mr.X', 100)"
	cmd = fmt.Sprintf("SELECT * FROM %s", tableName)
	rows, _ := DbConnection.Query(cmd)
	defer rows.Close()
	var pp []Person
	for rows.Next() {
		var p Person
		err = rows.Scan(&p.Name, &p.Age)
		if err != nil {
			log.Println(err)
		}
		pp = append(pp, p)
	}
	err = rows.Err()
	if err != nil {
		log.Fatalln(err)
	}
	for _, p := range pp {
		fmt.Println(p.Name, p.Age)
	}
}
