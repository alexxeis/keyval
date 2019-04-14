package client_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexxeis/keyval/api"
	"github.com/alexxeis/keyval/api/client"
)

func TestClient_Keys(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Error("wrong method ", r.Method)
		}
		if r.URL.String() != "/keys" {
			t.Error("wrong url:", r.URL.String())
		}
		w.Write([]byte(`["k1","k2"]`))
	}))
	defer server.Close()

	c := client.NewClient(server.URL, "go-client", server.Client())
	keys, err := c.Keys()
	if err != nil {
		t.Error(err)
	}

	if len(keys) != 2 {
		t.Errorf("keys len is %d", len(keys))
	}
	if keys[0] != "k1" {
		t.Error("keys[0] expepted k1, got ", keys[0])
	}
	if keys[1] != "k2" {
		t.Error("keys[1] expepted k2, got ", keys[1])
	}
}

func TestClient_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Error("wrong method ", r.Method)
		}
		if r.URL.String() != "/get/k" {
			t.Error("wrong url:", r.URL.String())
		}
		w.Write([]byte(`{"value":"v"}`))
	}))
	defer server.Close()

	c := client.NewClient(server.URL, "go-client", server.Client())
	val, err := c.Get("k")
	if err != nil {
		t.Error(err)
	}

	if val != "v" {
		t.Error("value expected v, got ", val)
	}
}

func TestClient_Set(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Error("wrong method ", r.Method)
		}
		if r.URL.String() != "/set/k" {
			t.Error("wrong url:", r.URL.String())
		}

		payload, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}
		r.Body.Close()

		var params api.SetParams
		if err = json.Unmarshal(payload, &params); err != nil {
			t.Error(err)
		}

		if params.Value != "v" || params.Ttl != 10 {
			t.Error("wrong payload ", string(payload))
		}
	}))
	defer server.Close()

	c := client.NewClient(server.URL, "go-client", server.Client())
	params := &api.SetParams{
		Value: "v",
		Ttl:   10,
	}
	if err := c.Set("k", params); err != nil {
		t.Error(err)
	}
}

func TestClient_Expire(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Error("wrong method ", r.Method)
		}
		if r.URL.String() != "/expire/k" {
			t.Error("wrong url:", r.URL.String())
		}

		payload, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}
		r.Body.Close()

		var ttl api.Ttl
		if err = json.Unmarshal(payload, &ttl); err != nil {
			t.Error(err)
		}

		if ttl.Ttl != 10 {
			t.Error("wrong payload ", string(payload))
		}
	}))
	defer server.Close()

	c := client.NewClient(server.URL, "go-client", server.Client())
	ttl := &api.Ttl{Ttl: 10}
	if err := c.Expire("k", ttl); err != nil {
		t.Error(err)
	}
}

func TestClient_Remove(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Error("wrong method ", r.Method)
		}
		if r.URL.String() != "/remove/k" {
			t.Error("wrong url:", r.URL.String())
		}
	}))
	defer server.Close()

	c := client.NewClient(server.URL, "go-client", server.Client())
	if err := c.Remove("k"); err != nil {
		t.Error(err)
	}
}

func TestClient_Hget(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Error("wrong method ", r.Method)
		}
		if r.URL.String() != "/hget/k/f" {
			t.Error("wrong url:", r.URL.String())
		}
		w.Write([]byte(`{"value":"v"}`))
	}))
	defer server.Close()

	c := client.NewClient(server.URL, "go-client", server.Client())
	val, err := c.Hget("k", "f")
	if err != nil {
		t.Error(err)
	}

	if val != "v" {
		t.Error("value expected v, got ", val)
	}
}

func TestClient_Hset(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Error("wrong method ", r.Method)
		}
		if r.URL.String() != "/hset/k/f" {
			t.Error("wrong url:", r.URL.String())
		}

		payload, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}
		r.Body.Close()

		var val api.Value
		if err = json.Unmarshal(payload, &val); err != nil {
			t.Error(err)
		}

		if val.Value != "v" {
			t.Error("wrong payload ", string(payload))
		}
	}))
	defer server.Close()

	c := client.NewClient(server.URL, "go-client", server.Client())
	if err := c.Hset("k", "f", "v"); err != nil {
		t.Error(err)
	}
}

func TestClient_Hdel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Error("wrong method ", r.Method)
		}
		if r.URL.String() != "/hdel/k/f" {
			t.Error("wrong url:", r.URL.String())
		}
	}))
	defer server.Close()

	c := client.NewClient(server.URL, "go-client", server.Client())
	if err := c.Hdel("k", "f"); err != nil {
		t.Error(err)
	}
}
