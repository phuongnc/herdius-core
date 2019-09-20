package checker

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestServer(ok bool) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !ok {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		_, _ = fmt.Fprintln(w, "")
	}))
	return ts
}

func TestBTCChecker(t *testing.T) {
	ts := newTestServer(true)
	defer ts.Close()

	c := New("BTC", ts.URL)
	assert.True(t, c.Check())

	ts = newTestServer(false)
	defer ts.Close()
	c = New("BTC", ts.URL)
	assert.False(t, c.Check())
}
