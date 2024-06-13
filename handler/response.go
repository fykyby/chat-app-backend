package handler

import (
	"encoding/json"
	"net/http"
)

func SendResponse(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	ok := statusCode >= 200 && statusCode < 300

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ok":      ok,
		"message": message,
		"data":    data,
	})
}
