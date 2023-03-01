package utils

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/gin-gonic/gin"
	"github.com/manubidegain/piggy-api/cmd/api/configuration"
	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/onflow/cadence/runtime/sema"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/access"
	"github.com/onflow/flow-go-sdk/access/grpc"
	"github.com/onflow/flow-go-sdk/crypto"
)

const (
	Development              = "dev"
	Production               = "prod"
	Test                     = "test"
	configurationPackagePath = "configfiles"
)

func CalculateProfile() string {
	scope := os.Getenv("SCOPE")

	if scope == "" {
		msg := "can not start application without 'SCOPE' environment variable"
		panic(msg)
	}

	scope = strings.ToLower(scope)
	switch scope {
	case Production:
		return Production
	case Test:
		return Test
	default:
		return Development
	}
}

func BuildConfig(profile string) *configuration.Config {
	path := configurationPackagePath + "/properties-" + profile + ".yml"
	return configuration.GetConfig(path)

}

func WaitForSeal(ctx context.Context, c access.Client, id flow.Identifier) *flow.TransactionResult {
	result, err := c.GetTransactionResult(ctx, id)
	if err != nil {
		msg := "Cannot get transaction result" + err.Error()
		panic(msg)
	}
	fmt.Printf("Waiting for transaction %s to be sealed...\n", id)

	for result.Status != flow.TransactionStatusSealed {
		time.Sleep(time.Second)
		fmt.Print(".")
		result, err = c.GetTransactionResult(ctx, id)
		if err != nil {
			msg := "Cannot get transaction result" + err.Error()
			panic(msg)
		}
	}
	fmt.Printf("Transaction %s sealed\n", id)
	return result
}

func PrintTransaction(ctx context.Context, c access.Client, id flow.Identifier, log *log.Logger) *flow.Transaction {
	result, err := c.GetTransaction(ctx, id)
	LogAndPanicError(log, err)

	return result
}

func GetReferenceBlockId(flowClient access.Client, log *log.Logger) flow.Identifier {
	block, err := flowClient.GetLatestBlock(context.Background(), true)
	LogAndPanicError(log, err)

	return block.ID
}

func GetAccount(ctx *gin.Context, client access.Client, address string) *flow.Account {
	addr := flow.HexToAddress(address)
	account, err := client.GetAccount(ctx, addr)
	if err != nil {
		msg := "error getting account" + err.Error()
		panic(msg)
	}
	return account
}

func ConnectToFlow(profile string, flowConfig *configuration.FlowConfig) (access.Client, error) {
	if profile == "dev" {
		flow, err := grpc.NewClient(grpc.EmulatorHost)
		if err != nil {
			panic("failed to establish connection with the Emulator")
		}
		return flow, nil
	}

	if profile == "test" {

		flow, err := grpc.NewClient(grpc.TestnetHost)
		if err != nil {
			panic("failed to establish connection with the Testnet network")
		}
		return flow, nil
	}

	if profile == "prod" {
		flow, err := grpc.NewClient(grpc.MainnetHost)
		if err != nil {
			panic("failed to establish connection with the Mainnet network")
		}
		return flow, nil
	}

	return nil, errors.New("profile is neccessary for connection")

}

func CloseConnection(client access.Client) {
	err := client.Close()
	if err != nil {
		panic(err)
	}
}

func HandleAndLogError(log *log.Logger, err error) error {
	if err != nil {
		log.Println(err.Error())
		fmt.Println("err:", err.Error())
		return err
	}
	return nil
}

func LogAndPanicError(log *log.Logger, err error) {
	if err != nil {
		log.Println(err.Error())
		fmt.Println("err:", err.Error())
		panic(err)
	}
}

func HandleCreateEventErr(log *log.Logger, err error) error {
	if err != nil {
		if !strings.Contains(err.Error(), "force assignment to non-nil resource-typed value") {
			log.Println(err.Error())
			fmt.Println("err:", err.Error())
			return err
		}
	}
	return nil
}

func HandleCreateTicketTypeErr(log *log.Logger, err error) error {
	if err != nil {
		if !strings.Contains(err.Error(), "force assignment to non-nil resource-typed value") {
			log.Println(err.Error())
			fmt.Println("err:", err.Error())
			return err
		}
	}
	return nil
}

func HandleAssignTicketTypeErr(log *log.Logger, err error) error {
	if err != nil {
		if !strings.Contains(err.Error(), "The TicketType has already beed added to the event.") {
			log.Println(err.Error())
			fmt.Println("err:", err.Error())
			return err
		}
	}
	return nil
}

func HandleMintTicketTypeErr(log *log.Logger, err error) (bool, error) {
	if err != nil {
		if !(strings.Contains(err.Error(), "Cannot borrow Event: Event doesn't exist") || (strings.Contains(err.Error(), "This TicketType doesn't exist"))) {
			log.Println(err.Error())
			fmt.Println("err:", err.Error())
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func PullMsgs(w io.Writer, subID string, projectConfig *configuration.ProjectConfig) error {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectConfig.ProjectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err.Error())
	}
	defer client.Close()

	//msg := make(chan *pubsub.Message, 1)
	sub := client.Subscription(subID)
	cctx, cancel := context.WithCancel(ctx)
	err = sub.Receive(cctx, func(ctx context.Context, m *pubsub.Message) {
		fmt.Fprintf(w, "Got message: %q\n", string(m.Data))
		data, err := base64.StdEncoding.DecodeString(string(m.Data))
		if err != nil {
			log.Fatalf("Base64: %v", err)
			cancel()
		}

		dataMap := make(map[string]interface{})
		if err := json.Unmarshal(data, &dataMap); err != nil {
			log.Fatalf("Json: %v", err)
			cancel()
		}
		printAllData(dataMap)
		m.Ack()
	})
	if err != nil {
		return err
	}
	/*for {
	select {
	case res := <-msg:
		fmt.Fprintf(w, "Got message: %q\n", string(res.Data))
		data, err := base64.StdEncoding.DecodeString(string(res.Data))
		if err != nil {
			log.Fatalf("Base64: %v", err)
			cancel()
			return err
		}

		dataMap := make(map[string]interface{})
		if err := json.Unmarshal(data, &dataMap); err != nil {
			log.Fatalf("Json: %v", err)
			cancel()
			return err
		}
		printAllData(dataMap)
		res.Ack()

	case <-time.After(3 * time.Second):
		fmt.Println("timeout")
		cancel()
	}
	*/
	return nil
}

func printAllData(dataMap map[string]interface{}) {
	for key, value := range dataMap {
		fmt.Printf("Key : %s And Value : %v ", key, value)
	}
}

func GetServiceAccount(flowClient access.Client, config *configuration.FlowConfig, profile string) (flow.Address, *flow.AccountKey, crypto.Signer) {
	// Handle this with secrets for PROD
	account := getServiceAccount(config, profile)
	privateKey, err := crypto.DecodePrivateKeyHex(crypto.ECDSA_P256, account.Key)
	if err != nil {
		msg := "error decoding privateKey" + err.Error()
		panic(msg)
	}

	addr := flow.HexToAddress(account.Address)
	acc, err := flowClient.GetAccount(context.Background(), addr)
	if err != nil {
		msg := "error getting service account" + err.Error()
		panic(msg)
	}
	accountKey := acc.Keys[0]
	signer, _ := crypto.NewInMemorySigner(privateKey, accountKey.HashAlgo)
	return addr, accountKey, signer
}

// RandomPrivateKey returns a randomly generated ECDSA P-256 private key.
func RandomPrivateKey(log *log.Logger) crypto.PrivateKey {
	seed := make([]byte, crypto.MinSeedLength)
	_, err := rand.Read(seed)
	LogAndPanicError(log, err)

	privateKey, err := crypto.GeneratePrivateKey(crypto.ECDSA_P256, seed)
	LogAndPanicError(log, err)

	return privateKey
}

func getServiceAccount(config *configuration.FlowConfig, profile string) configuration.FlowServiceAccount {
	if profile == "dev" {
		return config.Accounts.Emulator
	}

	if profile == "test" {

		return config.Accounts.Testnet
	}

	if profile == "prod" {
		return config.Accounts.Mainnet
	}

	return configuration.FlowServiceAccount{}
}

// Needs to get this from our blockchain-api
var mintTokensToAccountTemplate = `
import FungibleToken from 0x%s
import FlowToken from 0x%s

transaction(recipient: Address, amount: UFix64) {
	let tokenAdmin: &FlowToken.Administrator
	let tokenReceiver: &{FungibleToken.Receiver}

	prepare(signer: AuthAccount) {
		self.tokenAdmin = signer
			.borrow<&FlowToken.Administrator>(from: /storage/flowTokenAdmin)
			?? panic("Signer is not the token admin")

		self.tokenReceiver = getAccount(recipient)
			.getCapability(/public/flowTokenReceiver)
			.borrow<&{FungibleToken.Receiver}>()
			?? panic("Unable to borrow receiver reference")
	}

	execute {
		let minter <- self.tokenAdmin.createNewMinter(allowedAmount: amount)
		let mintedVault <- minter.mintTokens(amount: amount)

		self.tokenReceiver.deposit(from: <-mintedVault)

		destroy minter
	}
}
`

func SetupAccount(profile string, flowClient access.Client, config *configuration.FlowConfig, address flow.Address, amount float64) {
	//Needs to setup and fund account
}

func fundAccount(flowClient access.Client, config *configuration.FlowConfig, address flow.Address, amount float64) {
	//Transfer the funds from main account to new account
}

// Make the same for testnet and mainnet
func fundAccountInEmulator(flowClient access.Client, config *configuration.FlowConfig, address flow.Address, amount float64, profile string, log *log.Logger) {
	serviceAcctAddr, serviceAcctKey, serviceSigner := GetServiceAccount(flowClient, config, profile)

	referenceBlockID := GetReferenceBlockId(flowClient, log)

	fungibleTokenAddress := flow.HexToAddress(config.Contracts["FungibleToken"])
	flowTokenAddress := flow.HexToAddress(config.Contracts["FlowToken"])

	recipient := cadence.NewAddress(address)
	uintAmount := uint64(amount * sema.Fix64Factor)
	cadenceAmount := cadence.UFix64(uintAmount)

	fundAccountTx :=
		flow.NewTransaction().
			SetScript([]byte(fmt.Sprintf(mintTokensToAccountTemplate, fungibleTokenAddress, flowTokenAddress))).
			AddAuthorizer(serviceAcctAddr).
			AddRawArgument(jsoncdc.MustEncode(recipient)).
			AddRawArgument(jsoncdc.MustEncode(cadenceAmount)).
			SetProposalKey(serviceAcctAddr, serviceAcctKey.Index, serviceAcctKey.SequenceNumber).
			SetReferenceBlockID(referenceBlockID).
			SetPayer(serviceAcctAddr)

	err := fundAccountTx.SignEnvelope(serviceAcctAddr, serviceAcctKey.Index, serviceSigner)
	LogAndPanicError(log, err)

	ctx := context.Background()
	err = flowClient.SendTransaction(ctx, *fundAccountTx)
	LogAndPanicError(log, err)

	result := WaitForSeal(ctx, flowClient, fundAccountTx.ID())
	LogAndPanicError(log, result.Error)
}

func ReadFile(path string, log *log.Logger) string {
	contents, err := ioutil.ReadFile(path)
	LogAndPanicError(log, err)

	return string(contents)
}

func MakeInternalApiCall(method string, url string, body []byte, ctx *gin.Context, bearer string) ([]byte, error) {
	var req *http.Request
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if method != http.MethodGet {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(body))
	}
	if err != nil {
		return nil, err
	}
	//bearer := ctx.Request.Header.Get("Authorization")
	req.Header.Add("Authorization", bearer)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// change in 200 range
	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return nil, fmt.Errorf("unable to make api call %s , returns statuscode : %v with body %s", url, resp.StatusCode, string(bodyBytes[:]))
	}
	return bodyBytes, nil
}

func MakeInternalApiCallNoAuth(method string, url string, body []byte, ctx *gin.Context) ([]byte, error) {
	var req *http.Request
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if method != http.MethodGet {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(body))
	}
	if err != nil {
		return nil, err
	}
	//bearer := ctx.Request.Header.Get("Authorization")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// change in 200 range
	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return nil, fmt.Errorf("unable to make api call %s , returns statuscode : %v with body %s", url, resp.StatusCode, string(bodyBytes[:]))
	}
	return bodyBytes, nil
}
