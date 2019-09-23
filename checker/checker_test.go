package checker

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
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
	n := New("dev")
	require.True(t, n.Check("HER"))

	if !n.Check("BTC") {
		t.Log("Bitcoin network is down")
	}

}

func TestETHChecker(t *testing.T) {
	infuraID := os.Getenv("INFURAID")
	if len(infuraID) > 0 {
		c := &ETH{"https://ropsten.infura.io/v3/"}
		assert.True(t, c.Check())
	}
}

func TestETHCheckerFalse(t *testing.T) {
	c := &ETH{"https://wrong-ropsten.infura.io/v3/"}
	assert.False(t, c.Check())
}
