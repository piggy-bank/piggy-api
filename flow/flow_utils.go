package flow

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"testing"

	"github.com/manubidegain/piggy-api/utils"
	"github.com/onflow/cadence/runtime/sema"
	emulator "github.com/onflow/flow-emulator"
	"github.com/onflow/flow-emulator/types"

	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	ft_templates "github.com/onflow/flow-ft/lib/go/templates"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/access"
	"github.com/onflow/flow-go-sdk/access/grpc"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/flow-go-sdk/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	flow_crypto "github.com/onflow/flow-go/crypto"
)

/***********************************************
*
*    flow-core-contracts/lib/go/test/test.go
*
*    Provides common testing utilities for automated testing using the Flow emulator
*    such as setting up the emulator, submitting transactions and scripts,
*    constructing cadence values, creating accounts, and minting tokens
*
*    To use, import the `onflow/flow-core-contracts/lib/go/test` package
*    and call any of these functions, such as:
*
*    test.newTestSetup(t)
*
************************************************/

const (
	mainnetFTAddress                = "f233dcee88fe0abe"
	mainnetFlowTokenAddress         = "1654653399040a61"
	mainnetFlowFeesAddress          = "f919ee77447b7497"
	mainnetNonFungibleTokenAddress  = "1d7e57aa55817448"
	mainnetMetadataViewsAddress     = "1d7e57aa55817448"
	mainnetNFTStoreFrontAddress     = "e5a8b7f23e8b5483"
	mainnetPiggyAddress             = "09e8665388e90671"
	mainnetAnotherAccountAddress    = "e5a8b7f23e8b5486"
	testnetFTAddress                = "9a0766d93b6608b7"
	testnetFlowTokenAddress         = "7e60df042a9c0868"
	testnetFlowFeesAddress          = "912d5440f7e3769e"
	testnetNonFungibleTokenAddress  = "631e88ae7f1d7c20"
	testnetMetadataViewsAddress     = "631e88ae7f1d7c20"
	testnetNFTStoreFrontAddress     = "e5a8b7f23e8b5483"
	testnetPiggyAddress             = "36e55122ece3464c"
	testnetAnotherAccountAddress    = "e5a8b7f23e8b5486"
	emulatorFTAddress               = "f233dcee88fe0abe"
	emulatorFlowTokenAddress        = "1654653399040a61"
	emulatorFlowFeesAddress         = "f919ee77447b7497"
	emulatorNonFungibleTokenAddress = "1d7e57aa55817448"
	emulatorMetadataViewsAddress    = "1d7e57aa55817448"
	emulatorNFTStoreFrontAddress    = "e5a8b7f23e8b5483"
	emulatorPiggyAddress            = "36e55122ece3464c"
	emulatorAnotherAccountAddress   = "e5a8b7f23e8b5486"
	getExecutionEffortWeighs        = "FlowServiceAccount/scripts/get_execution_effort_weights.cdc"
	setExecutionEffortWeighs        = "FlowServiceAccount/set_execution_effort_weights.cdc"
)

type Environment struct {
	Network                    string
	NonFungibleTokenAddress    string
	PiggyAddress               string
	MetadataViewsAddress       string
	TicketAdminReceiverAddress string
	AnotherAccountAddress      string
	NFTStoreFrontAddress       string
	FungibleTokenAddress       string
	FlowTokenAddress           string
	IDTableAddress             string
	LockedTokensAddress        string
	StakingProxyAddress        string
	QuorumCertificateAddress   string
	DkgAddress                 string
	EpochAddress               string
	StorageFeesAddress         string
	FlowFeesAddress            string
	ServiceAccountAddress      string
}

func NewEnv(profile string) Environment {
	if profile == "prod" {
		return Environment{
			FungibleTokenAddress:    mainnetFTAddress,
			FlowTokenAddress:        mainnetFlowTokenAddress,
			NonFungibleTokenAddress: mainnetNonFungibleTokenAddress,
			PiggyAddress:            mainnetPiggyAddress,
			MetadataViewsAddress:    mainnetMetadataViewsAddress,
			NFTStoreFrontAddress:    mainnetNFTStoreFrontAddress,
			AnotherAccountAddress:   mainnetAnotherAccountAddress,
		}

	} else if profile == "test" {
		return Environment{
			FungibleTokenAddress:    testnetFTAddress,
			FlowTokenAddress:        testnetFlowTokenAddress,
			NonFungibleTokenAddress: testnetNonFungibleTokenAddress,
			PiggyAddress:            testnetPiggyAddress,
			MetadataViewsAddress:    testnetMetadataViewsAddress,
			NFTStoreFrontAddress:    testnetNFTStoreFrontAddress,
			AnotherAccountAddress:   testnetAnotherAccountAddress,
		}
	} else {
		return Environment{
			FungibleTokenAddress:    emulatorFTAddress,
			FlowTokenAddress:        emulatorFlowTokenAddress,
			NonFungibleTokenAddress: emulatorNonFungibleTokenAddress,
			PiggyAddress:            emulatorPiggyAddress,
			MetadataViewsAddress:    emulatorMetadataViewsAddress,
			NFTStoreFrontAddress:    emulatorNFTStoreFrontAddress,
			AnotherAccountAddress:   emulatorAnotherAccountAddress,
		}
	}
}

// Sets up testing and emulator objects and initialize the emulator default addresses
func newTestSetup(t *testing.T) (*emulator.Blockchain, *test.AccountKeys, Environment) {
	// Set for parallel processing
	//t.Parallel()

	// Create a new emulator instance
	b := newBlockchain()

	// Create a new account key generator object to generate keys
	// for test accounts
	accountKeys := test.AccountKeyGenerator()

	// Setup the env variable that stores import addresses for various contracts
	env := Environment{
		FungibleTokenAddress:    emulatorFTAddress,
		FlowTokenAddress:        emulatorFlowTokenAddress,
		NonFungibleTokenAddress: emulatorNonFungibleTokenAddress,
		PiggyAddress:            emulatorPiggyAddress,
		MetadataViewsAddress:    emulatorMetadataViewsAddress,
		NFTStoreFrontAddress:    emulatorNFTStoreFrontAddress,
		AnotherAccountAddress:   emulatorAnotherAccountAddress,
	}

	return b, accountKeys, env
}

func newTestNetSetup(t *testing.T) (access.Client, *test.AccountKeys, Environment) {
	// Set for parallel processing
	//t.Parallel()

	// Create a new emulator instance
	//b := newBlockchain()

	flowClient, err := grpc.NewClient(grpc.TestnetHost)

	if err != nil {
		panic(err)
	}

	// Create a new account key generator object to generate keys
	// for test accounts
	accountKeys := test.AccountKeyGenerator()

	// Setup the env variable that stores import addresses for various contracts
	env := Environment{
		FungibleTokenAddress:    emulatorFTAddress,
		FlowTokenAddress:        emulatorFlowTokenAddress,
		NonFungibleTokenAddress: emulatorNonFungibleTokenAddress,
		PiggyAddress:            emulatorPiggyAddress,
		MetadataViewsAddress:    emulatorMetadataViewsAddress,
		NFTStoreFrontAddress:    emulatorNFTStoreFrontAddress,
		AnotherAccountAddress:   emulatorAnotherAccountAddress,
	}

	return flowClient, accountKeys, env
}

// newBlockchain returns an emulator blockchain for testing.
func newBlockchain(opts ...emulator.Option) *emulator.Blockchain {
	b, err := emulator.NewBlockchain(
		append(
			[]emulator.Option{
				// No storage limit
				emulator.WithStorageLimitEnabled(false),
			},
			opts...,
		)...,
	)
	if err != nil {
		panic(err)
	}
	return b
}

// Create a new, empty account for testing
// and return the address, public keys, and signer objects
func newAccountWithAddress(b *emulator.Blockchain, accountKeys *test.AccountKeys) (flow.Address, *flow.AccountKey, crypto.Signer) {
	newAccountKey, newSigner := accountKeys.NewWithSigner()
	newAddress, err := b.CreateAccount([]*flow.AccountKey{newAccountKey}, nil)
	if err != nil {
		panic(err)
	}

	return newAddress, newAccountKey, newSigner
}

// Create a basic transaction template with the provided transaction script
// Sets the service account as the proposer and payer
// Uses the max gas limit
// authorizer address is the authorizer for the transaction (transaction has access to their AuthAccount object)
// Whoever authorizes here also needs to sign the envelope and payload when submitting the transaction
// returns the tx object so arguments can be added to it and it can be signed
func createTxWithTemplateAndAuthorizer(
	b *emulator.Blockchain,
	script []byte,
	authorizerAddress flow.Address,
) *flow.Transaction {

	tx := flow.NewTransaction().
		SetScript(script).
		SetGasLimit(9999).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(authorizerAddress)

	return tx
}

// signAndSubmit signs a transaction with an array of signers and adds their signatures to the transaction
// before submitting it to the emulator.
//
// If the private keys do not match up with the addresses, the transaction will not succeed.
//
// The shouldRevert parameter indicates whether the transaction should fail or not.
//
// This function asserts the correct result and commits the block if it passed.
func signAndSubmit(
	t *testing.T,
	b *emulator.Blockchain,
	tx *flow.Transaction,
	signerAddresses []flow.Address,
	signers []crypto.Signer,
	shouldRevert bool,
) *types.TransactionResult {
	// sign transaction with each signer
	for i := len(signerAddresses) - 1; i >= 0; i-- {
		signerAddress := signerAddresses[i]
		signer := signers[i]

		if i == 0 {
			err := tx.SignEnvelope(signerAddress, 0, signer)
			assert.NoError(t, err)
		} else {
			err := tx.SignPayload(signerAddress, 0, signer)
			assert.NoError(t, err)
		}
	}

	return Submit(t, b, tx, shouldRevert)
}

// Submit submits a transaction and checks if it fails or not, based on shouldRevert specification
func Submit(
	t *testing.T,
	b *emulator.Blockchain,
	tx *flow.Transaction,
	shouldRevert bool,
) *types.TransactionResult {
	// submit the signed transaction
	err := b.AddTransaction(*tx)
	require.NoError(t, err)

	// use the emulator to execute it
	result, err := b.ExecuteNextTransaction()
	require.NoError(t, err)

	// Check the status
	if shouldRevert {
		assert.True(t, result.Reverted())
	} else {
		if !assert.True(t, result.Succeeded()) {
			t.Log(result.Error.Error())
		}
	}

	_, err = b.CommitBlock()
	assert.NoError(t, err)

	return result
}

// executeScriptAndCheck executes a script and checks to make sure that it succeeded.
func executeScriptAndCheck(t *testing.T, b *emulator.Blockchain, script []byte, arguments [][]byte) cadence.Value {
	result, err := b.ExecuteScript(script, arguments)
	require.NoError(t, err)
	if !assert.True(t, result.Succeeded()) {
		t.Log(result.Error.Error())
	}

	return result.Value
}

// executeScriptAndCheck executes a script and checks to make sure that it succeeded.
func executeScriptAndCheckInTestnet(t *testing.T, client access.Client, script []byte, arguments []cadence.Value, ctx context.Context) cadence.Value {
	clientResult, err := client.ExecuteScriptAtLatestBlock(ctx, script, arguments)
	require.NoError(t, err)

	return clientResult
}

// Reads a file from the specified path
func readFile(path string) []byte {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return contents
}

// CadenceUFix64 returns a UFix64 value from a string representation
func CadenceUFix64(value string) cadence.Value {
	newValue, err := cadence.NewUFix64(value)

	if err != nil {
		panic(err)
	}

	return newValue
}

// CadenceUFix64 returns a UFix64 value from a string representation
func CadenceAddress(value flow.Address) cadence.Value {
	newValue := cadence.BytesToAddress(value.Bytes())
	return newValue
}

// CadenceString returns a string value from a string representation
func CadenceMapStringInt(value map[string]int) cadence.Value {
	keyValuesPairs := []cadence.KeyValuePair{}
	for key, element := range value {
		cadenceKey, err := cadence.NewString(key)
		if err != nil {
			panic(err)
		}
		cadenceElement := cadence.NewInt(element)
		keyValuePair := cadence.KeyValuePair{Key: cadenceKey, Value: cadenceElement}
		keyValuesPairs = append(keyValuesPairs, keyValuePair)

	}

	newValue := cadence.NewDictionary(keyValuesPairs)
	return newValue
}

// CadenceString returns a string value from a string representation
func CadenceMapStringString(value map[string]string) cadence.Value {
	keyValuesPairs := []cadence.KeyValuePair{}
	for key, element := range value {
		cadenceKey, err := cadence.NewString(key)
		if err != nil {
			panic(err)
		}
		cadenceElement, err := cadence.NewString(element)
		if err != nil {
			panic(err)
		}
		keyValuePair := cadence.KeyValuePair{Key: cadenceKey, Value: cadenceElement}
		keyValuesPairs = append(keyValuesPairs, keyValuePair)

	}

	newValue := cadence.NewDictionary(keyValuesPairs)
	return newValue
}

// CadenceUFix64 returns a UFix64 value from a string representation
func CadenceInt(value int) cadence.Value {
	newValue := cadence.NewInt(value)
	return newValue
}

// CadenceUFix64 returns a UFix64 value from a string representation
func CadenceArrayUFix64(vals []string) cadence.Array {
	ourArray := []cadence.Value{}
	for _, value := range vals {
		ufixvalue, err := cadence.NewUFix64(value)
		if err != nil {
			panic(err)
		}
		ourArray = append(ourArray, ufixvalue)
	}

	return cadence.NewArray(ourArray)
}

// CadenceUFix64 returns a UFix64 value from a string representation
func CadenceArrayUInt32(vals []int) cadence.Array {
	ourArray := []cadence.Value{}
	for _, value := range vals {
		ourArray = append(ourArray, cadence.NewUInt32(uint32(value)))
	}

	return cadence.NewArray(ourArray)
}

func CadenceArrayUInt64(vals []int) cadence.Array {
	ourArray := []cadence.Value{}
	for _, value := range vals {
		ourArray = append(ourArray, cadence.NewUInt64(uint64(value)))
	}

	return cadence.NewArray(ourArray)
}

// CadenceString returns a string value from a string representation
func CadenceString(value string) cadence.Value {
	newValue, err := cadence.NewString(value)

	if err != nil {
		panic(err)
	}

	return newValue
}

// Converts a byte array to a Cadence array of UInt8
func bytesToCadenceArray(b []byte) cadence.Array {
	values := make([]cadence.Value, len(b))

	for i, v := range b {
		values[i] = cadence.NewUInt8(v)
	}

	return cadence.NewArray(values)
}

// assertEqual asserts that two objects are equal.
//
//	assertEqual(t, 123, 123)
//
// Pointer variable equality is determined based on the equality of the
// referenced values (as opposed to the memory addresses). Function equality
// cannot be determined and will always fail.
func assertEqual(t *testing.T, expected, actual interface{}) bool {

	if assert.ObjectsAreEqual(expected, actual) {
		return true
	}

	message := fmt.Sprintf(
		"Not equal: \nexpected: %s\nactual  : %s",
		expected,
		actual,
	)

	return assert.Fail(t, message)
}

// Mints the specified amount of FLOW tokens for the specified account address
// Using the mint tokens template from the onflow/flow-ft repo
// signed by the service account
func mintTokensForAccount(t *testing.T, b *emulator.Blockchain, recipient flow.Address, amount string) {

	// Create a new mint FLOW transaction template authorized by the service account
	tx := createTxWithTemplateAndAuthorizer(b,
		ft_templates.GenerateMintTokensScript(flow.HexToAddress(emulatorFTAddress), flow.HexToAddress(emulatorFlowTokenAddress), "FlowToken"),
		b.ServiceKey().Address)

	// Add the recipient and amount as arguments
	_ = tx.AddArgument(cadence.NewAddress(recipient))
	_ = tx.AddArgument(CadenceUFix64(amount))

	cryptoSigner, _ := b.ServiceKey().Signer()

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address},
		[]crypto.Signer{cryptoSigner},
		false,
	)
}

// Creates multiple accounts and mints 1B tokens for each one
// Returns the addresses, keys, and signers for each account in matching arrays
func registerAndMintManyAccounts(
	t *testing.T,
	b *emulator.Blockchain,
	accountKeys *test.AccountKeys,
	numAccounts int) ([]flow.Address, []*flow.AccountKey, []crypto.Signer) {

	// make new addresses, keys, and signers
	var userAddresses = make([]flow.Address, numAccounts)
	var userPublicKeys = make([]*flow.AccountKey, numAccounts)
	var userSigners = make([]crypto.Signer, numAccounts)

	// Create each new account and mint 1B tokens for it
	for i := 0; i < numAccounts; i++ {
		userAddresses[i], userPublicKeys[i], userSigners[i] = newAccountWithAddress(b, accountKeys)
		mintTokensForAccount(t, b, userAddresses[i], "1000000000.0")
	}

	return userAddresses, userPublicKeys, userSigners
}

// Generates a new private/public key pair
func generateKeys(t *testing.T, algorithmName flow_crypto.SigningAlgorithm) (crypto.PrivateKey, crypto.PublicKey) {
	seedMinLength := 48
	seed := make([]byte, seedMinLength)
	n, err := rand.Read(seed)
	require.Equal(t, n, seedMinLength)
	require.NoError(t, err)
	sk, err := flow_crypto.GeneratePrivateKey(algorithmName, seed)
	require.NoError(t, err)

	publicKey := sk.PublicKey()

	return sk, publicKey
}

var transferTokensToAccountTemplate = `
import FungibleToken from 0x%s
import FlowToken from 0x%s

transaction(amount: UFix64, recipient: Address) {
  let sentVault: @FungibleToken.Vault
  prepare(signer: AuthAccount) {
    let vaultRef = signer.borrow<&FlowToken.Vault>(from: /storage/flowTokenVault)
      ?? panic("failed to borrow reference to sender vault")

    self.sentVault <- vaultRef.withdraw(amount: amount)
  }

  execute {
    let receiverRef =  getAccount(recipient)
      .getCapability(/public/flowTokenReceiver)
      .borrow<&{FungibleToken.Receiver}>()
        ?? panic("failed to borrow reference to recipient vault")

    receiverRef.deposit(from: <-self.sentVault)
  }
}
`

func FundAccount(flowClient access.Client, e Environment, recip flow.Address, amount float64, authorizer flow.Address, accountKey *flow.AccountKey, log *log.Logger) *flow.Transaction {
	referenceBlockID := utils.GetReferenceBlockId(flowClient, log)

	// Take this from ENV
	fungibleTokenAddress := flow.HexToAddress("0x" + e.FungibleTokenAddress)
	flowTokenAddress := flow.HexToAddress("0x" + e.FlowTokenAddress)

	recipient := cadence.NewAddress(recip)
	uintAmount := uint64(amount * sema.Fix64Factor)
	cadenceAmount := cadence.UFix64(uintAmount)

	fundAccountTx :=
		flow.NewTransaction().
			SetScript([]byte(fmt.Sprintf(transferTokensToAccountTemplate, fungibleTokenAddress, flowTokenAddress))).
			AddAuthorizer(authorizer).
			AddRawArgument(jsoncdc.MustEncode(cadenceAmount)).
			AddRawArgument(jsoncdc.MustEncode(recipient)).
			SetProposalKey(authorizer, accountKey.Index, accountKey.SequenceNumber).
			SetReferenceBlockID(referenceBlockID).
			SetPayer(authorizer)

	return fundAccountTx
}

func SetupAccount(client access.Client, e Environment, address flow.Address, accountKey *flow.AccountKey, log *log.Logger) *flow.Transaction {
	referenceBlockID := utils.GetReferenceBlockId(client, log)
	tx := flow.NewTransaction().
		SetScript(GenerateSetupAccount(e)).
		SetGasLimit(9999).
		SetProposalKey(address, accountKey.Index, accountKey.SequenceNumber).
		SetReferenceBlockID(referenceBlockID).
		SetPayer(address).
		AddAuthorizer(address)

	return tx
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

func CreatePiggy(client access.Client, e Environment, address flow.Address, metadata map[string]string, accountKey *flow.AccountKey, log *log.Logger) (*flow.Transaction, error) {

	referenceBlockID := utils.GetReferenceBlockId(client, log)

	tx := flow.NewTransaction().
		SetScript(GenerateCreatePiggy(e)).
		SetGasLimit(9999).
		SetProposalKey(address, accountKey.Index, accountKey.SequenceNumber).
		SetReferenceBlockID(referenceBlockID).
		SetPayer(address).
		AddAuthorizer(address)

	err := tx.AddArgument(CadenceMapStringString(metadata))
	if err != nil {
		return nil, err
	}
	return tx, nil

}

func mintDonation(b *emulator.Blockchain, e Environment, address flow.Address, piggyID int, donationComment string, recipientAddress flow.Address) (*flow.Transaction, error) {
	tx := createTxWithTemplateAndAuthorizer(b,
		GenerateMintDonation(e),
		address)

	err := tx.AddArgument(cadence.NewUInt64(uint64(piggyID)))
	if err != nil {
		return nil, err
	}
	err = tx.AddArgument(CadenceString(donationComment))
	if err != nil {
		return nil, err
	}
	err = tx.AddArgument(CadenceAddress(recipientAddress))
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func MintDonation(client access.Client, e Environment, address flow.Address, piggyID int, donationComment string, recipientAddress flow.Address, accountKey *flow.AccountKey, log *log.Logger) (*flow.Transaction, error) {
	referenceBlockID := utils.GetReferenceBlockId(client, log)

	tx := flow.NewTransaction().
		SetScript(GenerateMintDonation(e)).
		SetGasLimit(9999).
		SetProposalKey(address, accountKey.Index, accountKey.SequenceNumber).
		SetReferenceBlockID(referenceBlockID).
		SetPayer(address).
		AddAuthorizer(address)

	err := tx.AddArgument(cadence.NewUInt32(uint32(piggyID)))
	if err != nil {
		return nil, err
	}
	err = tx.AddArgument(CadenceString(donationComment))
	if err != nil {
		return nil, err
	}
	err = tx.AddArgument(CadenceAddress(recipientAddress))
	if err != nil {
		return nil, err
	}
	return tx, nil
}
func setupAccount(b *emulator.Blockchain, e Environment, address flow.Address) *flow.Transaction {

	tx := createTxWithTemplateAndAuthorizer(b,
		GenerateSetupAccount(e),
		address)

	return tx
}

func transferAdmin(b *emulator.Blockchain, e Environment, address flow.Address) (*flow.Transaction, error) {
	tx := createTxWithTemplateAndAuthorizer(b,
		GenerateTransferAdmin(e),
		address)
	return tx, nil
}

func TranferAdmin(client access.Client, e Environment, address flow.Address, accountKey *flow.AccountKey, log *log.Logger) (*flow.Transaction, error) {
	referenceBlockID := utils.GetReferenceBlockId(client, log)

	tx := flow.NewTransaction().
		SetScript(GenerateTransferAdmin(e)).
		SetGasLimit(9999).
		SetProposalKey(address, accountKey.Index, accountKey.SequenceNumber).
		SetReferenceBlockID(referenceBlockID).
		SetPayer(address).
		AddAuthorizer(address)

	return tx, nil
}
