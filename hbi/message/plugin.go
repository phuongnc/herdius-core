package message

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/herdius/herdius-core/accounts/account"
	"github.com/herdius/herdius-core/blockchain"
	"github.com/herdius/herdius-core/checker"
	"github.com/herdius/herdius-core/storage/mempool"

	"github.com/herdius/herdius-core/hbi/protobuf"
	protoplugin "github.com/herdius/herdius-core/hbi/protobuf"
	"github.com/herdius/herdius-core/libs/common"
	cmn "github.com/herdius/herdius-core/libs/common"
	plog "github.com/herdius/herdius-core/p2p/log"
	"github.com/herdius/herdius-core/p2p/network"
)

// TxType ...
//go:generate stringer -type=TxType
type TxType int

const (
	Update TxType = iota
	Lock
	Redeem
)

// BlockMessagePlugin will receive all Block specific messages.
type BlockMessagePlugin struct {
	*network.Plugin
}

// AccountMessagePlugin will receive all Account specific messages.
type AccountMessagePlugin struct {
	*network.Plugin
}

// TransactionMessagePlugin will receive all Transaction specific messages.
type TransactionMessagePlugin struct {
	*network.Plugin
	env string
}

// NewTransactionMessagePlugin returns a TransactionMessagePllugin with given env.
func NewTransactionMessagePlugin(env string) *TransactionMessagePlugin {
	return &TransactionMessagePlugin{env: env}
}

// Receive handles account specific messages
func (state *AccountMessagePlugin) Receive(ctx *network.PluginContext) error {
	switch msg := ctx.Message().(type) {
	case *protoplugin.AccountRequest:
		return getAccount(msg.Address, ctx)

	case *protoplugin.AccountResponse:
		plog.Info().Msgf("Account Response: %v", msg)
	}
	return nil
}

// Receive handles block specific received messages
func (state *BlockMessagePlugin) Receive(ctx *network.PluginContext) error {
	switch msg := ctx.Message().(type) {

	case *protoplugin.LastBlockRequest:
		return getLastBlock(ctx)

	case *protoplugin.BlockHeightRequest:
		return getBlock(msg.BlockHeight, ctx)

	case *protoplugin.BlockResponse:
		plog.Info().Msgf("Block Response: %v", msg)
	}
	return nil
}

// Receive to handle transaction requests
func (state *TransactionMessagePlugin) Receive(ctx *network.PluginContext) error {
	switch msg := ctx.Message().(type) {

	case *protoplugin.TxsByAssetAndAddressRequest:
		address := msg.GetAddress()
		asset := msg.GetAsset()
		accSrv := account.NewAccountService()
		acc, err := accSrv.GetAccountByAddress(address)
		if err != nil {
			if err := ctx.Reply(network.WithSignMessage(context.Background(), true), &protoplugin.TxsResponse{}); err != nil {
				return fmt.Errorf(fmt.Sprintf("Failed to reply to client: %v", err))
			}
			return errors.New("Couldn't find the acc due to: " + err.Error())
		}
		if acc != nil && strings.EqualFold(acc.Address, address) {
			return getTxsByAssetAndAccount(asset, address, ctx)
		}

		if err := ctx.Reply(network.WithSignMessage(context.Background(), true), &protoplugin.TxsResponse{}); err != nil {
			return fmt.Errorf(fmt.Sprintf("Failed to reply to client: %v", err))
		}
		return nil

	case *protoplugin.TxsByAddressRequest:
		address := msg.GetAddress()
		accSrv := account.NewAccountService()
		acc, err := accSrv.GetAccountByAddress(address)
		if err != nil {
			if err := ctx.Reply(network.WithSignMessage(context.Background(), true), &protoplugin.TxsResponse{}); err != nil {
				return fmt.Errorf(fmt.Sprintf("Failed to reply to client: %v", err))
			}
			return errors.New("Couldn't find the acc due to: " + err.Error())
		}

		if acc != nil && strings.EqualFold(acc.Address, address) {
			return getTxs(address, ctx)
		}

		if err := ctx.Reply(network.WithSignMessage(context.Background(), true), &protoplugin.TxsResponse{}); err != nil {
			return fmt.Errorf(fmt.Sprintf("Failed to reply to client: %v", err))
		}
		return nil

	case *protoplugin.TxsByBlockHeightRequest:
		return getTxsByblockHeight(msg.GetBlockHeight(), ctx)

	case *protoplugin.TxDetailRequest:
		return getTx(msg.GetTxId(), ctx)

	case *protoplugin.TxRequest:
		tx := msg.GetTx()
		assetChecker := checker.New(state.env)
		if !assetChecker.Check(tx.Asset.Symbol) {
			err := fmt.Errorf("%s network is not available.", tx.Asset.Symbol)
			if err := ctx.Reply(network.WithSignMessage(context.Background(), true), &protoplugin.TxResponse{
				TxId: "", Status: "failed", Queued: 0, Pending: 0,
				Message: err.Error(),
			}); err != nil {
				return fmt.Errorf(fmt.Sprintf("Failed to reply to client: %v", err))
			}
			return err
		}
		accSrv := account.NewAccountService()
		accSrv.SetReceiverAddress(tx.RecieverAddress)
		accSrv.SetAssetSymbol(tx.Asset.Symbol)
		accSrv.SetExtAddress(tx.Asset.ExternalSenderAddress)
		accSrv.SetTxValue(tx.Asset.Value)
		accSrv.SetTxLockedAmount(tx.Asset.LockedAmount)
		accSrv.SetTxRedeemAmount(tx.Asset.RedeemedAmount)
		accSrv.SetTxType(tx.Type)
		acc, err := accSrv.GetAccountByAddress(msg.Tx.GetSenderAddress())
		if err != nil {
			if err := ctx.Reply(network.WithSignMessage(context.Background(), true), &protoplugin.TxResponse{
				TxId: "", Status: "failed", Queued: 0, Pending: 0,
				Message: "Couldn't find the acc due to : " + err.Error(),
			}); err != nil {
				return fmt.Errorf(fmt.Sprintf("Failed to reply to client: %v", err))
			}
			return errors.New("Couldn't find the acc due to: " + err.Error())
		}

		// Check Tx.Nonce > acc.Nonce
		if acc != nil {
			if !accSrv.VerifyAccountNonce(acc, tx.GetAsset().Nonce) {

				txNonce := strconv.FormatUint(msg.Tx.GetAsset().Nonce, 10)
				accountNonce := strconv.FormatUint(acc.Nonce, 10)
				failedVerificationMsg := "Transaction nonce " + txNonce +
					" should be greater than acc nonce " + accountNonce

				if err := ctx.Reply(network.WithSignMessage(context.Background(), true), &protoplugin.TxResponse{
					TxId: "", Status: "failed", Queued: 0, Pending: 0,
					Message: failedVerificationMsg,
				}); err != nil {
					return fmt.Errorf(fmt.Sprintf("Failed to reply to client :%v", err))
				}
				return errors.New(failedVerificationMsg)
			}
		}
		accSrv.SetAccount(acc)
		// Check if tx is of type acc update
		// and verify external address exists
		update := Update.String()
		if strings.EqualFold(tx.Type, update) {
			if accSrv.AccountExternalAddressExist() {
				failedVerificationMsg := "External acc existed: " + tx.Asset.ExternalSenderAddress
				if err := ctx.Reply(network.WithSignMessage(context.Background(), true), &protoplugin.TxResponse{
					TxId: "", Status: "failed", Queued: 0, Pending: 0,
					Message: failedVerificationMsg,
				}); err != nil {
					return fmt.Errorf(fmt.Sprintf("Failed to reply to client :%v", err))
				}
				return errors.New(failedVerificationMsg)
			}
			if accSrv.AccountEBalancePerAssetReachLimit() {
				failedVerificationMsg := "Account reached number of addresses limit"
				if err := ctx.Reply(network.WithSignMessage(context.Background(), true), &protoplugin.TxResponse{
					TxId: "", Status: "failed", Queued: 0, Pending: 0,
					Message: failedVerificationMsg,
				}); err != nil {
					return fmt.Errorf(fmt.Sprintf("Failed to reply to client :%v", err))
				}
				return errors.New(failedVerificationMsg)
			}
			return postAccountUpdateTx(tx, ctx, accSrv)
		}

		// Check if tx is of type lock or redeem
		// verify if external acc address doesn't exists
		// verify if receiver address is herdius zero address
		lock := Lock.String()
		redeem := Redeem.String()
		isLockTx := strings.EqualFold(tx.Type, lock)
		isRedeemTx := strings.EqualFold(tx.Type, redeem)
		if isLockTx || isRedeemTx {
			if !accSrv.AccountExternalAddressExist() {
				failedVerificationMsg := "external address does not exist: " + tx.Asset.ExternalSenderAddress
				if err := ctx.Reply(network.WithSignMessage(context.Background(), true), &protoplugin.TxResponse{
					TxId: "", Status: "failed", Queued: 0, Pending: 0,
					Message: failedVerificationMsg,
				}); err != nil {
					return fmt.Errorf(fmt.Sprintf("Failed to reply to client :%v", err))
				}
				return errors.New(failedVerificationMsg)
			}
			if !accSrv.IsHerdiusZeroAddress() {
				failedVerificationMsg := "incorrect Herdius zero address: " + tx.RecieverAddress
				if err := ctx.Reply(network.WithSignMessage(context.Background(), true), &protoplugin.TxResponse{
					TxId: "", Status: "failed", Queued: 0, Pending: 0,
					Message: failedVerificationMsg,
				}); err != nil {
					return fmt.Errorf(fmt.Sprintf("Failed to reply to client :%v", err))
				}
				return errors.New(failedVerificationMsg)
			}
		}

		if isLockTx && !accSrv.VerifyLockedAmount() {
			failedVerificationMsg := "account does not have enough locked amount"
			if err := ctx.Reply(network.WithSignMessage(context.Background(), true), &protoplugin.TxResponse{
				TxId: "", Status: "failed", Queued: 0, Pending: 0,
				Message: failedVerificationMsg,
			}); err != nil {
				return fmt.Errorf(fmt.Sprintf("Failed to reply to client :%v", err))
			}
			return errors.New(failedVerificationMsg)
		}

		if isRedeemTx && !accSrv.VerifyRedeemAmount() {
			failedVerificationMsg := "account does not have enough amount to redeem"
			if err := ctx.Reply(network.WithSignMessage(context.Background(), true), &protoplugin.TxResponse{
				TxId: "", Status: "failed", Queued: 0, Pending: 0,
				Message: failedVerificationMsg,
			}); err != nil {
				return fmt.Errorf(fmt.Sprintf("Failed to reply to client :%v", err))
			}
			return errors.New(failedVerificationMsg)
		}

		// Check if asset has enough balance
		// acc.Balance > Tx.Value
		if strings.EqualFold(tx.Type, update) && !accSrv.VerifyAccountBalance() {
			if err := ctx.Reply(network.WithSignMessage(context.Background(), true), &protoplugin.TxResponse{
				TxId: "", Status: "failed", Queued: 0, Pending: 0,
				Message: "Not enough balance: " + string(msg.Tx.GetAsset().Value),
			}); err != nil {
				return fmt.Errorf(fmt.Sprintf("Failed to reply to client :%v", err))
			}
			return errors.New("Not enough balance: " + strconv.FormatUint(uint64(msg.Tx.GetAsset().Value), 10))
		}

		// Add Tx to Mempool
		mp := mempool.GetMemPool()
		txbz, err := cdc.MarshalJSON(tx)
		if err != nil {
			if err := ctx.Reply(network.WithSignMessage(context.Background(), true), &protoplugin.TxResponse{
				TxId: "", Status: "failed", Queued: 0, Pending: 0,
				Message: "Incorrect Transaction format : " + msg.Tx.GetSenderAddress(),
			}); err != nil {
				return fmt.Errorf(fmt.Sprintf("Failed to reply to client :%v", err))
			}
			return errors.New("Failed to Masshal Tx: " + msg.Tx.GetSenderAddress())
		}
		log.Println("Add tx to mempool")
		pending, queue := mp.AddTx(tx, accSrv)

		plog.Info().Msgf("Remaining mempool pending, queue: %+v %+v", pending, queue)

		// Create the Transaction ID
		txID := cmn.CreateTxID(txbz)
		plog.Info().Msgf("Tx ID : %v", txID)

		// Send Tx ID to client who sent the TX
		if err = ctx.Reply(network.WithSignMessage(context.Background(), true), &protoplugin.TxResponse{
			TxId: txID, Status: "success", Queued: int64(queue), Pending: int64(pending),
		}); err != nil {
			return fmt.Errorf("Failed to reply to client :%v", err)
		}
	case *protoplugin.TxUpdateRequest:
		log.Println("Update request received")
		id := msg.GetTxId()
		newTx := msg.GetTx()
		log.Println("Processing request to update Tx, ID:", id)
		updatedID, updatedTx, err := putTxUpdateRequest(id, newTx)
		if err != nil {
			if errRep := ctx.Reply(network.WithSignMessage(context.Background(), true),
				&protoplugin.TxUpdateResponse{
					Error:  err.Error(),
					Status: false,
				}); errRep != nil {
				return fmt.Errorf("could not reply to API client, transaction not updated: %v", errRep)
			}

			return fmt.Errorf("could not update request: %v", err)
		}
		if errRep := ctx.Reply(network.WithSignMessage(context.Background(), true), &protoplugin.TxUpdateResponse{
			Status: true,
			TxId:   updatedID,
			Tx:     updatedTx,
		}); errRep != nil {
			return fmt.Errorf("could not reply to API client, but transaction was updated: %v", errRep)
		}

		return nil
	case *protoplugin.TxDeleteRequest:
		mp := mempool.GetMemPool()
		succ := mp.DeleteTx(msg.TxId)
		if succ {
			return ctx.Reply(network.WithSignMessage(context.Background(), true), &protoplugin.TxUpdateResponse{
				TxId:   msg.TxId,
				Status: succ,
			})
		}
		return ctx.Reply(network.WithSignMessage(context.Background(), true), &protoplugin.TxUpdateResponse{
			Status: succ,
			Error:  fmt.Sprintf("Unable to find Tx (id: %v) in memory pool", msg.TxId),
		})
	case *protoplugin.TxLockedRequest:
		return getLockedTxsByBlockNumber(ctx, msg.BlockNumber)
	case *protoplugin.TxRedeemRequest:
		return getRedeemTxsByBlockNumber(ctx, msg.BlockNumber)
	}
	return nil
}

func postAccountUpdateTx(tx *protoplugin.Tx, ctx *network.PluginContext, as account.ServiceI) error {
	// Add Tx to Mempool
	mp := mempool.GetMemPool()
	txbz, err := cdc.MarshalJSON(tx)
	if err != nil {
		if err := ctx.Reply(network.WithSignMessage(context.Background(), true),
			&protoplugin.TxResponse{
				TxId: "", Status: "failed", Queued: 0, Pending: 0,
				Message: "Transaction format incorrect : " + err.Error(),
			}); err != nil {
			return fmt.Errorf(fmt.Sprintf("Failed to reply to client :%v", err))
		}
		return errors.New("Failed to Marshal Tx: " + err.Error())
	}

	pending, queue := mp.AddTx(tx, as)

	plog.Info().Msgf("Remaining mempool pending, queue: %+v %+v", pending, queue)

	// Create the Transaction ID
	txID := cmn.CreateTxID(txbz)
	plog.Info().Msgf("Tx ID : %v", txID)

	// Send Tx ID to client who sent the TX
	if err := ctx.Reply(network.WithSignMessage(context.Background(), true),
		&protoplugin.TxResponse{
			TxId: txID, Status: "success", Queued: 0, Pending: 0,
		}); err != nil {
		return fmt.Errorf(fmt.Sprintf("Failed to reply to client :%v", err))
	}
	return nil
}

func getBlock(height int64, ctx *network.PluginContext) error {
	blockchainSvc := &blockchain.Service{}
	block, err := blockchainSvc.GetBlockByHeight(height)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("Failed to retrieve the Block: :%v", err))
	}

	supervisorAdd := ctx.Client().ID.Address
	if block.Header != nil {
		timestamp := &protoplugin.Timestamp{
			Nanos:   block.GetHeader().GetTime().GetNanos(),
			Seconds: block.GetHeader().GetTime().GetSeconds(),
		}

		totalTxs := block.GetHeader().TotalTxs

		blockRes := protoplugin.BlockResponse{
			BlockHeight:       block.GetHeader().GetHeight(),
			TotalTxs:          totalTxs,
			Time:              timestamp,
			SupervisorAddress: supervisorAdd,
		}

		plog.Info().Msgf("Block Response at processor: %v", blockRes)

		if err := ctx.Reply(network.WithSignMessage(context.Background(), true), &blockRes); err != nil {
			return fmt.Errorf(fmt.Sprintf("Failed to reply to client :%v", err))
		}
		return nil
	}

	if err := ctx.Reply(network.WithSignMessage(context.Background(), true), &protoplugin.BlockResponse{}); err != nil {
		return fmt.Errorf(fmt.Sprintf("Failed to reply to client :%v", err))
	}
	return nil
}

func getAccount(address string, ctx *network.PluginContext) error {
	accountSvc := &account.Service{}
	acc, err := accountSvc.GetAccountByAddress(address)
	if err != nil {
		plog.Error().Msgf("Failed to retrieve the Account: %v", err)
	}

	if acc == nil {
		if err := ctx.Reply(network.WithSignMessage(context.Background(), true), &protoplugin.AccountResponse{}); err != nil {
			return fmt.Errorf(fmt.Sprintf("Failed to reply to client :%v", err))
		}
	}

	if acc != nil {
		eBalances := make(map[string]*protoplugin.EBalanceAsset)

		for asset, assetAccount := range acc.EBalances {
			eBalances[asset] = &protoplugin.EBalanceAsset{}
			eBalances[asset].Asset = make(map[string]*protobuf.EBalance)
			for _, eb := range assetAccount.Asset {
				eBalanceRes := &protoplugin.EBalance{
					Address:         eb.Address,
					Balance:         eb.Balance,
					LastBlockHeight: eb.LastBlockHeight,
					Nonce:           eb.Nonce,
				}
				eBalances[asset].Asset[eb.Address] = eBalanceRes
			}
		}
		accountResp := protoplugin.AccountResponse{
			Address:              address,
			Nonce:                acc.Nonce,
			Balance:              acc.Balance,
			StorageRoot:          acc.StorageRoot,
			PublicKey:            acc.PublicKey,
			EBalances:            eBalances,
			Erc20Address:         acc.Erc20Address,
			ExternalNonce:        acc.ExternalNonce,
			LastBlockHeight:      acc.LastBlockHeight,
			FirstExternalAddress: acc.FirstExternalAddress,
		}
		if err := ctx.Reply(network.WithSignMessage(context.Background(), true), &accountResp); err != nil {
			return fmt.Errorf(fmt.Sprintf("Failed to reply to client :%v", err))
		}
	}
	return nil
}

func getTx(id string, ctx *network.PluginContext) error {
	txSvc := &blockchain.TxService{}
	txDetailRes, err := txSvc.GetTx(id)
	if err != nil {
		if err := ctx.Reply(network.WithSignMessage(context.Background(), true),
			&protoplugin.TxDetailResponse{}); err != nil {
			return fmt.Errorf(fmt.Sprintf("Failed to reply to client: %v", err))
		}
		return errors.New("Failed due to: " + err.Error())
	}
	log.Println("txDetailRes: ", txDetailRes)
	if err := ctx.Reply(network.WithSignMessage(context.Background(), true), txDetailRes); err != nil {
		return fmt.Errorf(fmt.Sprintf("Failed to reply to client: %v", err))
	}
	return nil
}

func getTxs(address string, ctx *network.PluginContext) error {
	txSvc := &blockchain.TxService{}
	txs, err := txSvc.GetTxs(address)
	if err != nil {
		if err := ctx.Reply(network.WithSignMessage(context.Background(), true), &protoplugin.TxDetailResponse{}); err != nil {
			return fmt.Errorf(fmt.Sprintf("Failed to reply to client: %v", err))
		}
		return errors.New("Failed due to: " + err.Error())
	}

	if err := ctx.Reply(network.WithSignMessage(context.Background(), true), txs); err != nil {
		return fmt.Errorf(fmt.Sprintf("Failed to reply to client: %v", err))
	}
	return nil
}

func getTxsByAssetAndAccount(asset, address string, ctx *network.PluginContext) error {
	txSvc := &blockchain.TxService{}
	txs, err := txSvc.GetTxsByAssetAndAddress(asset, address)
	if err != nil {
		if err := ctx.Reply(network.WithSignMessage(context.Background(), true), &protoplugin.TxDetailResponse{}); err != nil {
			return fmt.Errorf(fmt.Sprintf("Failed to reply to client: %v", err))
		}
		return errors.New("Failed due to: " + err.Error())
	}

	if err := ctx.Reply(network.WithSignMessage(context.Background(), true), txs); err != nil {
		return fmt.Errorf(fmt.Sprintf("Failed to reply to client: %v", err))
	}
	return nil
}

// putTxUpdateRequest updates the Tx with the input string. After updating, calculates new Tx ID
func putTxUpdateRequest(id string, newTx *protoplugin.Tx) (string, *protoplugin.Tx, error) {
	mp := mempool.GetMemPool()
	i, origTx, err := mp.GetTx(id)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get original transaction details (id: %v), err: %v", id, err)
	}
	if origTx == nil {
		return "", nil, fmt.Errorf("requested Tx (id: %v) does not exist in memory pool; it may have been flushed from the memory pool into a block", id)
	}
	updatedTx, err := mp.UpdateTx(i, newTx)
	if err != nil {
		return "", nil, fmt.Errorf("failed to update Tx in MemPool with new values: %v", err)
	}
	updatedBz, err := cdc.MarshalJSON(updatedTx)
	if err != nil {
		return "", nil, fmt.Errorf("could not marshal updated transaction back into memory pool: %v", err)
	}
	newID := common.CreateTxID(updatedBz)
	return newID, updatedTx, nil
}

func getLockedTxsByBlockNumber(ctx *network.PluginContext, blockNumber int64) error {
	txSvc := &blockchain.TxService{}
	txs, err := txSvc.GetLockedTxsByBlockNumber(blockNumber)
	if err != nil {
		return fmt.Errorf("error getting lost transactions: %v", err)
	}
	if err := ctx.Reply(network.WithSignMessage(context.Background(), true), txs); err != nil {
		return fmt.Errorf("error replying to apiClient: %v", err)
	}
	return nil
}
func getRedeemTxsByBlockNumber(ctx *network.PluginContext, blockNumber int64) error {
	txSvc := &blockchain.TxService{}
	txs, err := txSvc.GetRedeemTxsByBlockNumber(blockNumber)
	if err != nil {
		return fmt.Errorf("error getting lost transactions: %v", err)
	}
	if err := ctx.Reply(network.WithSignMessage(context.Background(), true), txs); err != nil {
		return fmt.Errorf("error replying to apiClient: %v", err)
	}
	return nil
}

func getTxsByblockHeight(height int64, ctx *network.PluginContext) error {
	txSvc := &blockchain.TxService{}
	txs, err := txSvc.GetTxsByHeight(height)
	if err != nil {
		if err := ctx.Reply(network.WithSignMessage(context.Background(), true), &protoplugin.TxDetailResponse{}); err != nil {
			return fmt.Errorf(fmt.Sprintf("Failed to reply to client: %v", err))
		}
		return errors.New("Failed due to: " + err.Error())
	}

	if err := ctx.Reply(network.WithSignMessage(context.Background(), true), txs); err != nil {
		return fmt.Errorf(fmt.Sprintf("Failed to reply to client: %v", err))
	}
	return nil
}

func getLastBlock(ctx *network.PluginContext) error {
	blockchainSvc := &blockchain.Service{}
	block := blockchainSvc.GetLastBlock()

	supervisorAdd := ctx.Client().ID.Address
	if block.Header != nil {
		timestamp := &protoplugin.Timestamp{
			Nanos:   block.GetHeader().GetTime().GetNanos(),
			Seconds: block.GetHeader().GetTime().GetSeconds(),
		}

		totalTxs := block.GetHeader().TotalTxs

		blockRes := protoplugin.BlockResponse{
			BlockHeight:       block.GetHeader().GetHeight(),
			TotalTxs:          totalTxs,
			Time:              timestamp,
			SupervisorAddress: supervisorAdd,
		}

		plog.Info().Msgf("Block Response at processor: %v", blockRes)

		if err := ctx.Reply(network.WithSignMessage(context.Background(), true), &blockRes); err != nil {
			return fmt.Errorf(fmt.Sprintf("Failed to reply to client :%v", err))
		}
		return nil
	}

	if err := ctx.Reply(network.WithSignMessage(context.Background(), true), &protoplugin.BlockResponse{}); err != nil {
		return fmt.Errorf(fmt.Sprintf("Failed to reply to client :%v", err))
	}
	return nil
}
