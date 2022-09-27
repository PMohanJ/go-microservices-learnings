package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/pmohanj/go-microservices/product-images/files"
	"github.com/pmohanj/go-microservices/product-images/handlers"
)

//var bindAddress = env.String("BIND_ADDRESS", false, ":9090", "Bind address for the server")
//var logLevel = env.String("LOG_LEVEL", false, "debug", "Log output level for the server [debug, info, trace]")
//var basePath = env.String("BASE_PATH", false, "./imagestore", "Base path to save images")

func main() {

	//	env.Parse()

	l := hclog.New(
		&hclog.LoggerOptions{
			Name:  "product-images",
			Level: hclog.LevelFromString("debug"),
		},
	)
	sl := l.StandardLogger(&hclog.StandardLoggerOptions{InferLevels: true})

	// create the storage class, use local storage
	// max filesize 5MB
	storage, err := files.NewLocal("./imagestore", 1024*1000*5)
	if err != nil {
		l.Error("Unable to crate storage", "error", err)
		os.Exit(1)
	}
	fh := handlers.NewFiles(l, storage)

	sm := mux.NewRouter()

	posth := sm.Methods(http.MethodPost).Subrouter()
	posth.HandleFunc("/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-z]{3}}", fh.ServeHTTP)

	//get files
	gh := sm.Methods(http.MethodGet).Subrouter()
	gh.Handle("/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-z]{3}}", http.StripPrefix("/images/", http.FileServer(http.Dir("./imagestore"))))

	s := http.Server{
		Addr:         ":9090",
		Handler:      sm,
		ErrorLog:     sl,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// start the server
	go func() {
		l.Info("Starting server", "bind_address")
		err := s.ListenAndServe()
		if err != nil {
			l.Error("Unable to start server", "error", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interrupt and shutdown server gracefully
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a singal is received
	sig := <-c
	l.Info("Shutting down the server with", "singal", sig)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}
