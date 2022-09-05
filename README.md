## enclave

## How to use
```go
package main

func main() {
	wallet := "0x72864219334298fF2F97183b8Dd9444B4CD8625C"
	user := &ToSignUser{
		Nonce:          []byte(uuid.NewString()),
		Amount:         big.NewInt(0.002 * 1e18),
		User:           common.HexToAddress(wallet),
		ExpireAt:       big.NewInt(1662958800),
		RewardContract: common.HexToAddress("0xf3105df48f330D1AfDb45d41F4946D4bA77f6153"),
		UseFor:         big.NewInt(int64(1)),
	}

	logger, _ := zap.NewDevelopment(zap.Development())
	pk, err := MnemonicToPK(logger, "your mnemonic")
	if err != nil {
		logger.Error("Error converting mnemonic to private key",
			zap.Error(err),
		)
	}
	r, err := SignRewards(logger, user, pk)
	if err != nil {
		logger.Error("Error signing rewards",
			zap.Error(err),
		)
	}
	log.Println("-----")
	log.Println(r.User)
	log.Println(r.Amount)
	u, err := uuid.Parse(string(r.Nonce))
	if err != nil {
		log.Fatal(err)
	}
	log.Println(u.String())
	log.Println(hexutil.Encode([]byte(u.String())))
	log.Println(r.ExpireAt)
	log.Println(r.Signature)
}
```