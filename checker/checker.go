package checker

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/herdius/herdius-core/config"
	"github.com/herdius/herdius-core/symbol"
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
	n.checker[symbol.BTC] = &BTC{c.CheckerBtcURL}
	n.checker[symbol.HER] = &HER{}
	n.checker[symbol.ETH] = &ETH{c.EthRPCURL}
	n.checker[symbol.HBTC] = &ETH{c.EthRPCURL}
	n.checker[symbol.HTZX] = &ETH{c.EthRPCURL}
	n.checker[symbol.HLTC] = &ETH{c.EthRPCURL}
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
	ethURL := e.url
	if strings.Index(e.url, ".infura.io") > -1 {
		ethURL += os.Getenv("INFURAID")
	}
	client, err := getEthClient(ethURL)
	if err != nil {
		log.Println(fmt.Sprintf("Failed to create ethereum client %s.", err))
		return false
	}
	_, err = client.BlockByNumber(context.Background(), nil)
	if err != nil {
		log.Println(fmt.Sprintf("Failed to get block detail %s.", err))
		return false
	}
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
	ethURL := hbtc.url
	if strings.Index(hbtc.url, ".infura.io") > -1 {
		ethURL += os.Getenv("INFURAID")
	}
	client, err := getEthClient(ethURL)
	if err != nil {
		log.Println(fmt.Sprintf("Failed to create ethereum client %s.", err))
		return false
	}
	_, err = client.BlockByNumber(context.Background(), nil)
	if err != nil {
		return false
	}
	return true
}

// HTZX represents checker for HTZX.
type HTZX struct {
	url string
}

func (htzx *HTZX) Check() bool {
	ethURL := htzx.url
	if strings.Index(htzx.url, ".infura.io") > -1 {
		ethURL += os.Getenv("INFURAID")
	}
	client, err := getEthClient(ethURL)
	if err != nil {
		log.Println(fmt.Sprintf("Failed to create ethereum client %s.", err))
		return false
	}
	_, err = client.BlockByNumber(context.Background(), nil)
	if err != nil {
		return false
	}
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
	ethURL := hltc.url
	if strings.Index(hltc.url, ".infura.io") > -1 {
		ethURL += os.Getenv("INFURAID")
	}
	client, err := getEthClient(ethURL)
	if err != nil {
		log.Println(fmt.Sprintf("Failed to create ethereum client %s.", err))
		return false
	}
	_, err = client.BlockByNumber(context.Background(), nil)
	if err != nil {
		return false
	}
	return true
}

func getEthClient(ethURL string) (*ethclient.Client, error) {
	clientChan := make(chan *ethclient.Client, 1)
	errChan := make(chan error, 1)
	go func() {
		client, err := ethclient.Dial(ethURL)
		if err != nil {
			errChan <- err
		}
		clientChan <- client
	}()
	select {
	case res := <-clientChan:
		return res, nil
	case ethErr := <-errChan:
		return nil, ethErr
	case <-time.After(2 * time.Second):
		return nil, fmt.Errorf(fmt.Sprintf("Could not connect to eth client due to timeout"))
	}
}
