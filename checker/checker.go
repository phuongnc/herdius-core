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
func New(env string) *Network {
	c := config.GetConfiguration(env)
	n := &Network{}
	n.checker = make(map[string]Checker)
	n.checker["BTC"] = &BTC{c.CheckerBtcURL}
	n.checker["HER"] = &HER{}

	return n
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

// ETH represents checker for Ethereum.
type ETH struct {
	url string
}

// Check reports whether ETH network is available for processing.
func (e *ETH) Check() bool {
	// TODO: implement later.
	return true
}

// HER represents checker for Herdius.
type HER struct{}

// Check reports whether Herdius network is available for processing.
func (h *HER) Check() bool {
	// Herdius network is always available
	return true
}

// HBTC represents checker for HBTC.
type HBTC struct {
	url string
}

func (hbtc *HBTC) Check() bool {
	// TODO: implement later.
	return true
}

// HTZX represents checker for HTZX.
type HTZX struct {
	url string
}

func (htzx *HTZX) Check() bool {
	// TODO: implement later.
	return true
}

// TZX represents checker for Tezos.
type TZX struct {
	url string
}

func (tzx *TZX) Check() bool {
	// TODO: implement later.
	return true
}

// LTC represents checker for lite coin.
type LTC struct {
	url string
}

func (l *LTC) Check() bool {
	// TODO: implement later.
	return true
}

// HLTC represents checker for HLTC.
type HLTC struct {
	url string
}

func (hltc *HLTC) Check() bool {
	// TODO: implement later.
	return true
}
