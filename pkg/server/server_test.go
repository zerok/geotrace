package server_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
	"github.com/zerok/geotrace/pkg/server"
	"github.com/zerok/geotrace/pkg/store"
)

const testRequest = `{
  "locations": [
    {
      "type": "Feature",
      "geometry": {
        "type": "Point",
        "coordinates": [
          -122.030581, 
          37.331800
        ]
      },
      "properties": {
        "timestamp": "2015-10-01T08:00:00-0700",
        "altitude": 0,
        "speed": 4,
        "horizontal_accuracy": 30,
        "vertical_accuracy": -1,
        "motion": ["driving","stationary"],
        "pauses": false,
        "activity": "other_navigation",
        "desired_accuracy": 100,
        "deferred": 1000,
        "significant_change": "disabled",
        "locations_in_payload": 1,
        "battery_state": "charging",
        "battery_level": 0.89,
        "device_id": "",
        "wifi": ""
      }
    }
  ]
}
`

func TestServer(t *testing.T) {
	w := httptest.NewRecorder()
	st := store.NewCSVFileStore(afero.NewMemMapFs(), "data.csv")
	body := bytes.NewBufferString(testRequest)
	r := httptest.NewRequest(http.MethodPost, "/", body)
	s := server.New(st, "apikey")
	s.ServeHTTP(w, r)
	require.Equal(t, http.StatusUnauthorized, w.Code)

	body = bytes.NewBufferString(testRequest)
	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodPost, "/?apikey=apikey", body)
	s.ServeHTTP(w, r)
	require.Equal(t, http.StatusOK, w.Code)
}
