package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi"

	"github.com/zerok/geotrace/pkg/store"
)

const tsFormat = "2006-01-02T15:04:05Z0700"

type server struct {
	s      store.Store
	apiKey string
	router chi.Router
}

func New(s store.Store, apiKey string) http.Handler {
	r := chi.NewRouter()
	srv := &server{
		s:      s,
		apiKey: apiKey,
		router: r,
	}
	r.With(srv.requireAPIKey()).Post("/", srv.handlePost)
	return srv
}

func (s *server) requireAPIKey() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if s.apiKey == "" {
				next.ServeHTTP(w, r)
				return
			}
			q := r.URL.Query()
			apiKey := q.Get("apikey")
			if apiKey != s.apiKey {
				http.Error(w, "Invalid api key", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) handlePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST supported", http.StatusMethodNotAllowed)
		return
	}
	req := Request{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	for _, l := range req.Locations {
		ts, err := time.Parse(tsFormat, l.Properties.Timestamp)
		if err != nil {
			http.Error(w, "Bad timestamp", http.StatusBadRequest)
			return
		}
		if err := s.s.Add(ts, l.Geometry.Coordinates, l.Properties.DeviceID); err != nil {
			http.Error(w, "Failed to store", http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&Response{
		Result: "ok",
	})
}
