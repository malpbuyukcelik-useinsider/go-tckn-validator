package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

type ValidationRequest struct {
	TCKN      string `json:"tckn"`
	Ad        string `json:"ad"`
	Soyad     string `json:"soyad"`
	DogumYili int    `json:"dogumYili"`
}

type ValidationResponse struct {
	Valid bool   `json:"valid"`
	Error string `json:"error,omitempty"`
}

func validateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ValidationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.TCKN) != 11 {
		sendResponse(w, ValidationResponse{Valid: false, Error: "TCKN must be 11 digits"})
		return
	}

	if !validateTCKN(req.TCKN) {
		sendResponse(w, ValidationResponse{Valid: false, Error: "Invalid TCKN format"})
		return
	}

	if req.Ad == "" || req.Soyad == "" || req.DogumYili == 0 {
		sendResponse(w, ValidationResponse{Valid: false, Error: "Ad, Soyad and DogumYili are required"})
		return
	}

	valid, err := validateWithNVI(req.TCKN, strings.ToUpper(req.Ad), strings.ToUpper(req.Soyad), req.DogumYili)
	if err != nil {
		sendResponse(w, ValidationResponse{Valid: false, Error: "Error validating TCKN: " + err.Error()})
		return
	}

	sendResponse(w, ValidationResponse{Valid: valid})
}

func sendResponse(w http.ResponseWriter, response ValidationResponse) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func enableCORS(handler http.HandlerFunc) http.HandlerFunc {
	allowedOrigins := map[string]bool{
		"https://extensions.shopifycdn.com":     true,
		"https://shopiapp-dev.myshopify.com":   true,
		"https://admin.shopify.com":            true,
		"https://shopify.com":                  true,
		"http://localhost:3000":                true,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Cache-Control, Pragma, X-Requested-With, X-HTTP-Method-Override, If-Match, If-None-Match, If-Modified-Since, If-Unmodified-Since")
			w.Header().Set("Access-Control-Max-Age", "3600")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		handler(w, r)
	}
}

func main() {
	http.HandleFunc("/validate", enableCORS(validateHandler))
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Printf("Server starting on :%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
} 