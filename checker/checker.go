package checker

import (
	"io"
	"io/ioutil"
	"net/http"
)

// Checker checks availability of asset network
type Checker interface {
	Check() bool
}

// New returns new Checker for given asset.
func New(asset, url string) Checker {
	switch asset {
	case "BTC":
		return &BTC{url}
	}
	return nil
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
