package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/store")
	if err != nil {
		log.Println(err)
	}
	var (
		ccnum, date, cvv, exp string
		amount                float32
	)
	rows, err := db.Query("SELECT ccnum, date, amount, cvv, exp FROM transactions")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&ccnum, &date, &amount, &cvv, &exp)
		if err != nil {
			log.Println(err)
		}
		fmt.Println(ccnum, date, amount, cvv, exp)
		if rows.Err() != nil {
			log.Println(err)
		}
	}
}
