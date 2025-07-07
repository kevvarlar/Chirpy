package main

import (
	"net/http"
	"strconv"

)

func readiness(res http.ResponseWriter, _ *http.Request) {
	res.WriteHeader(200)
	res.Write([]byte("OK"))
}

func (cfg *apiConfig) metrics(res http.ResponseWriter, _ *http.Request) {
	hits := cfg.fileserverHits.Load()
	stringHits := strconv.FormatInt(int64(hits), 10)
	res.WriteHeader(200)
	res.Write([]byte("Hits: " + stringHits))
}

func (cfg *apiConfig) reset(res http.ResponseWriter, _ *http.Request) {
	cfg.fileserverHits.Store(0)
	res.WriteHeader(200)
	res.Write([]byte("OK"))
}