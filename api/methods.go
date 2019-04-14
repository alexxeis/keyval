package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Value is a struct for JSON value object
type Value struct {
	Value string `json:"value"`
}

// Ttl is a struct for JSON ttl object
type Ttl struct {
	Ttl time.Duration `json:"ttl"`
}

// SetParams is a struct for JSON setParams object
type SetParams struct {
	Value string        `json:"value"`
	Ttl   time.Duration `json:"ttl"`
}

func (h *handler) Keys(w http.ResponseWriter, r *http.Request) {
	keys := h.storage.Keys()
	writeContent(w, keys)
}

func (h *handler) Set(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key, ok := vars["key"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var params SetParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if params.Value == "" || params.Ttl < 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.storage.Set(key, params.Value, params.Ttl*time.Millisecond)
}

func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key, ok := vars["key"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	val, err := h.storage.Get(key)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if val == "" {
		w.WriteHeader(http.StatusNotFound)
	} else {
		writeContent(w, Value{val})
	}
}

func (h *handler) Remove(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key, ok := vars["key"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.storage.Remove(key)
}

func (h *handler) Expire(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key, ok := vars["key"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var params Ttl
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if params.Ttl < 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if ok = h.storage.Expire(key, params.Ttl*time.Millisecond); !ok {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (h *handler) Hget(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key, ok := vars["key"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	field, ok := vars["field"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	val, err := h.storage.Hget(key, field)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if val == "" {
		w.WriteHeader(http.StatusNotFound)
	} else {
		writeContent(w, Value{val})
	}
}

func (h *handler) Hset(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key, ok := vars["key"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	field, ok := vars["field"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var params Value
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if params.Value == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.storage.Hset(key, field, params.Value); err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (h *handler) Hdel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key, ok := vars["key"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	field, ok := vars["field"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.storage.Hdel(key, field); err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
}
