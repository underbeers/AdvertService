package main

import (
	"github.com/sirupsen/logrus"
	"github.com/underbeers/AdvertService/pkg/handler"
	"github.com/underbeers/AdvertService/pkg/repository"
	advertService "github.com/underbeers/AdvertService/pkg/server"
	"github.com/underbeers/AdvertService/pkg/service"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	cfg := repository.GetConfig()
	db, err := repository.NewPostgresDB(*cfg.DB)
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos, cfg)
	handlers := handler.NewHandler(services)

	srv := new(advertService.Server)
	if err := srv.Run(handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error occured while running http server: %s", err.Error())
	}
}
