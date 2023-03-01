import PiggyBanks from 0xPIGGYADDRESS
import PiggyBankAdminReceiver from 0xADMINRECEIVERADDRESS

// this transaction takes a Piggy bank Admin resource and 
// saves it to the account storage of the account
// where the contract is deployed
transaction {

    // Local variable for the topshot Admin object
    let adminRef: @PiggyBanks.Admin

    prepare(acct: AuthAccount) {

        self.adminRef <- acct.load acct.borrow<&PiggyBanks.Admin>(from: /storage/PiggyBanksAdmin)
            ?? panic("Could not borrow a reference to the Admin resource")
    }

    execute {

        PiggyBankAdminReceiver.storeAdmin(newAdmin: <-self.adminRef)
        
    }
}