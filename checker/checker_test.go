package checker

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	c := &BTC{ts.URL}
	assert.True(t, c.Check())

	ts = newTestServer(false)
	defer ts.Close()
	c.url = ts.URL
	assert.False(t, c.Check())
}

func TestCheckerNetwork(t *testing.T) {
	n, err := New("dev")
	require.Nil(t, err)
	if !n.Check("BTC") {
		t.Log("Bitcoin network is down")
	}

}
