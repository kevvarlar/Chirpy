package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

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
	type cleanedBody struct{
		CleanedBody string `json:"cleaned_body"`
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
	profane := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}
	body := strings.Split(params.Body, " ")
	isValid := valid{
		Valid: true,
	}
	for i, w := range body {
		check := strings.ToLower(w)
		for _, p := range profane{
			if check == p {
				body[i] = "****"
				isValid.Valid = false
			}
		}
	}
	cleaned := strings.Join(body, " ")
	if !isValid.Valid {
		data, err := json.Marshal(cleanedBody{
			CleanedBody: cleaned,
		})
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			res.WriteHeader(500)
			return
		}

		res.WriteHeader(200)
		res.Write(data)
		return
	}
	data, err := json.Marshal(isValid)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		res.WriteHeader(500)
		return
	}

	res.WriteHeader(200)
	res.Write(data)

}

func (cfg *apiConfig) createUser(res http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}
	requestBody := parameters{}
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		log.Printf("Error decoding request body: %s", err)
		res.WriteHeader(500)
		return
	}
	user, err := cfg.db.CreateUser(req.Context(), requestBody.Email)
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
	responseBody := User(user)
	jsonResponseBody, err := json.Marshal(responseBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		res.WriteHeader(500)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(201)
	res.Write(jsonResponseBody)
}