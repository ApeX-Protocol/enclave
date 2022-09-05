package enclave

import (
	"crypto/ecdsa"
	"errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"go.uber.org/zap"
	"math/big"
)

type ToSignUser struct {
	Nonce          []byte         `json:"nonce"`
	Amount         *big.Int       `json:"amount"` // amount 1e18
	User           common.Address `json:"user"`
	ExpireAt       *big.Int       `json:"expireAt"` // second
	RewardContract common.Address `json:"rewardContract"`
	UseFor         *big.Int       `json:"useFor"`
}

type SignedUser struct {
	Nonce     []byte         `json:"nonce"`
	User      common.Address `json:"user"`
	Amount    *big.Int       `json:"amount"`   // amount 1e18
	ExpireAt  *big.Int       `json:"expireAt"` // second
	Signature string         `json:"signature"`
}

func SignRewards(logger *zap.Logger, toSignUser *ToSignUser, privateKey *ecdsa.PrivateKey) (*SignedUser, error) {
	uint256Ty, _ := abi.NewType("uint256", "uint256", nil)
	bytesTy, _ := abi.NewType("bytes", "bytes", nil)
	addressTy, _ := abi.NewType("address", "address", nil)

	// support solidity abi.encode(user, useFor, amount, expireAt, nonce, address(this))
	arguments := abi.Arguments{
		{
			Type: addressTy,
		},
		{
			Type: uint256Ty,
		},
		{
			Type: uint256Ty,
		},
		{
			Type: uint256Ty,
		},
		{
			Type: bytesTy,
		},
		{
			Type: addressTy,
		},
	}

	abiEncode, _ := arguments.Pack(toSignUser.User, toSignUser.UseFor, toSignUser.Amount, toSignUser.ExpireAt, toSignUser.Nonce, toSignUser.RewardContract)

	hash := crypto.Keccak256Hash(abiEncode)
	ethHash := crypto.Keccak256([]byte("\x19Ethereum Signed Message:\n32"), hash.Bytes())

	signature, err := crypto.Sign(ethHash, privateKey)

	if err != nil {
		logger.Error("Error signing",
			zap.Error(err),
		)

		return nil, err
	}

	if signature[64] == 0 || signature[64] == 1 {
		signature[64] += 27
	}

	return &SignedUser{
		Nonce:     toSignUser.Nonce,
		User:      toSignUser.User,
		ExpireAt:  toSignUser.ExpireAt,
		Signature: hexutil.Encode(signature),
		Amount:    toSignUser.Amount,
	}, nil
}

func MnemonicToPK(logger *zap.Logger, mnemonic string) (*ecdsa.PrivateKey, error) {
	privateKey, err := crypto.HexToECDSA(mnemonic)

	if err != nil {
		logger.Error("Error converting hex to private key",
			zap.Error(err),
		)

		return nil, errors.New("Error converting hex to private key")
	}

	return privateKey, nil
}
