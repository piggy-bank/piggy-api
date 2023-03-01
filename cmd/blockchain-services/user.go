package blockchainservices

// Create account and store account on datastore

// GetKey to sign

import (
	"crypto/rand"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/manubidegain/piggy-api/cmd/api/configuration"
	flowUtils "github.com/manubidegain/piggy-api/flow"
	"github.com/manubidegain/piggy-api/utils"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/flow-go-sdk/examples"
	"github.com/onflow/flow-go-sdk/templates"
)

func GetAccount(ctx *gin.Context, profile string, config *configuration.FlowConfig) {
	client, err := utils.ConnectToFlow(profile, config)
	if err != nil {
		msg := "Cannot connect to flow" + err.Error()
		panic(msg)
	}
	address := ctx.Param("address")
	account := utils.GetAccount(ctx, client, address)
	if account == nil {
		ctx.IndentedJSON(http.StatusNotFound, "Account not found")
		return
	}
	utils.CloseConnection(client)
	ctx.IndentedJSON(http.StatusOK, account)
}

func TemporaryGetValue(ctx *gin.Context, profile string, projectConfig *configuration.ProjectConfig) {
	address := ctx.Param("address")
	dataStoreClient, err := utils.SetupDataStoreClient("dev", projectConfig)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, nil)
		return
	}
	readed, err := utils.GetValue(ctx, dataStoreClient, address, "Account", profile, projectConfig)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, nil)
		return
	}
	ctx.IndentedJSON(http.StatusOK, readed)
}

func CreateAccount(ctx *gin.Context, profile string, config *configuration.FlowConfig, log *log.Logger, projectConfig *configuration.ProjectConfig) (string, error) {
	client, err := utils.ConnectToFlow(profile, config)
	env := flowUtils.NewEnv(profile)
	if err != nil {
		log.Println("Cannot connect to flow with profile " + profile)
		msg := "Cannot connect to flow" + err.Error()
		panic(msg)
	}

	serviceAcctAddr, serviceAcctKey, serviceSigner := utils.GetServiceAccount(client, config, profile)

	//Handle more options
	myPrivateKey := createRandomPrivateKey()
	newAcctKey := flow.NewAccountKey().
		FromPrivateKey(myPrivateKey).
		SetHashAlgo(crypto.SHA3_256).
		SetWeight(flow.AccountKeyWeightThreshold)

	//privateKey, err := crypto.DecodePrivateKeyHex(crypto.ECDSA_P256, myPrivateKey.String())

	anotherSigner, err := crypto.NewInMemorySigner(myPrivateKey, newAcctKey.HashAlgo)
	if err != nil {
		return "", err
	}

	referenceBlockID := utils.GetReferenceBlockId(client, log)
	createAccountTx, err := templates.CreateAccount([]*flow.AccountKey{newAcctKey}, nil, serviceAcctAddr)
	createAccountTx.SetProposalKey(
		serviceAcctAddr,
		serviceAcctKey.Index,
		serviceAcctKey.SequenceNumber,
	)
	createAccountTx.SetReferenceBlockID(referenceBlockID)
	createAccountTx.SetPayer(serviceAcctAddr)

	if err != nil {
		msg := "cannot generate the transaction: " + err.Error()
		log.Println(msg)
		panic(msg)
		return "", err
	}

	err = createAccountTx.SignEnvelope(serviceAcctAddr, serviceAcctKey.Index, serviceSigner)
	if err != nil {
		msg := "cannot sign envelope : " + err.Error()
		log.Println(msg)
		panic(msg)
		return "", err
	}

	// Send the transaction to the network
	err = client.SendTransaction(ctx, *createAccountTx)
	if err != nil {
		msg := "error sending transaction" + err.Error()
		log.Println(msg)
		panic(msg)
		return "", err
	}

	accountCreationTxRes := examples.WaitForSeal(ctx, client, createAccountTx.ID())

	var newAddress flow.Address

	for _, event := range accountCreationTxRes.Events {
		if event.Type == flow.EventAccountCreated {
			accountCreatedEvent := flow.AccountCreatedEvent(event)
			newAddress = accountCreatedEvent.Address()
		}
	}

	newServiceAccountKey := serviceAcctKey
	newServiceAccountKey.SequenceNumber++

	tx := flowUtils.FundAccount(client, env, newAddress, 0.001, serviceAcctAddr, newServiceAccountKey, log)

	err = tx.SignEnvelope(serviceAcctAddr, serviceAcctKey.Index, serviceSigner)
	err = utils.HandleAndLogError(log, err)
	if err != nil {
		return "", err
	}
	err = client.SendTransaction(ctx, *tx)
	err = utils.HandleAndLogError(log, err)
	if err != nil {
		return "", err
	}

	createMinterTxResp := utils.WaitForSeal(ctx, client, tx.ID())
	err = utils.HandleAndLogError(log, createMinterTxResp.Error)
	if err != nil {
		return "", err
	}

	// Setup acc
	tx2 := flowUtils.SetupAccount(client, env, newAddress, newAcctKey, log)

	err = tx2.SignEnvelope(newAddress, newAcctKey.Index, anotherSigner)
	err = utils.HandleAndLogError(log, err)
	if err != nil {
		return "", err
	}
	err = client.SendTransaction(ctx, *tx2)
	err = utils.HandleAndLogError(log, err)
	if err != nil {
		return "", err
	}

	createMinterTxResp = utils.WaitForSeal(ctx, client, tx2.ID())
	err = utils.HandleAndLogError(log, createMinterTxResp.Error)
	if err != nil {
		return "", err
	}

	// Fund acc

	utils.CloseConnection(client)

	dataStoreClient, err := utils.SetupDataStoreClient("dev", projectConfig)
	if err != nil {
		log.Println(err)
		return "", err
	}
	entry, err := utils.CreateNewEntry(newAddress.Hex(), newAcctKey.PublicKey.String(), myPrivateKey.String(), projectConfig)
	if err != nil {
		log.Println(err)
		return "", err
	}

	key, err := utils.UploadValue(ctx, entry, "Account", dataStoreClient)
	if err != nil {
		log.Println(err)
		return "", err
	}
	dataStoreClient.Close()

	return key, nil
}

func createRandomPrivateKey() crypto.PrivateKey {
	seed := make([]byte, crypto.MinSeedLength)
	_, err := rand.Read(seed)
	if err != nil {
		msg := "error reading seed" + err.Error()
		panic(msg)
	}

	privateKey, err := crypto.GeneratePrivateKey(crypto.ECDSA_P256, seed)
	if err != nil {
		msg := "error generating private key" + err.Error()
		panic(msg)
	}

	return privateKey
}
