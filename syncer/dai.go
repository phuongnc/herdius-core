package sync

import (
	"context"
	"errors"

	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/herdius/herdius-core/p2p/log"
	"github.com/herdius/herdius-core/syncer/contract"
)

// DaiSyncer syncs all DAI external accounts
type DaiSyncer struct {
	RPC                  string
	syncer               *ExternalSyncer
	TokenContractAddress string
}

func newDAISyncer() *DaiSyncer {
	d := &DaiSyncer{}
	d.syncer = newExternalSyncer("DAI")

	return d
}

// GetExtBalance ...
func (ds *DaiSyncer) GetExtBalance() error {
	// If DAI account exists
	daiAccount, ok := ds.syncer.Account.EBalances[ds.syncer.assetSymbol]
	if !ok {
		return errors.New("DAI account does not exists")
	}

	for _, ta := range daiAccount {
		client, err := ethclient.Dial(ds.RPC)
		if err != nil {
			log.Error().Msgf("Error connecting DAI RPC: %v", err)
			ds.syncer.addressError[ta.Address] = true
			continue
		}
		tokenAddress := common.HexToAddress(ds.TokenContractAddress)
		address := common.HexToAddress(ta.Address)

		// Get latest block number
		latestBlockNumber, err := ds.getLatestBlockNumber(client)
		if err != nil {
			log.Error().Err(err).Msg("Error getting Latest Block Number from RPC")
			ds.syncer.addressError[ta.Address] = true
			continue
		}

		//Get nonce
		nonce, err := ds.getNonce(client, address, latestBlockNumber)
		if err != nil {
			log.Error().Err(err).Msg("Error getting TOKEN Account nonce from RPC")
			ds.syncer.addressError[ta.Address] = true
			continue
		}

		instance, err := contract.NewToken(tokenAddress, client)
		if err != nil {
			log.Error().Err(err).Msg("Error getting TOKEN instance")
			ds.syncer.addressError[ta.Address] = true
			continue
		}
		ctx, cancel := context.WithTimeout(context.Background(), rpcTimeout)
		defer cancel()

		balance, err := instance.BalanceOf(&bind.CallOpts{BlockNumber: latestBlockNumber, Context: ctx}, address)
		if err != nil {
			log.Error().Err(err).Msgf("Error getting DAI Balance from RPC: %v", err)
			ds.syncer.addressError[ta.Address] = true
			continue
		}

		ds.syncer.ExtBalance[ta.Address] = balance
		ds.syncer.BlockHeight[ta.Address] = latestBlockNumber
		ds.syncer.Nonce[ta.Address] = nonce
		ds.syncer.addressError[ta.Address] = false
	}

	return nil
}

// Update updates accounts in cache as and when external balances
// external chains are updated.
func (ds *DaiSyncer) Update() {
	for _, daiAccount := range ds.syncer.Account.EBalances[ds.syncer.assetSymbol] {
		if ds.syncer.addressError[daiAccount.Address] {
			log.Warn().Msgf("Dai account info is not available at this moment, skip sync: %s", daiAccount.Address)
			continue
		}
		ds.syncer.update(daiAccount.Address)
	}
}

func (ds *DaiSyncer) getLatestBlockNumber(client *ethclient.Client) (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), rpcTimeout)
	defer cancel()
	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, err
	}
	return header.Number, nil
}

func (ds *DaiSyncer) getNonce(client *ethclient.Client, account common.Address, block *big.Int) (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), rpcTimeout)
	defer cancel()
	nonce, err := client.NonceAt(ctx, account, block)
	if err != nil {
		return 0, err
	}
	return nonce, nil
}
