package sync

import (
	"github.com/herdius/herdius-core/p2p/log"
)

// LTCSyncer syncs all XTZ external accounts
type LTCSyncer struct {
	RPC    string
	syncer *ExternalSyncer
}

func newLTCSyncer() *LTCSyncer {
	t := &LTCSyncer{}
	t.syncer = newExternalSyncer("LTC")

	return t
}

// GetExtBalance syncs lite coin account.
func (ls *LTCSyncer) GetExtBalance() error {
	// TODO implements later.
	return nil
}

// Update updates accounts in cache as and when external balances
// external chains are updated.
func (ls *LTCSyncer) Update() {
	for _, ltcAccount := range ls.syncer.Account.EBalances[ls.syncer.assetSymbol] {
		if ls.syncer.addressError[ltcAccount.Address] {
			log.Warn().Msgf("LTC account info is not available at this moment, skip sync: %s", ltcAccount.Address)
			continue
		}
		ls.syncer.update(ltcAccount.Address)
	}
}
