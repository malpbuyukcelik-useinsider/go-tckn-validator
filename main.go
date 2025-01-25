package main

import (
	"encoding/json"
	"log"
	"net/http"
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

	// Basic validation
	if len(req.TCKN) != 11 {
		sendResponse(w, ValidationResponse{
			Valid: false,
			Error: "TCKN must be 11 digits",
		})
		return
	}

	// First do basic validation
	if !validateTCKN(req.TCKN) {
		sendResponse(w, ValidationResponse{
			Valid: false,
			Error: "Invalid TCKN format",
		})
		return
	}

	// Validate required fields
	if req.Ad == "" || req.Soyad == "" || req.DogumYili == 0 {
		sendResponse(w, ValidationResponse{
			Valid: false,
			Error: "Ad, Soyad and DogumYili are required",
		})
		return
	}

	// Convert names to uppercase for NVI service
	ad := strings.ToUpper(req.Ad)
	soyad := strings.ToUpper(req.Soyad)

	// Then validate with NVI service
	valid, err := validateWithNVI(req.TCKN, ad, soyad, req.DogumYili)
	if err != nil {
		sendResponse(w, ValidationResponse{
			Valid: false,
			Error: "Error validating TCKN: " + err.Error(),
		})
		return
	}

	sendResponse(w, ValidationResponse{
		Valid: valid,
		Error: "",
	})
}

func sendResponse(w http.ResponseWriter, response ValidationResponse) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/validate", validateHandler)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
} 