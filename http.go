package main

import "net/http"

func readiness(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(200)
	res.Write([]byte("OK"))
}