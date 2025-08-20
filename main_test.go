package traefik_plugin_response_cache_control_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	traefik_plugin_response_cache_control "github.com/NETCOREXT/traefik-plugin-response-cache-control"
)

func TestResponseCacheControl(t *testing.T) {
	cfg := traefik_plugin_response_cache_control.CreateConfig()
	cfg.Value = "public, max-age=3600"
	cfg.Override = true

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// 寫入一些內容，或至少設置狀態碼，這樣才會觸發 WriteHeader
		rw.WriteHeader(http.StatusOK)
	})

	handler, err := traefik_plugin_response_cache_control.New(ctx, next, cfg, "test")

	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)

	assertHeader(t, recorder, "Cache-Control", "public, max-age=3600")
}

func assertHeader(t *testing.T, recorder *httptest.ResponseRecorder, key, expected string) {
	t.Helper()

	if recorder.Header().Get(key) != expected {
		t.Errorf("invalid header value: %s", recorder.Header().Get(key))
	}
}
