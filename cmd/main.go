package main

import (
	todo "github.com/Hudayberdyyev/Rest_ToDo"
	"github.com/Hudayberdyyev/Rest_ToDo/pkg/handler"
	"github.com/Hudayberdyyev/Rest_ToDo/pkg/repository"
	"github.com/Hudayberdyyev/Rest_ToDo/pkg/service"
	"github.com/spf13/viper"
	"log"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s\n", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host: "localhost",
		Port: 5432,
		Username: "postgres",
		Password: "qwerty",
		DBName: "postgres",
		SSLMode: "disable",
	})

	if err != nil {
		log.Fatalf("error to initialize db: %s\n", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(todo.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
