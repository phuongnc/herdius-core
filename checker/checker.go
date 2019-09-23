package checker

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/herdius/herdius-core/config"
)

// Checker checks availability of asset network
type Checker interface {
	Check() bool
}

// New returns new Checker for given asset.
func New(env string) (*Network, error) {
	c := config.GetConfiguration(env)
	n := &Network{}
	n.checker = make(map[string]Checker)
	n.checker["BTC"] = &BTC{c.CheckerBtcURL}

	return n, nil
}

type Network struct {
	checker map[string]Checker
}

func (n *Network) Check(asset string) bool {
	c, ok := n.checker[asset]
	if !ok {
		return false
	}

	return c.Check()
}

// BTC represents checker for bitcoin.
type BTC struct {
	url string
}

// Check reports whether BTC network is available for processing.
func (b *BTC) Check() bool {
	if b.url == "" {
		return false
	}

	resp, err := http.Get(b.url)
	if err != nil {
		return false
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	_, _ = io.Copy(ioutil.Discard, resp.Body)
	return resp.StatusCode == http.StatusOK
}
