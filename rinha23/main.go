package main

import (
	"log"
	"net/http"
	"os"
	"rinha23/controllers"
	"rinha23/helpers"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	godotenv.Load(".env")
	helpers.SetupLog()	
	logrus.Info(">>>>>>>>>>   " + os.Getenv("SERVER_NAME")  + "   <<<<<<<<<< ")

	// DB
	helpers.TestDBConnection()

	// Redis
	helpers.TestRedisConnection()

	routes := controllers.SetupRoutes()

	logrus.Debug("Starting server at port ", os.Getenv("WEB_PORT"))
    if err := http.ListenAndServe(":" + os.Getenv("WEB_PORT"), routes); err != nil {
        log.Fatal(err)
    }
}


