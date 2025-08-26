package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

// Payload recebido do app Android
type Ingest struct {
	Device string                   `json:"device"`
	Source string                   `json:"source"`
	Events []map[string]interface{} `json:"events"`
}

// Handler é o ponto de entrada da Vercel Function
func Handler(w http.ResponseWriter, r *http.Request) {
	// CORS básico (útil p/ testes web; para Android nativo não interfere)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	switch r.URL.Path {
	case "/":
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte("its working on vercel"))
		return

	case "/health/ingest":
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		// limite defensivo de 5MB
		r.Body = http.MaxBytesReader(w, r.Body, 5<<20)

		var p Ingest
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		log.Printf("received %d events from %s/%s", len(p.Events), p.Device, p.Source)

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"ok":       true,
			"received": len(p.Events),
		})
		return

	default:
		http.NotFound(w, r)
		return
	}
}
