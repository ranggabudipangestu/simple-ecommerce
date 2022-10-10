package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

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

	port := fmt.Sprintf(`:%s`, os.Getenv("APP_PORT"))
	log.Println("Listening Server On Port " + port)

	err = http.ListenAndServe(port, mux)

	if err != nil {
		log.Fatalln("Failed to start Server")
	}

}
