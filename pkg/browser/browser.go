package browser

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/zerok/geotrace/pkg/store"
)

type Browser struct {
	store   store.Store
	router  chi.Router
	webRoot string
}

type Option func(*Browser)

func WithWebRoot(path string) Option {
	return func(b *Browser) {
		b.webRoot = path
	}
}

func New(st store.Store, options ...Option) *Browser {
	b := &Browser{
		store: st,
	}
	for _, o := range options {
		o(b)
	}
	router := chi.NewRouter()
	router.Get("/api/v1/latest", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if err := st.Open(ctx); err != nil {
			http.Error(w, "Failed to open datastore", http.StatusInternalServerError)
			return
		}
		defer st.Close(ctx)
		start := time.Now().Add(time.Hour * -24)
		traces, err := st.GetTracesSince(ctx, start)
		if err != nil {
			http.Error(w, "Failed to fetch traces", http.StatusInternalServerError)
		}
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(traces)
	})
	if b.webRoot != "" {
		router.Mount("/", http.FileServer(http.Dir(b.webRoot)))
	}
	b.router = router
	return b
}

func (b *Browser) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b.router.ServeHTTP(w, r)
}
