package ramchi

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Etwodev/ramchi/router"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}
	defer resp.Body.Close()

	return resp, string(respBody)
}

func TestBasicServer(t *testing.T) {
	ts := New()

	pingAll := func(w http.ResponseWriter, r *http.Request) {
		res, _ := json.Marshal(map[string]string{"success": "ping"})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)
		if _, err := w.Write(res); err != nil {
			t.Fatal(err)
		}
	}

	testRoutes := func() []router.Route {
		return []router.Route{
			router.NewGetRoute("/ping", true, false, pingAll),
		}
	}

	testRouters := func() []router.Router {
		return []router.Router{
			router.NewRouter(testRoutes(), true),
		}
	}

	ts.LoadRouter(testRouters())

	instance := httptest.NewServer(ts.handler())
	defer instance.Close()

	if _, body := testRequest(t, instance, http.MethodGet, "/ping", nil); body != `{"success":"ping"}` {
		t.Fatalf(body)
	}
}
