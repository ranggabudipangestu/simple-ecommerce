package main

import (
	"log"
	"net/http"

	"github.com/ranggabudipangestu/simple-ecommerce/database"
	"github.com/ranggabudipangestu/simple-ecommerce/internal/factory"
)

func main() {

	db, err := database.Connect()

	if err != nil {
		log.Fatalln(err)
	}
	mux := http.NewServeMux()

	factory.RegisterHandlers(mux, db)

	log.Println("Listening...")

	err = http.ListenAndServe(":3000", mux)
	if err != nil {
		log.Fatalln("Failed to start Server")
	}

}
