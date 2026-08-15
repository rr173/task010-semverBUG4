package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Regression: decodeJSON must reject request bodies that contain
// trailing content after a valid JSON object, whether the trailing
// content is valid JSON or invalid garbage.

func TestDecodeJSONRejectsTrailingGarbage(t *testing.T) {
	api := New()
	srv := httptest.NewServer(api.Handler())
	defer srv.Close()

	cases := []struct {
		name string
		body string
	}{
		{"trailing-text", `{"version":"1.2.3"}garbage`},
		{"trailing-whitespace-and-text", `{"version":"1.2.3"} oops`},
		{"trailing-number", `{"version":"1.2.3"}123`},
		{"trailing-array", `{"version":"1.2.3"}[1,2]`},
	}
	for _, c := range cases {
		resp, err := http.Post(srv.URL+"/api/parse", "application/json",
			strings.NewReader(c.body))
		if err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("%s: status=%d, want 400", c.name, resp.StatusCode)
		}
	}
}
