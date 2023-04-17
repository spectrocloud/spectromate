package endpoints

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

// NewHandlerContext returns a new CounterRoute with a database connection.
func NewHealthHandlerContext(ctx context.Context) *HealthRoute {
	return &HealthRoute{ctx}
}

func (health *HealthRoute) HealthHTTPHandler(writer http.ResponseWriter, request *http.Request) {
	log.Debug().Msg("Health check request received.")
	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	var payload []byte

	value, err := health.getHandler(request)
	if err != nil {
		log.Debug().Msg("Error getting counter value.")
		http.Error(writer, "Error getting counter value.", http.StatusInternalServerError)
	}
	writer.WriteHeader(http.StatusOK)
	payload = value
	_, err = writer.Write([]byte(payload))
	if err != nil {
		log.Error().Err(err).Msg("Error writing response to the health endpoint.")
	}
}

// getHandler returns a health check response.
func (health *HealthRoute) getHandler(r *http.Request) ([]byte, error) {
	return json.Marshal(map[string]string{"status": "OK"})
}
