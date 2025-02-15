package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	config "github.com/dgt4l/avito_shop/configs/avito_shop"
	"github.com/dgt4l/avito_shop/internal/avito_shop/auth"
	"github.com/dgt4l/avito_shop/internal/avito_shop/controller"
	handler "github.com/dgt4l/avito_shop/internal/avito_shop/handler"
	repository "github.com/dgt4l/avito_shop/internal/avito_shop/repository/pgsql"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	cfg, err := config.LoadConfig(".")
	if err != nil {
		logrus.Fatalf("Failed to load config: %v", err)
	}

	db, err := repository.NewRepository(&cfg)
	if err != nil {
		logrus.Fatalf("Failed to init db: %v", err)
	}

	auth := auth.NewAuth(&cfg)
	srv := controller.NewShopService(db, auth, &cfg)

	sh := handler.NewShopHandler(srv, auth)

	go sh.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := db.Close(); err != nil {
		logrus.Fatal(err)
	}
	if err := sh.Close(ctx); err != nil {
		logrus.Fatal(err)
	}
}
