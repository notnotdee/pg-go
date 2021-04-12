package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	server "github.com/dl-watson/pg-go/controller"
	db "github.com/dl-watson/pg-go/db/sqlc"
)

func main() {
	conn, err := sql.Open("postgres", "qfwfq:begaydocrime@localhost:5432/pg_go")
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	store := db.NewStore(conn)
	server.NewServer(*store)
}