package ramchi

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	c "github.com/Etwodev/ramchi/config"
	"github.com/Etwodev/ramchi/log"
	"github.com/Etwodev/ramchi/middleware"
	"github.com/Etwodev/ramchi/router"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader, headers map[string]string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	return resp, string(respBody)
}

func TestBasicServer(t *testing.T) {
	const ERROR_STATUS_CODE = 418
	const ERROR_MESSAGE = "Example error has occurred"
	const ERROR_RESPONSE = "test error pass-through"

	defaultConfig := &c.Config{
		Port:                 "7000",
		Address:              "127.0.0.1",
		Experimental:         false,
		ReadTimeout:          15,
		WriteTimeout:         15,
		IdleTimeout:          60,
		LogLevel:             "debug",
		MaxHeaderBytes:       1048576,
		EnableTLS:            false,
		TLSCertFile:          "",
		TLSKeyFile:           "",
		ShutdownTimeout:      5,
		EnableCORS:           true,
		AllowedOrigins:       []string{"http://example.com"},
		EnableRequestLogging: true,
	}

	err := c.Create(defaultConfig)
	if err != nil {
		t.Fatal(err)
	}

	ts := New()

	// Simulate what Start() would do — apply middlewares based on config
	if c.EnableCORS() && len(c.AllowedOrigins()) > 0 {
		corsMw := middleware.NewCORSMiddleware(c.AllowedOrigins())
		ts.LoadMiddleware([]middleware.Middleware{corsMw})
	}
	if c.EnableRequestLogging() {
		loggingMw := middleware.NewLoggingMiddleware(ts.Logger())
		ts.LoadMiddleware([]middleware.Middleware{loggingMw})
	}

	// Handlers
	pingAll := func(w http.ResponseWriter, r *http.Request) {
		// Confirm logger middleware injected the logger
		logger := log.FromContext(r.Context())
		if logger == nil {
			t.Error("Expected logger to be injected into context via middleware")
		}

		res, _ := json.Marshal(map[string]string{"success": "ping"})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)

		_, err := w.Write(res)
		if err != nil {
			t.Fatal(err)
		}

	}

	errorAll := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(ERROR_STATUS_CODE)
		response := map[string]string{
			"error":   ERROR_MESSAGE,
			"details": ERROR_RESPONSE,
		}

		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			t.Fatal(err)
		}

	}

	// Routes
	testRoutes := func() []router.Route {
		return []router.Route{
			router.NewGetRoute("ping", true, false, pingAll, nil),
			router.NewGetRoute("error", true, false, errorAll, nil),
		}
	}

	// Routers
	testRouters := func() []router.Router {
		return []router.Router{
			router.NewRouter("test", testRoutes(), true, nil),
		}
	}
	ts.LoadRouter(testRouters())

	instance := httptest.NewServer(ts.handler())
	defer instance.Close()

	// ─── Test /ping ─────────────────────────────────────────────────────
	resp, body := testRequest(t, instance, http.MethodGet, "/test/ping", nil, nil)
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected 201 status, got %d", resp.StatusCode)
	}
	if body != `{"success":"ping"}` {
		t.Fatalf("Unexpected ping response: %s", body)
	}

	// ─── Test /error ────────────────────────────────────────────────────
	resp, body = testRequest(t, instance, http.MethodGet, "/test/error", nil, nil)
	if resp.StatusCode != ERROR_STATUS_CODE {
		t.Errorf("Expected status %d, got %d", ERROR_STATUS_CODE, resp.StatusCode)
	}
	expected := `{"details":"test error pass-through","error":"Example error has occurred"}` + "\n"
	if body != expected {
		t.Fatalf("Unexpected error response: %s", body)
	}

	// ─── Test CORS ──────────────────────────────────────────────────────
	req, _ := http.NewRequest(http.MethodOptions, instance.URL+"/test/ping", nil)
	req.Header.Set("Origin", "http://example.com")
	req.Header.Set("Access-Control-Request-Method", "GET")
	corsResp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if got := corsResp.Header.Get("Access-Control-Allow-Origin"); got != "http://example.com" && got != "*" {
		t.Errorf("CORS header mismatch: got '%s', expected 'http://example.com'", got)
	}

	if got := corsResp.Header.Get("Access-Control-Allow-Methods"); got == "" {
		t.Errorf("CORS Allow-Methods header missing")
	}
}
