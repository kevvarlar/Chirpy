package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}
	ServeMux := http.NewServeMux()

	ServeMux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(".")))))

	ServeMux.HandleFunc("/healthz", readiness)
	ServeMux.HandleFunc("/metrics", apiCfg.metrics)
	ServeMux.HandleFunc("/reset", apiCfg.reset)
	server := &http.Server{
		Handler: ServeMux,
		Addr: ":8080",
	}
	fmt.Println("Server running on http://localhost" + server.Addr)
	log.Fatal(server.ListenAndServe())
}
