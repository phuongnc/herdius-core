package sync

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"strconv"
	"strings"

	"github.com/herdius/herdius-core/p2p/log"
	"github.com/herdius/herdius-core/storage/state/statedb"
	"github.com/herdius/herdius-core/symbol"
)

// HTokenSyncer syncs all HToken external accounts
type HTokenSyncer struct {
	RPC                  string
	symbol, parentSymbol string
	syncer               *ExternalSyncer
}

func newHBTCSyncer() *HTokenSyncer {
	return newHTokenSyncer(symbol.HBTC, symbol.ETH)
}

func newHLTCSyncer() *HTokenSyncer {
	return newHTokenSyncer(symbol.HLTC, symbol.ETH)
}

func newHBNBSyncer() *HTokenSyncer {
	return newHTokenSyncer(symbol.HBNB, symbol.ETH)
}

func newHXTZSyncer() *HTokenSyncer {
	return newHTokenSyncer(symbol.HXTZ, symbol.ETH)
}

func newHTokenSyncer(hSymbol, parentSymbol string) *HTokenSyncer {
	h := &HTokenSyncer{symbol: hSymbol, parentSymbol: parentSymbol}
	h.syncer = newExternalSyncer(hSymbol)

	return h
}

// GetExtBalance ...
func (hs *HTokenSyncer) GetExtBalance() error {
	// If ETH account exists
	parentSymbolAccount, ok := hs.syncer.Account.EBalances[hs.parentSymbol]
	if !ok {
		log.Warn().Msgf("%s depends on %[2]s, but no %[2]s account available", hs.symbol, hs.parentSymbol)
		return fmt.Errorf("%s account does not exists", hs.parentSymbol)
	}

	hTokenAccount, ok := parentSymbolAccount[hs.syncer.Account.FirstExternalAddress[hs.parentSymbol]]
	if !ok {
		msg := fmt.Sprintf("%s does not exist", hs.symbol)
		log.Warn().Msg(msg)
		return errors.New(msg)
	}

	httpClient := newHTTPClient()
	resp, err := httpClient.Get(fmt.Sprintf("%s/%s", hs.RPC, hTokenAccount.Address))
	if err != nil {
		log.Error().Err(err).Msgf("failed to get %s balance", hs.symbol)
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("failed to read response body")
		return err
	}

	balance, err := strconv.ParseInt(strings.TrimSuffix(string(body), "\n"), 10, 64)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse response body")
		return err
	}

	hs.syncer.ExtBalance[hTokenAccount.Address] = big.NewInt(balance)

	return nil
}

// Update updates accounts in cache as and when external balances
// external chains are updated.
func (hs *HTokenSyncer) Update() {
	if hs.syncer.Account.EBalances[hs.parentSymbol] == nil {
		log.Warn().Msgf("No %s account available, skip", hs.parentSymbol)
		return
	}
	if hs.syncer.Account.EBalances[hs.symbol] == nil {
		hs.syncer.Account.EBalances[hs.symbol] = make(map[string]statedb.EBalance)
		hs.syncer.Account.EBalances[hs.symbol][hs.syncer.Account.FirstExternalAddress[hs.parentSymbol]] = statedb.EBalance{Address: hs.syncer.Account.FirstExternalAddress[hs.parentSymbol]}
	}

	// HToken account is first ETH account of user.
	parentSymbolAccount := hs.syncer.Account.EBalances[hs.symbol][hs.syncer.Account.FirstExternalAddress[hs.parentSymbol]]
	hs.syncer.update(parentSymbolAccount.Address)
}
