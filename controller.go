package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/kevvarlar/Chirpy/internal/database"
)

type error struct{
	Error string `json:"error"`
}

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

func (cfg *apiConfig) reset(res http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		res.WriteHeader(403)
		res.Write([]byte("403 Forbidden"))
		return
	}
	_, err := cfg.db.ResetUsers(req.Context())
	if err != nil {
		res.WriteHeader(500)
		return
	}
	cfg.fileserverHits.Store(0)
	res.WriteHeader(200)
	res.Write([]byte("OK"))
}

func (cfg *apiConfig) createChirp(res http.ResponseWriter, req *http.Request) {
	type parameters struct{
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		res.WriteHeader(500)
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
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(400)
		res.Write(jsonErr)
		return
	}
	profane := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}
	body := strings.Split(params.Body, " ")
	for i, w := range body {
		check := strings.ToLower(w)
		for _, p := range profane{
			if check == p {
				body[i] = "****"
			}
		}
	}
	cleanedBody := strings.Join(body, " ")
	chirp, err := cfg.db.CreateChirp(req.Context(), database.CreateChirpParams{
		Body: cleanedBody,
		UserID: params.UserID,
	})
	if err != nil {
		log.Printf("Error creating chirp: %s", err)
		jsonError, err := json.Marshal(error{
			Error: fmt.Sprintf("error while creating chirp: %s", err),
		})
		if err != nil {
			res.WriteHeader(500)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(400)
		res.Write(jsonError)
		return
	}
	response := Chirp(chirp)
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		res.WriteHeader(500)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(201)
	res.Write(jsonResponse)

}

func (cfg *apiConfig) getAllChirps (res http.ResponseWriter, req *http.Request) {
	response := []Chirp{}
	chirps, err := cfg.db.GetAllChirps(req.Context())
	if err != nil {
		log.Printf("Error getting all chirps: %s", err)
		jsonError, err := json.Marshal(error{
			Error: fmt.Sprintf("Error getting all chirps: %s", err),
		})
		if err != nil {
			res.WriteHeader(500)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(500)
		res.Write(jsonError)
		return
	}
	for _, chirp := range chirps {
		response = append(response, Chirp(chirp))
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error while marshalling response")
		res.WriteHeader(500)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(200)
	res.Write(jsonResponse)
}

func (cfg *apiConfig) getChirp(res http.ResponseWriter, req *http.Request) {
	path := req.PathValue("chirpID")
	chirpID, err := uuid.Parse(path)
	if err != nil {
		log.Printf("invalid uuid: %s", chirpID)
		jsonError, err := json.Marshal(error{
			Error: fmt.Sprintf("Invalid UUID %s", err),
		})
		if err != nil {
			res.WriteHeader(500)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(500)
		res.Write(jsonError)
		return
	}
	chirp, err := cfg.db.GetChirp(req.Context(), uuid.UUID(chirpID))
	if err != nil {
		if len(chirp.Body) == 0 {
			log.Printf("could not find chirp")
			res.WriteHeader(404)
			res.Write([]byte("404 Not Found"))
		} else {
			log.Printf("error getting chirp: %s", err)
			res.WriteHeader(500)
		}
		return
	}
	response := Chirp(chirp)
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		res.WriteHeader(500)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(200)
	res.Write(jsonResponse)
}

func (cfg *apiConfig) createUser(res http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}
	params := parameters{}
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&params); err != nil {
		log.Printf("Error decoding request body: %s", err)
		res.WriteHeader(500)
		return
	}
	user, err := cfg.db.CreateUser(req.Context(), params.Email)
	if err != nil {
		log.Printf("Error creating user: %s", err)
		jsonError, err := json.Marshal(error{
			Error: fmt.Sprintf("error while creating user: %s", err),
		})
		if err != nil {
			res.WriteHeader(500)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		res.Write(jsonError)
		res.WriteHeader(400)
		return
	}
	response := User(user)
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		res.WriteHeader(500)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(201)
	res.Write(jsonResponse)
}
