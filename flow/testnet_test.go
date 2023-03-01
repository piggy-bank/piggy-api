package flow

import (
	"context"
	"fmt"
	"log"
	"testing"

	"cloud.google.com/go/logging"
	"github.com/manubidegain/piggy-api/utils"
	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/stretchr/testify/require"
)

// This test tests the contracts associated with the service account, such as
// NonFungibleToken, MetaData, and Ticker

func TestNetTest(t *testing.T) {

	client, _, e := newTestNetSetup(t)
	fmt.Print(e)

	//adminAccountKey, ourSigner := k.NewWithSigner()
	adminAddress := flow.HexToAddress("36e55122ece3464c")

	privateKey, err := crypto.DecodePrivateKeyHex(crypto.ECDSA_P256, "ae02806e92a9a4581ef5ea28166053afddc187329934692dd24b4a8352bb0d24")
	if err != nil {
		msg := "error decoding privateKey" + err.Error()
		panic(msg)
	}
	acc, err := client.GetAccount(context.Background(), adminAddress)
	if err != nil {
		msg := "error getting service account" + err.Error()
		panic(msg)
	}
	accountKey := acc.Keys[0]

	fmt.Print(privateKey)
	fmt.Print(accountKey)

	signer, _ := crypto.NewInMemorySigner(privateKey, accountKey.HashAlgo)
	fmt.Print(signer)

	ctx := context.Background()

	anotherPrivateKey, err := crypto.DecodePrivateKeyHex(crypto.ECDSA_P256, "af59d90a2692bf78d1f13086c5e7f586a7faf091f427f9416f60bd65b3bc0b44")
	if err != nil {
		msg := "error decoding privateKey" + err.Error()
		panic(msg)
	}

	anotherAddress := flow.HexToAddress("52dbb96b337d3765")
	anotherAcc, err := client.GetAccount(context.Background(), anotherAddress)
	if err != nil {
		msg := "error getting service account" + err.Error()
		panic(msg)
	}

	e.AnotherAccountAddress = anotherAddress.String()

	anotherAccountKey := anotherAcc.Keys[0]

	anotherSigner, _ := crypto.NewInMemorySigner(anotherPrivateKey, anotherAccountKey.HashAlgo)
	fmt.Printf(anotherSigner.PrivateKey.String())

	projectID := "zinc-involution-379214"

	// Creates a client.
	loggingClient, err := logging.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer loggingClient.Close()

	// Sets the name of the log to write to.
	logName := "my-log"

	logger := loggingClient.Logger(logName).StandardLogger(logging.Info)
	fmt.Printf(logger.Prefix())

	/*t.Run("Creating Piggy...", func(t *testing.T) {
		metadata := make(map[string]string)
		metadata["test"] = "yes"
		tx, err := CreatePiggy(client, e, adminAddress, metadata, accountKey, logger)
		require.NoError(t, err)

		err = tx.SignEnvelope(adminAddress, accountKey.Index, signer)
		utils.LogAndPanicError(logger, err)
		err = client.SendTransaction(ctx, *tx)
		utils.LogAndPanicError(logger, err)

		createMinterTxResp := utils.WaitForSeal(ctx, client, tx.ID())
		utils.LogAndPanicError(logger, createMinterTxResp.Error)
	})

	t.Run("Get number of piggies..", func(t *testing.T) {
		nextPiggyID := cadence.NewUInt32(12)
		result := executeScriptAndCheckInTestnet(t, client, GenerateGetNextPiggyID(e), nil, ctx)
		assertEqual(t, nextPiggyID, result)

	}) */

	t.Run("Minting Donation...", func(t *testing.T) {
		piggyID := 12
		donationComment := "First donation on testnet"
		tx, err := MintDonation(client, e, adminAddress, piggyID, donationComment, adminAddress, accountKey, logger)
		require.NoError(t, err)

		err = tx.SignEnvelope(adminAddress, accountKey.Index, signer)
		utils.LogAndPanicError(logger, err)
		err = client.SendTransaction(ctx, *tx)
		utils.LogAndPanicError(logger, err)

		createMinterTxResp := utils.WaitForSeal(ctx, client, tx.ID())
		utils.LogAndPanicError(logger, createMinterTxResp.Error)
	})

	t.Run("Get number of donations..", func(t *testing.T) {
		totalSupply := cadence.NewUInt64(2)
		result := executeScriptAndCheckInTestnet(t, client, GenerateGetTotalSupply(e), nil, ctx)
		assertEqual(t, totalSupply, result)

	})

}
