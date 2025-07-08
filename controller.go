package main

import (
	"fmt"
	"net/http"
	// "strconv"
)

func readiness(res http.ResponseWriter, _ *http.Request) {
	res.WriteHeader(200)
	res.Write([]byte("OK"))
}

func (cfg *apiConfig) metrics(res http.ResponseWriter, _ *http.Request) {
	hits := cfg.fileserverHits.Load()
	// stringHits := strconv.FormatInt(int64(hits), 10)
	res.WriteHeader(200)
	res.Header().Set("Content-Type", "text/html")
	metricsHtml := fmt.Sprintf(`
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>
	`, hits)
	res.Write([]byte(metricsHtml))
}

func (cfg *apiConfig) reset(res http.ResponseWriter, _ *http.Request) {
	cfg.fileserverHits.Store(0)
	res.WriteHeader(200)
	res.Write([]byte("OK"))
}