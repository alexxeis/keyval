package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexxeis/keyval/api"
	"github.com/alexxeis/keyval/cluster"
	"github.com/gorilla/mux"
)

func main() {
	port := flag.String("p", "8000", "listening port")
	count := flag.Int("c", 100, "cluster instances count")
	cleanInterval := flag.Int64("i", 1000, "clean interval in milliseconds")
	flag.Parse()

	if *port == "" || *count < 1 || *cleanInterval < 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	ci := time.Duration(*cleanInterval) * time.Millisecond
	c := cluster.NewCluster(*count, ci)
	handler := api.NewHandler(c)

	router := mux.NewRouter()
	router.HandleFunc("/api/keys", handler.Keys).Methods(http.MethodGet)
	router.HandleFunc("/api/get/{key}", handler.Get).Methods(http.MethodGet)
	router.HandleFunc("/api/set/{key}", handler.Set).Methods(http.MethodPost)
	router.HandleFunc("/api/remove/{key}", handler.Remove).Methods(http.MethodPost)
	router.HandleFunc("/api/expire/{key}", handler.Expire).Methods(http.MethodPost)
	router.HandleFunc("/api/hget/{key}/{field}", handler.Hget).Methods(http.MethodGet)
	router.HandleFunc("/api/hset/{key}/{field}", handler.Hset).Methods(http.MethodPost)
	router.HandleFunc("/api/hdel/{key}/{field}", handler.Hdel).Methods(http.MethodPost)

	// TODO: graceful shutdown
	log.Fatal(http.ListenAndServe(":"+*port, router))
}
