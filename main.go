package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	"github.com/kevvarlar/Chirpy/internal/database"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db *database.Queries
	platform string
}

func main() {
	godotenv.Load()
	platform := os.Getenv("PLATFORM")
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Error opening connection to database")
	}
	dbQueries := database.New(db)
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db: dbQueries,
		platform: platform,
	}
	ServeMux := http.NewServeMux()

	ServeMux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(".")))))

	ServeMux.HandleFunc("GET /api/healthz", readiness)
	ServeMux.HandleFunc("POST /api/chirps", apiCfg.createChirp)
	ServeMux.HandleFunc("GET /api/chirps", apiCfg.getAllChirps)
	ServeMux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.getChirp)
	ServeMux.HandleFunc("POST /api/users", apiCfg.createUser)
	ServeMux.HandleFunc("GET /admin/metrics", apiCfg.metrics)
	ServeMux.HandleFunc("POST /admin/reset", apiCfg.reset)
	server := &http.Server{
		Handler: ServeMux,
		Addr: ":8080",
	}
	fmt.Println("Server running on http://localhost" + server.Addr)
	log.Fatal(server.ListenAndServe())
}
