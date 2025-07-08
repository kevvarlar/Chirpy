package main

import (
	"encoding/json"
	"fmt"
	"log"
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

func validateChirp(res http.ResponseWriter, req *http.Request) {
	type parameters struct{
		Body string `json:"body"`
	}
	type valid struct{
		Valid bool `json:"valid"`
	}
	type error struct{
		Error string `json:"error"`
	}

	res.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		jsonErr, err := json.Marshal(error{
			Error: "Error Decoding parameters",
		})
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			res.WriteHeader(500)
			return
		}
		res.WriteHeader(500)
		res.Write(jsonErr)
		return
	}
	if len(params.Body) > 140 {
		jsonErr, err := json.Marshal(error{
			Error: "Chirp is too long",
		})
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			res.WriteHeader(500)
			return
		}
		res.WriteHeader(400)
		res.Write(jsonErr)
		return
	}
	data, err := json.Marshal(valid{
		Valid: true,
	})
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		res.WriteHeader(500)
		return
	}
	res.WriteHeader(200)
	res.Write(data)
}