package flow

const (

	// USER
	setupAccountFilename = "blockchain/transactions/user/setup_account.cdc"

	// ADMIN
	createPiggyFilename   = "blockchain/transactions/admin/create_piggy.cdc"
	mintDonationFilename  = "blockchain/transactions/admin/mint_donation.cdc"
	transferAdminFilename = "blockchain/transactions/admin/transfer_admin.cdc"

	// SCRIPTS
	nextPiggyIDFilename    = "blockchain/transactions/scripts/get_nextPiggyID.cdc"
	getTotalSupplyFilename = "blockchain/transactions/scripts/get_totalSupply.cdc"
)

func GenerateSetupAccount(env Environment) []byte {
	code := MustAssetString(setupAccountFilename)

	return []byte(replaceAddresses(code, env))
}

func GenerateCreatePiggy(env Environment) []byte {
	code := MustAssetString(createPiggyFilename)

	return []byte(replaceAddresses(code, env))
}

func GenerateMintDonation(env Environment) []byte {
	code := MustAssetString(mintDonationFilename)

	return []byte(replaceAddresses(code, env))
}

func GenerateTransferAdmin(env Environment) []byte {
	code := MustAssetString(transferAdminFilename)

	return []byte(replaceAddresses(code, env))
}

func GenerateGetNextPiggyID(env Environment) []byte {
	code := MustAssetString(nextPiggyIDFilename)

	return []byte(replaceAddresses(code, env))
}

func GenerateGetTotalSupply(env Environment) []byte {
	code := MustAssetString(getTotalSupplyFilename)

	return []byte(replaceAddresses(code, env))
}
