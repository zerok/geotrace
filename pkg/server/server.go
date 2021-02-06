package server

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/zerok/geotrace/pkg/store"
)

const tsFormat = "2006-01-02T15:04:05Z0700"

type server struct {
	s                 store.Store
	apiKey            string
	router            chi.Router
	exposeMetrics     bool
	metricRegistry    *prometheus.Registry
	metricPointsTotal prometheus.GaugeFunc
}

type Configurator func(srv *server)

func ExposeMetrics(value bool) Configurator {
	return func(srv *server) {
		srv.exposeMetrics = value
	}
}

func New(s store.Store, apiKey string, options ...Configurator) http.Handler {
	r := chi.NewRouter()
	srv := &server{
		s:              s,
		apiKey:         apiKey,
		router:         r,
		metricRegistry: prometheus.NewRegistry(),
	}
	srv.metricPointsTotal = prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Name: "geotrace_points_total",
		Help: "Total number of stored points",
	}, srv.retrievePointsTotalMetric)
	srv.metricRegistry.MustRegister(srv.metricPointsTotal)

	for _, opt := range options {
		opt(srv)
	}

	if srv.exposeMetrics {
		r.Handle("/metrics", promhttp.HandlerFor(srv.metricRegistry, promhttp.HandlerOpts{}))
	}
	r.With(srv.requireAPIKey()).Post("/", srv.handlePost)
	return srv
}

func (s *server) retrievePointsTotalMetric() float64 {
	ctx := context.Background()
	count, _ := s.s.Count(ctx)
	return float64(count)
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
