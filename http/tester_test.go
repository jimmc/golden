package http_test

import (
  "net/http"
  "testing"

  goldenhttp "github.com/jimmc/golden/http"
)

type handler struct {}

func (h *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
  w.WriteHeader(http.StatusOK)
  w.Write([]byte("Sample response"))
}

func createTestHandler(r *goldenhttp.Tester) http.Handler {
  return &handler{}
}

func TestHttpTester(t *testing.T) {
  request := func() (*http.Request, error) {
    return http.NewRequest("GET", "/api/foo/", nil)
  }
  r := goldenhttp.NewTester(createTestHandler)
  if err := r.Run("foo", request); err != nil {
    t.Fatalf("Error in Run: %s", err)
  }
}
