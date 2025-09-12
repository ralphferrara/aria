package responses

import (
	"encoding/json"
	"log"
	"net/http"
)

//||------------------------------------------------------------------------------------------------||
//|| Success Response
//||------------------------------------------------------------------------------------------------||

func Success(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := map[string]any{
		"success": true,
		"data":    data,
	}

	json.NewEncoder(w).Encode(resp)
}

//||------------------------------------------------------------------------------------------------||
//|| Error Response
//||------------------------------------------------------------------------------------------------||

func Error(w http.ResponseWriter, status int, message string) {
	log.Printf("❌ %d ERROR: %s\n", status, message)

	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]any{
		"success": false,
		"message": message,
	})
}
