import PiggyBanks from 0xPIGGYADDRESS

// This transaction is what an admin would use to mint a single new Donation
// and deposit it in a user's collection
// Parameters:
//
// piggyID: the ID of a event containing the target play
// donationComment: an string related to the donation
// recipientAddr: the Flow address of the account receiving the newly minted donation

transaction(piggyID: UInt32, donationComment: String, recipientAddr: Address) {
    // local variable for the admin reference
    let adminRef: &PiggyBanks.Admin

    prepare(acct: AuthAccount) {
        // borrow a reference to the Admin resource in storage
        self.adminRef = acct.borrow<&PiggyBanks.Admin>(from: /storage/PiggyBanksAdmin)
            ?? panic("Could not borrow a reference to the Admin resource")
    }

    execute {
        // Borrow a reference to the specified event
        let piggyRef = self.adminRef.borrowPiggy(piggyID: piggyID)

        // Mint a new NFT
        let donation1 <- piggyRef.mintDonation(piggyID: piggyID, donationComment: donationComment)

        // get the public account object for the recipient
        let recipient = getAccount(recipientAddr)

        // get the Collection reference for the receiver
        let receiverRef = recipient.getCapability(/public/DonationCollection).borrow<&{PiggyBanks.DonationCollectionPublic}>()
            ?? panic("Cannot borrow a reference to the recipient's collection")

        // deposit the NFT in the receivers collection
        receiverRef.deposit(token: <-donation1)
    }
}