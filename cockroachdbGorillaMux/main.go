package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/carlosm27/blog/cockroachdb-gorillamux/model"
	"github.com/gorilla/mux"
)

var (
	addr = flag.String("addr", "postgresql://postgres:Barcelona.1899@localhost:5432/expenses?sslmode=disable", "the address of the database")
)

func main() {
	flag.Parse()
	db, err := model.SetupDB()

	if err != nil {
		log.Println("Failed setting up database")
	}

	router := mux.NewRouter()

	server := NewServer(db)

	server.RegisterRouter(router)
	log.Fatal(http.ListenAndServe(":8000", router))

}
