package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mar4ehk0/notes"
	"github.com/mar4ehk0/notes/pkg/handler"
	"github.com/mar4ehk0/notes/pkg/repository"
	"github.com/mar4ehk0/notes/pkg/service"
	"github.com/spf13/viper"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	router := gin.Default()
	repository := repository.NewRepository()
	service := service.NewService(repository)
	handler := handler.New(router, service)

	srv := notes.Server{}
	if err := srv.Run(viper.GetString("port"), handler.InitRoutes()); err != nil {
		log.Fatalf("error occurred while running http server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
