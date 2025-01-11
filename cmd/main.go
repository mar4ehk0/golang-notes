package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mar4ehk0/notes"
	"github.com/mar4ehk0/notes/pkg/handler"
)

func main() {
	router := gin.Default()
	handler := handler.New(router)

	srv := notes.Server{}
	if err := srv.Run("8080", handler.InitRoutes()); err != nil {
		log.Fatalf("error occurred while running http server: %s", err.Error())
	}
}
