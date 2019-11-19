package sync

import (
	"encoding/json"
	"errors"
	"math/big"
	"strconv"

	"github.com/herdius/herdius-core/p2p/log"
	"github.com/herdius/herdius-core/symbol"
)

// BNBSyncer syncs all BNB external accounts
type BNBSyncer struct {
	RPC    string
	syncer *ExternalSyncer
}

func newBNBSyncer() *BNBSyncer {
	t := &BNBSyncer{}
	t.syncer = newExternalSyncer(symbol.BNB)

	return t
}

type Balance struct {
	Free   string
	Symbol string
}

type bnbBalance struct {
	Balances []Balance
}

// GetExtBalance syncs bnb account.
func (bs *BNBSyncer) GetExtBalance() error {
	bsAccount, ok := bs.syncer.Account.EBalances[bs.syncer.assetSymbol]
	if !ok {
		return errors.New("BTC account does not exists")
	}

	httpClient := newHTTPClient()

	for _, ba := range bsAccount {
		bs.syncer.addressError[ba.Address] = true
		resp, err := httpClient.Get(bs.RPC + "/" + ba.Address)
		if err != nil {
			log.Error().Err(err).Msg("Error connecting lite coin network")
			continue
		}

		balanceResp := &bnbBalance{}
		if err := json.NewDecoder(resp.Body).Decode(balanceResp); err != nil {
			log.Error().Err(err).Msg("failed to decode response body")
		}

		for _, b := range balanceResp.Balances {
			if b.Symbol == "BNB" {
				balance, err := strconv.ParseFloat(b.Free, 64)
				if err == nil {
					bs.syncer.addressError[ba.Address] = false
					bs.syncer.ExtBalance[ba.Address] = big.NewInt(int64(balance * mutez))
					bs.syncer.BlockHeight[ba.Address] = big.NewInt(0)
				}
			}
		}
		_ = resp.Body.Close()
	}
	return nil
}

// Update updates accounts in cache as and when external balances
// external chains are updated.
func (bs *BNBSyncer) Update() {
	for _, bnbAccount := range bs.syncer.Account.EBalances[bs.syncer.assetSymbol] {
		if bs.syncer.addressError[bnbAccount.Address] {
			//log.Warn().Msgf("BNB account info is not available at this moment, skip sync: %s", bnbAccount.Address)
			continue
		}
		bs.syncer.update(bnbAccount.Address)
	}
}
