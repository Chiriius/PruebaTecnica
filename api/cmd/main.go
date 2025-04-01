package main

import (
	"context"
	"log"
	"os"
	"prueba_tecnica/api/server"
	"time"

	_ "prueba_tecnica/api/docs" //

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @title API de Gestión de Eventos
// @version 1.0
// @description API para la gestión y clasificación de eventos
// @host localhost:8080
// @BasePath /api/v1
func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	logger := logrus.StandardLogger()
	logger.SetFormatter(&logrus.JSONFormatter{})

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		dbUrl = "mongodb://mongodb:27017"
	}

	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbUrl))
	if err != nil {
		log.Fatal(err)
	}

	srv := server.NewServer(client, logger)
	srv.Run()
}
