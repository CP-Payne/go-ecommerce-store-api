package utils

import (
	"encoding/json"
	"net/http"

	"github.com/CP-Payne/ecomstore/internal/config"
	"go.uber.org/zap"
)

func RespondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		config.GetLogger().Error("error mashalling json: %v", zap.Error(err))
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(data)
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	type errorResponse struct {
		Error string `json:"error"`
	}
	RespondWithJson(w, code, errorResponse{
		Error: message,
	})
}
