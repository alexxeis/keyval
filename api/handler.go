package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/alexxeis/keyval/storage"
)

// handler is a API handler struct
type handler struct {
	storage storage.Storage
}

// NewHandler returns new API handler
func NewHandler(s storage.Storage) *handler {
	return &handler{s}
}

// writeContent writes payload to writer
func writeContent(w http.ResponseWriter, content interface{}) {
	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(content); err != nil {
		log.Print(err)
	}
}
