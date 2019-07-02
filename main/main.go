package main

import (
	"fmt"
	"github.com/ravik-karn/Dataweave/factory"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/ravik-karn/Dataweave/config"
	"github.com/ravik-karn/Dataweave/handlers"
)

func main() {
	logger := logrus.New()
	config, err := config.New("config.json")
	if err != nil {
		logger.Fatal("read config file: %s", err)
	}

	factory := factory.New(config, *logger)

	r := mux.NewRouter()
	r.HandleFunc("/", handlers.HandleMain())
	r.HandleFunc("/products", handlers.HandleProduct(*logger, factory.ProductFetcher()))

	err = http.ListenAndServe(fmt.Sprintf(":%s", config.AppPort), nil)
	if err != nil {
		logger.Fatal("start server: %s", err)
	}
	logger.Infof("Running on: %s", config.AppPort)
}
