package tx

import (
	b64 "encoding/base64"
	"encoding/json"

	"github.com/herdius/herdius-core/crypto/merkle"
	"github.com/herdius/herdius-core/crypto/secp256k1"
	"github.com/herdius/herdius-core/hbi/protobuf"
	cmn "github.com/herdius/herdius-core/libs/common"
)

// Tx is an arbitrary byte array.
type Tx []byte

// Txs is a slice of Tx.
// TODO make Txs of type `[]Tx`
type Txs [][]byte

// Proof represents a Merkle proof of the presence of a transaction in the Merkle tree.
type Proof struct {
	RootHash cmn.HexBytes
	Data     Tx
	Proof    merkle.SimpleProof
}

func VerifySignature(tx *protobuf.Tx, pubKey secp256k1.PubKeySecp256k1) (bool, error) {
	asset := &protobuf.Asset{
		Category:              tx.Asset.Category,
		Symbol:                tx.Asset.Symbol,
		Network:               tx.Asset.Network,
		Value:                 tx.Asset.Value,
		Fee:                   tx.Asset.Fee,
		Nonce:                 tx.Asset.Nonce,
		ExternalSenderAddress: tx.Asset.ExternalSenderAddress,
		LockedAmount:          tx.Asset.LockedAmount,
		RedeemedAmount:        tx.Asset.RedeemedAmount,
	}
	verifiableTx := protobuf.Tx{
		SenderAddress:   tx.SenderAddress,
		SenderPubkey:    tx.SenderPubkey,
		RecieverAddress: tx.RecieverAddress,
		Asset:           asset,
		Message:         tx.Message,
		Type:            tx.Type,
		Data:            tx.Data,
		ExternalAddress: tx.ExternalAddress,
	}

	txbBeforeSign, err := json.Marshal(verifiableTx)
	if err != nil {
		return false, err
	}

	decodedSig, err := b64.StdEncoding.DecodeString(tx.Sign)

	if err != nil {
		return false, err
	}

	return pubKey.VerifyBytes(txbBeforeSign, decodedSig), nil
}
