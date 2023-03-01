import FungibleToken from 0x9a0766d93b6608b7
import NonFungibleToken from 0x631e88ae7f1d7c20
import MetadataViews from 0x631e88ae7f1d7c20

pub contract PiggyBanks: NonFungibleToken {
    // The network the contract is deployed on
    //pub fun Network() : String { return ${NETWORK} }

    // -----------------------------------------------------------------------
    // PiggyBanks contract Events
    // -----------------------------------------------------------------------
    // Emitted when the PiggyBanks contract is created
    pub event ContractInitialized()

    // Emitted when a new Piggy Bank struct is created
    pub event PiggyCreated(id: UInt32, metadata: {String:String})
    // Emitted when a Donation is made to a piggy and a conmemorative NFT is minted
    pub event DonationMinted(donationID: UInt64, piggyID: UInt32, serialNumber: UInt32)
    // Emitted when a break a piggy, piggy will not be deleted.
    pub event SetBroken(piggyID: UInt32)
    // Emitted when a Donation is destroyed
    pub event DonationDestroyed(id: UInt64)
    // Events for Collection-related actions
    //
    // Emitted when a donation is withdrawn from a Collection
    pub event Withdraw(id: UInt64, from: Address?)
    // Emitted when a donation is deposited into a Collection
    pub event Deposit(id: UInt64, to: Address?)

    

    // -----------------------------------------------------------------------
    // PiggyBanks contract-level fields.
    // These contain actual values that are stored in the smart contract.
    // -----------------------------------------------------------------------
   
    // Variable size dictionary of Piggy structs
    access(self) var piggiesDatas: {UInt32: Piggy}

     // Mapping of Piggy IDs that indicates the number of Donations     
    // that have been done for specific Piggies.
    // When a Donation is made, this value is stored in the Donation to
    // show its place in the Piggy, eg. 13 of 60.
    access(contract) var numberOfDonationsPerPig: {UInt32: UInt32}

    // The ID that is used to create Piggies. 
    // Every time a Piggy is created, piggyID is assigned 
    // to the new Piggy's ID and then is incremented by 1.
    pub var nextpiggyID: UInt32

    pub var totalSupply: UInt64

    pub struct Piggy {

        // The unique ID for the Piggy Bank
        pub let piggyID: UInt32

        // Stores all the metadata about the piggy bank as a string mapping
        // This is not the long term way NFT metadata will be stored. It's a temporary
        // construct while we figure out a better way to do metadata.
        //
        pub let metadata: {String: String}

        // Indicates if the Piggy is currently broken.
        // When a Piggy is created, it is not broken
        // and Donations are allowed to be made to it.
        // When a piggy is broken, Donations cannot be done.
        // A Piggy can never be changed from broken to unbroken,
        // the decision to break a Piggy it is final.
        pub var broken: Bool

        // Piggy collected amount, will be zero until pig is broken.
        // IN CENTS
        pub var collectedAmount: UInt64

        // Piggy royalty to breaker.
        // IN CENTS
        pub var breakerRoyalty: UInt64

        init(metadata: {String: String}) {
            pre {
                metadata.length != 0: "New Piggy metadata cannot be empty"
            }
            self.piggyID = PiggyBanks.nextpiggyID
            self.metadata = metadata
            self.collectedAmount = 0
            self.breakerRoyalty = 0
            self.broken = false
            PiggyBanks.piggiesDatas[self.piggyID] = self
        }

        // breakPiggy() break the Pig so that no more Donations can be made
        //
        // Pre-Conditions:
        // The Piggy should not be broken
        pub fun breakPiggy(collectedAmount: UInt64, breakerRoyalty : UInt64) {
            if !self.broken {
                self.broken = true
                self.collectedAmount = collectedAmount
                self.breakerRoyalty = breakerRoyalty
                emit SetBroken(piggyID: self.piggyID)
            }
        }

        pub fun mintDonation(piggyID: UInt32, donationComment : String): @NFT {
            pre {
                !self.broken: "Cannot make the donation to the Piggy because the piggy is broken."           
            }

            // Gets the number of Donation that have been made for this Piggy
            // to use as this Donation's serial number
            let numInPiggy = PiggyBanks.numberOfDonationsPerPig[piggyID]!

            // Mint the new Donation
            let newDonation: @NFT <- create NFT(serialNumber: numInPiggy + UInt32(1),
                                              piggyID: piggyID, donationComment: donationComment)

            // Increment the count of Donations minted for this Piggy
            PiggyBanks.numberOfDonationsPerPig[piggyID] = numInPiggy + UInt32(1)

            return <-newDonation
        }

        pub fun getBroken(): Bool {
            return self.broken
        }

        pub fun getNumOfDonations(piggyID: UInt32): UInt32? {
            return PiggyBanks.numberOfDonationsPerPig[self.piggyID]
        }
    }

    pub struct DonationData {
        // The ID of the Piggy that the Donation references
        pub let piggyID: UInt32

        // The place in the piggy that this Donation was minted
        // Otherwise know as the serial number
        pub let serialNumber: UInt32

        // Donation reason, by default is empty
        pub let donationComment : String

        init(piggyID: UInt32, serialNumber: UInt32, donationComment: String) {
            self.piggyID = piggyID
            self.serialNumber = serialNumber
            self.donationComment = donationComment
        }

    }
    
    // The resource that represents the Donation NFTs
    //
    pub resource NFT: NonFungibleToken.INFT, MetadataViews.Resolver {

        // Global unique donation ID
        pub let id: UInt64
        
        // Struct of Donation metadata
        pub let data: DonationData

        init(serialNumber: UInt32, piggyID: UInt32, donationComment: String) {
            // Increment the global Moment IDs
            PiggyBanks.totalSupply = PiggyBanks.totalSupply + UInt64(1)

            self.id = PiggyBanks.totalSupply

            // Set the metadata struct
            self.data = DonationData(piggyID: piggyID, serialNumber: serialNumber, donationComment: donationComment)

            emit DonationMinted(donationID: self.id,
                              piggyID: piggyID,
                              serialNumber: self.data.serialNumber)        }

        // If the Donation is destroyed, emit an event to indicate 
        // to outside ovbservers that it has been destroyed
        destroy() {
            emit DonationDestroyed(id: self.id)
        }

        // Add more metadata, for example. donation for piggy reason or whatever, define all metadata.
        pub fun name(): String {
            return "Donation number"
                .concat(self.data.serialNumber.toString())
                .concat("for piggy with id :")
                .concat(self.data.serialNumber.toString())
        }

        pub fun description(): String {
            return self.data.donationComment.length > 0 ? self.data.donationComment : "No reason"
        }

        // All supported metadata views for the Moment including the Core NFT Views
        pub fun getViews(): [Type] {
            return [
                Type<MetadataViews.Display>(),
                Type<MetadataViews.ExternalURL>(),
                Type<MetadataViews.NFTCollectionDisplay>(),
                Type<MetadataViews.Serial>(),
                Type<MetadataViews.Medias>()
            ]
        }

       

        pub fun resolveView(_ view: Type): AnyStruct? {
            switch view {
                case Type<MetadataViews.Display>():
                    return MetadataViews.Display(
                        name: self.name(),
                        description: self.description(),
                        thumbnail: MetadataViews.HTTPFile(url: self.thumbnail())
                    )
                case Type<MetadataViews.Serial>():
                    return MetadataViews.Serial(
                        UInt64(self.data.serialNumber)
                    )
                case Type<MetadataViews.ExternalURL>():
                    return MetadataViews.ExternalURL(self.getDonationURL())
                case Type<MetadataViews.NFTCollectionDisplay>():
                    let bannerImage = MetadataViews.Media(
                        file: MetadataViews.HTTPFile(
                            url: "https://gopiggy.com/static/img/some-image.svg"
                        ),
                        mediaType: "image/svg+xml"
                    )
                    let squareImage = MetadataViews.Media(
                        file: MetadataViews.HTTPFile(
                            url: "https://gopiggy.com/static/img/some-image.png"
                        ),
                        mediaType: "image/png"
                    )
                    return MetadataViews.NFTCollectionDisplay(
                        name: "PiggyBanks",
                        description: "Piggy Bank is to chance to change someone life while you get a valuable NFT and a chance to obtain more than your gift!!!",
                        externalURL: MetadataViews.ExternalURL("https://gopiggy.com"),
                        squareImage: squareImage,
                        bannerImage: bannerImage,
                        socials: {
                            "twitter": MetadataViews.ExternalURL("https://twitter.com/gopiggy"),
                            "discord": MetadataViews.ExternalURL("https://discord.com/invite/gopiggy"),
                            "instagram": MetadataViews.ExternalURL("https://www.instagram.com/gopiggy")
                        }
                    )
                case Type<MetadataViews.Medias>():
                    return MetadataViews.Medias(
                        items: [
                            MetadataViews.Media(
                                file: MetadataViews.HTTPFile(
                                    url: self.mediumimage()
                                ),
                                mediaType: "image/jpeg"
                            ),
                            MetadataViews.Media(
                                file: MetadataViews.HTTPFile(
                                    url: self.video()
                                ),
                                mediaType: "video/mp4"
                            )
                        ]
                    )
            }

            return nil
        }

        // getMomentURL 
        // Returns: The computed external url of the donation
        pub fun getDonationURL(): String {
            return "https://gopiggy.com/donation/".concat(self.id.toString())
        }

        pub fun assetPath(): String {
            return "https://assets.gopiggy.com/media/".concat(self.id.toString())
        }

        // returns a url to disPiggy an medium sized image
        pub fun mediumimage(): String {
            let url = self.assetPath().concat("?width=512")
            return self.appendOptionalParams(url: url, firstDelim: "&")
        }

        // a url to disPiggy a thumbnail associated with the donation
        pub fun thumbnail(): String {
            let url = self.assetPath().concat("?width=256")
            return self.appendOptionalParams(url: url, firstDelim: "&")
        }

        // a url to disPiggy a video associated with the donation
        pub fun video(): String {
            let url = self.assetPath().concat("/video")
            return self.appendOptionalParams(url: url, firstDelim: "?")
        }

        // appends and optional network param needed to resolve the media
        pub fun appendOptionalParams(url: String, firstDelim: String): String {
            //if (PiggyBanks.Network() == "testnet") {
            //    return url.concat(firstDelim).concat("testnet")
            //}
            return url
        }
    }

   
    
    // Admin is a special authorization resource that 
    // allows the owner to perform important functions to modify the 
    // various aspects of the Piggies and donations
    //
    pub resource Admin {

        // createPiggy creates a new Piggy struct 
        // and stores it in the Piggies dictionary in the PiggyBanks smart contract
        //
        // Parameters: metadata: A dictionary mapping metadata titles to their data
        //                       example: {"ToBe": "Defined", "When": "ASAP"}
        //
        // Returns: the ID of the new Piggy object
        //
        pub fun createPiggy(metadata: {String: String}): UInt32 {
            // Create the new Piggy
            var newPiggy = Piggy(metadata: metadata)
            let newID = newPiggy.piggyID

            // Increment the ID so that it isn't used again
            PiggyBanks.nextpiggyID = PiggyBanks.nextpiggyID + UInt32(1)
            PiggyBanks.numberOfDonationsPerPig[newPiggy.piggyID] = 0

            emit PiggyCreated(id: newPiggy.piggyID, metadata: metadata)

            // Store it in the contract storage
            PiggyBanks.piggiesDatas[newID] = newPiggy

            return newID
        }

         pub fun borrowPiggy(piggyID: UInt32): &Piggy {
            pre {
                PiggyBanks.piggiesDatas[piggyID] != nil: "Cannot borrow Piggy: Piggy doesn't exist"
            }
            
            // Get a reference to the event and return it
            // use `&` to indicate the reference to the object and type
            return (&PiggyBanks.piggiesDatas[piggyID] as &Piggy?)!
        }

        // createNewAdmin creates a new Admin resource
        //
        pub fun createNewAdmin(): @Admin {
            return <-create Admin()
        }
    }

    // This is the interface that users can cast their Donations Collection as
    // to allow others to deposit Donations into their Collection. It also allows for reading
    // the IDs of Donations in the Collection.
    pub resource interface DonationCollectionPublic {
        pub fun deposit(token: @NonFungibleToken.NFT)
        pub fun batchDeposit(tokens: @NonFungibleToken.Collection)
        pub fun getIDs(): [UInt64]
        pub fun borrowNFT(id: UInt64): &NonFungibleToken.NFT
        pub fun borrowDonation(id: UInt64): &PiggyBanks.NFT? {
            // If the result isn't nil, the id of the returned reference
            // should be the same as the argument to the function
            post {
                (result == nil) || (result?.id == id): 
                    "Cannot borrow Donation reference: The ID of the returned reference is incorrect"
            }
        }
    }

    // Collection is a resource that every user who owns NFTs 
    // will store in their account to manage their NFTS
    //
    pub resource Collection: DonationCollectionPublic, NonFungibleToken.Provider, NonFungibleToken.Receiver, NonFungibleToken.CollectionPublic, MetadataViews.ResolverCollection { 
        // Dictionary of Donations conforming tokens
        // NFT is a resource type with a UInt64 ID field
        pub var ownedNFTs: @{UInt64: NonFungibleToken.NFT}

        init() {
            self.ownedNFTs <- {}
        }

        // withdraw removes an Donation from the Collection and moves it to the caller
        //
        // Parameters: withdrawID: The ID of the NFT 
        // that is to be removed from the Collection
        //
        // returns: @NonFungibleToken.NFT the token that was withdrawn
        pub fun withdraw(withdrawID: UInt64): @NonFungibleToken.NFT {

            // Borrow nft and check if locked
            let nft = self.borrowNFT(id: withdrawID)

            // Remove the nft from the Collection
            let token <- self.ownedNFTs.remove(key: withdrawID) 
                ?? panic("Cannot withdraw: Donation does not exist in the collection")

            emit Withdraw(id: token.id, from: self.owner?.address)
            
            // Return the withdrawn token
            return <-token
        }

        // batchWithdraw withdraws multiple tokens and returns them as a Collection
        //
        // Parameters: ids: An array of IDs to withdraw
        //
        // Returns: @NonFungibleToken.Collection: A collection that contains
        //                                        the withdrawn moments
        //
        pub fun batchWithdraw(ids: [UInt64]): @NonFungibleToken.Collection {
            // Create a new empty Collection
            var batchCollection <- create Collection()
            
            // Iterate through the ids and withdraw them from the Collection
            for id in ids {
                batchCollection.deposit(token: <-self.withdraw(withdrawID: id))
            }
            
            // Return the withdrawn tokens
            return <-batchCollection
        }

        // deposit takes a Moment and adds it to the Collections dictionary
        //
        // Paramters: token: the NFT to be deposited in the collection
        //
        pub fun deposit(token: @NonFungibleToken.NFT) {
            
            // Cast the deposited token as a PiggyBanks NFT to make sure
            // it is the correct type
            let token <- token as! @PiggyBanks.NFT

            // Get the token's ID
            let id = token.id

            // Add the new token to the dictionary
            let oldToken <- self.ownedNFTs[id] <- token

            // Only emit a deposit event if the Collection 
            // is in an account's storage
            if self.owner?.address != nil {
                emit Deposit(id: id, to: self.owner?.address)
            }

            // Destroy the empty old token that was "removed"
            destroy oldToken
        }

        // batchDeposit takes a Collection object as an argument
        // and deposits each contained NFT into this Collection
        pub fun batchDeposit(tokens: @NonFungibleToken.Collection) {

            // Get an array of the IDs to be deposited
            let keys = tokens.getIDs()

            // Iterate through the keys in the collection and deposit each one
            for key in keys {
                self.deposit(token: <-tokens.withdraw(withdrawID: key))
            }

            // Destroy the empty Collection
            destroy tokens
        }

        // getIDs returns an array of the IDs that are in the Collection
        pub fun getIDs(): [UInt64] {
            return self.ownedNFTs.keys
        }

        // borrowNFT Returns a borrowed reference to a Donation in the Collection
        // so that the caller can read its ID
        //
        // Parameters: id: The ID of the NFT to get the reference for
        //
        // Returns: A reference to the NFT
        //
        // Note: This only allows the caller to read the ID of the NFT,
        // not any PiggyBanks specific data. Please use borrowDonation to 
        // read Donation data.
        //
        pub fun borrowNFT(id: UInt64): &NonFungibleToken.NFT {
            return (&self.ownedNFTs[id] as &NonFungibleToken.NFT?)!
        }

        // Safe way to borrow a reference to an NFT that does not panic
        // Also now part of the NonFungibleToken.PublicCollection interface
        //
        // Parameters: id: The ID of the NFT to get the reference for
        //
        // Returns: An optional reference to the desired NFT, will be nil if the passed ID does not exist
        pub fun borrowNFTSafe(id: UInt64): &NonFungibleToken.NFT? {
            if let nftRef = &self.ownedNFTs[id] as &NonFungibleToken.NFT? {
                return nftRef
            }
            return nil
        }

        // borrowDonation returns a borrowed reference to a Donation
        // so that the caller can read data and call methods from it.
        // They can use this to read its piggyID, serialNumber,
        // or any of the Piggy or the donation data associated with it
        //
        // Parameters: id: The ID of the NFT to get the reference for
        //
        // Returns: A reference to the NFT
        pub fun borrowDonation(id: UInt64): &PiggyBanks.NFT? {
            if self.ownedNFTs[id] != nil {
                let ref = (&self.ownedNFTs[id] as auth &NonFungibleToken.NFT?)!
                return ref as! &PiggyBanks.NFT
            } else {
                return nil
            }
        }

        pub fun borrowViewResolver(id: UInt64): &AnyResource{MetadataViews.Resolver} {
            let nft = (&self.ownedNFTs[id] as auth &NonFungibleToken.NFT?)! 
            let PiggyBanksNFT = nft as! &PiggyBanks.NFT
            return PiggyBanksNFT as &AnyResource{MetadataViews.Resolver}
        }

        // If a transaction destroys the Collection object,
        // All the NFTs contained within are also destroyed!
        //
        destroy() {
            destroy self.ownedNFTs
        }
    }

    // -----------------------------------------------------------------------
    // PiggyBanks contract-level function definitions
    // -----------------------------------------------------------------------
    // createEmptyCollection creates a new, empty Collection object so that
    // a user can store it in their account storage.
    // Once they have a Collection in their storage, they are able to receive
    // Moments in transactions.
    //
    pub fun createEmptyCollection(): @NonFungibleToken.Collection {
        return <-create PiggyBanks.Collection()
    }

    // getAllPiggies returns all the piggies in PiggyBanks
    //
    // Returns: An array of all the piggies that have been created
    pub fun getAllPiggies(): [PiggyBanks.Piggy] {
        return PiggyBanks.piggiesDatas.values
    }

    // getPiggyMetaData returns all the metadata associated with a specific Piggy
    // 
    // Parameters: piggyID: The id of the Piggy that is being searched
    //
    // Returns: The metadata as a String to String mapping optional
    pub fun getPiggyMetaData(piggyID: UInt32): {String: String}? {
        return self.piggiesDatas[piggyID]?.metadata
    }

    // getPiggyMetaDataByField returns the metadata associated with a 
    //                        specific field of the metadata
    // 
    // Parameters: piggyID: The id of the Piggy that is being searched
    //             field: The field to search for
    //
    // Returns: The metadata field as a String Optional
    pub fun getPiggyMetaDataByField(piggyID: UInt32, field: String): String? {
        // Don't force a revert if the piggyID or field is invalid
        if let piggy = PiggyBanks.piggiesDatas[piggyID] {
            return piggy.metadata[field]
        } else {
            return nil
        }
    }

    // -----------------------------------------------------------------------
    // PiggyBanks initialization function
    // -----------------------------------------------------------------------
    //
    init() {
        // Initialize contract fields
        self.piggiesDatas = {}
        self.numberOfDonationsPerPig = {}
        self.nextpiggyID = 1
        self.totalSupply = 0

        // Put a new Collection in storage
        self.account.save<@Collection>(<- create Collection(), to: /storage/DonationCollection)

        // Create a public capability for the Collection
        self.account.link<&{DonationCollectionPublic}>(/public/DonationCollection, target: /storage/DonationCollection)

        // Put the Minter in storage
        self.account.save<@Admin>(<- create Admin(), to: /storage/PiggyBanksAdmin)

        emit ContractInitialized()
    }

}