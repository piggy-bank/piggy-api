import PiggyBanks from 0xPIGGYADDRESS


// This transaction is for the admin to create a new Piggy

transaction(metadata: {String: String}) {
    
    // Local variable for the piggy banks Admin object
    let adminRef: &PiggyBanks.Admin

    prepare(acct: AuthAccount) {

        // borrow a reference to the Admin resource in storage
        self.adminRef = acct.borrow<&PiggyBanks.Admin>(from: /storage/PiggyBanksAdmin)
            ?? panic("Could not borrow a reference to the Admin resource")
    }

    execute {
        
        // Create a set with the specified name
        self.adminRef.createPiggy(metadata:metadata)
    }
}