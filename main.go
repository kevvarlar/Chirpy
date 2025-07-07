package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	ServeMux := http.NewServeMux()
	ServeMux.Handle("/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	ServeMux.HandleFunc("/healthz", readiness)
	server := &http.Server{
		Handler: ServeMux,
		Addr: ":8080",
	}
	fmt.Println("Server running on http://localhost" + server.Addr)
	log.Fatal(server.ListenAndServe())
}
