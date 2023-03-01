import NonFungibleToken from 0xNONFUNGIBLETOKENADDRESS
import PiggyBanks from 0xPIGGYADDRESS
import MetadataViews from 0xMETADATAVIEWSADDRESS

// This transaction sets up an account to use Piggy Banks
// by storing an empty ticket collection and creating
// a public capability for it
transaction {

    prepare(acct: AuthAccount) {

        // First, check to see if a moment collection already exists
        if acct.borrow<&PiggyBanks.Collection>(from: /storage/DonationCollection) == nil {

            // create a new Donation  Collection
            let collection <- PiggyBanks.createEmptyCollection() as! @PiggyBanks.Collection

            // Put the new Collection in storage
            acct.save(<-collection, to: /storage/DonationCollection)

            // create a public capability for the collection
            acct.link<&{NonFungibleToken.CollectionPublic, PiggyBanks.DonationCollectionPublic, MetadataViews.ResolverCollection}>(/public/DonationCollection, target: /storage/DonationCollection)
        }
    }
}