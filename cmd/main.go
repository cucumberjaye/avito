package main

import (
	"github.com/cucumberjaye/avito/pkg/handler"
	"github.com/cucumberjaye/avito/pkg/repository"
	"log"
	"net/http"
)

func main() {
	/*db, err := repository.NewPostgresDB()
	if err != nil {
		return
	}*/
	repos := repository.NewRepository(nil)
	handlers := handler.NewHandler(repos)
	err := http.ListenAndServe(":8000", handlers.InitRoutes())
	if err != nil {
		log.Fatalf("error occured while listen and serve %s", err.Error())
	}
}
