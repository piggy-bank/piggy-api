package flow

import (
	"fmt"
	"strings"

	_ "github.com/kevinburke/go-bindata"

	_ "github.com/kevinburke/go-bindata"
)

/// This package contains utility functions to get contract code for the contracts in this repo
/// To use this package, import the `flow-core-contracts/lib/go/contracts` package,
/// then use the contracts package to call one of these functions.
/// They will return the byte array version of the contract.
///
/// Example
///
/// flowTokenCode := contracts.FlowToken(fungibleTokenAddr)
///

const (
	nonFungibleTokenFilename   = "blockchain/contracts/NonFungibleToken.cdc"
	metadataViewsFilename      = "blockchain/contracts/MetadataViews.cdc"
	piggyFilename              = "blockchain/contracts/piggy.cdc"
	flowServiceAccountFilename = "../blockchain-api/FlowServiceAccount.cdc"
	// Test contracts

	// Each contract has placeholder addresses that need to be replaced
	// depending on which network they are being used with
	placeholderNonFungibleTokenAddress = "0xNONFUNGIBLETOKENADDRESS"
	placeholderMetadataViewsAddress    = "0xMETADATAVIEWSADDRESS"
	placeholderTicketAddress           = "0xPIGGYADDRESS"
)

// Adds a `0x` prefix to the provided address string
func withHexPrefix(address string) string {
	if address == "" {
		return ""
	}

	if address[0:2] == "0x" {
		return address
	}

	return fmt.Sprintf("0x%s", address)
}

// NonFungibleToken returns the NonFungibleToken contract.
//
// The returned contract will import the NonFungibleToken contract from the specified address.
func NonFungibleToken() []byte {
	code := MustAssetString(nonFungibleTokenFilename)
	return []byte(code)
}

// MetadataViews returns the MetadataViews contract.
//
// The returned contract will import the MetadataViews contract from the specified addresses.
func MetadataViews() []byte {
	code := MustAssetString(metadataViewsFilename)
	return []byte(code)
}

func Piggy(nonFungibleTokenAddress, metadataViewsAddress string) []byte {
	code := MustAssetString(piggyFilename)

	code = strings.ReplaceAll(
		code,
		placeholderNonFungibleTokenAddress,
		withHexPrefix(nonFungibleTokenAddress),
	)

	code = strings.ReplaceAll(
		code,
		placeholderMetadataViewsAddress,
		withHexPrefix(metadataViewsAddress),
	)

	return []byte(code)
}

// FlowServiceAccount returns the FlowServiceAccount contract.
//
// The returned contract will import the FungibleToken, FlowToken, FlowFees, and FlowStorageFees
// contracts from the specified addresses.
func FlowServiceAccount(nonFungibleTokenAddress, metadataViewsAddress, ticketAddress string) []byte {
	code := MustAssetString(flowServiceAccountFilename)

	code = strings.ReplaceAll(
		code,
		placeholderNonFungibleTokenAddress,
		withHexPrefix(nonFungibleTokenAddress),
	)

	code = strings.ReplaceAll(
		code,
		placeholderMetadataViewsAddress,
		withHexPrefix(metadataViewsAddress),
	)

	code = strings.ReplaceAll(
		code,
		placeholderTicketAddress,
		withHexPrefix(ticketAddress),
	)

	return []byte(code)
}

func replaceAddresses(code string, env Environment) string {

	code = strings.ReplaceAll(
		code,
		placeholderNonFungibleTokenAddress,
		withHexPrefix(env.NonFungibleTokenAddress),
	)

	code = strings.ReplaceAll(
		code,
		placeholderMetadataViewsAddress,
		withHexPrefix(env.MetadataViewsAddress),
	)

	code = strings.ReplaceAll(
		code,
		placeholderTicketAddress,
		withHexPrefix(env.PiggyAddress),
	)
	return code
}
