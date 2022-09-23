package main

import (
	"github.com/cucumberjaye/balanceAPI/pkg/handler"
	"github.com/cucumberjaye/balanceAPI/pkg/repository"
	"log"
	"net/http"
)

func main() {
	db, err := repository.NewPostgresDB()
	if err != nil {
		return
	}
	repos := repository.NewRepository(db)
	handlers := handler.NewHandler(repos)
	err = http.ListenAndServe(":8000", handlers.InitRoutes())
	if err != nil {
		log.Fatalf("error occured while listen and serve %s", err.Error())
	}
}
