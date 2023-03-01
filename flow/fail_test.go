package flow

/*
import (
	"testing"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	sdktemplates "github.com/onflow/flow-go-sdk/templates"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFails(t *testing.T) {

	b, k, e := newTestSetup(t)

	// OUR CONTRACTS
	nonFungibleTokenAccountKey, ourSigner := k.NewWithSigner()
	nonFungibleTokenCode := NonFungibleToken()
	nonFungibleTokenAddress, err := b.CreateAccount([]*flow.AccountKey{nonFungibleTokenAccountKey}, []sdktemplates.Contract{
		{
			Name:   "NonFungibleToken",
			Source: string(nonFungibleTokenCode),
		},
	})
	assert.NoError(t, err)
	_, err = b.CommitBlock()
	assert.NoError(t, err)

	metaDataViewsCode := MetadataViews()
	metaDataAddress, err := b.CreateAccount([]*flow.AccountKey{nonFungibleTokenAccountKey}, []sdktemplates.Contract{
		{
			Name:   "MetadataViews",
			Source: string(metaDataViewsCode),
		},
	})
	assert.NoError(t, err)
	_, err = b.CommitBlock()
	assert.NoError(t, err)

	ticketCode := Ticket(nonFungibleTokenAddress.String(), metaDataAddress.String())
	ticketAddress, err := b.CreateAccount([]*flow.AccountKey{nonFungibleTokenAccountKey}, []sdktemplates.Contract{
		{
			Name:   "Tickets",
			Source: string(ticketCode),
		},
	})
	assert.NoError(t, err)
	_, err = b.CommitBlock()
	assert.NoError(t, err)

	anotherAccountKey, anotherSigner := k.NewWithSigner()
	anotherAddress, err := b.CreateAccount([]*flow.AccountKey{anotherAccountKey}, []sdktemplates.Contract{})
	assert.NoError(t, err)
	_, err = b.CommitBlock()
	assert.NoError(t, err)

	e.NonFungibleTokenAddress = nonFungibleTokenAddress.String()
	e.MetadataViewsAddress = metaDataAddress.String()
	e.TicketAddress = ticketAddress.String()
	e.AnotherAccountAddress = anotherAddress.String()

	t.Run("Trying to mint a ticket whitout an event nor metadata...", func(t *testing.T) {
		tx, err := mintTicket(b, e, ticketAddress, 1, 1, "123.0", ticketAddress)
		require.NoError(t, err)
		cryptoSigner, err := b.ServiceKey().Signer()

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, ticketAddress},
			[]crypto.Signer{cryptoSigner, ourSigner},
			true,
		)
	})

	t.Run("Creating new event.. ", func(t *testing.T) {

		ourMap := make(map[string]int)
		ourMap["VIP"] = 2
		ourMap["GA"] = 2
		tx, err := createEvent(b, e, ticketAddress, 1, "FirstBlockchainTestingIujuuu", 24, ourMap)
		require.NoError(t, err)
		cryptoSigner, err := b.ServiceKey().Signer()

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, ticketAddress},
			[]crypto.Signer{cryptoSigner, ourSigner},
			false,
		)
	})

	t.Run("Trying to mint a ticket whitout metadata...", func(t *testing.T) {
		tx, err := mintTicket(b, e, ticketAddress, 1, 1, "123.0", ticketAddress)
		require.NoError(t, err)
		cryptoSigner, err := b.ServiceKey().Signer()

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, ticketAddress},
			[]crypto.Signer{cryptoSigner, ourSigner},
			true,
		)
	})

	t.Run("Creating Metadata ...", func(t *testing.T) {

		ourMetadataMap := make(map[string]string)
		ourMetadataMap["SEAT"] = "12"
		ourMetadataMap["ROW"] = "1"
		tx, err := createMetadata(b, e, ticketAddress, 1, 1, "FirstMetadata", "FirstImage", "VIP", "EventType", "Venue", "EventDescription", "Date", ourMetadataMap)
		require.NoError(t, err)
		cryptoSigner, err := b.ServiceKey().Signer()

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, ticketAddress},
			[]crypto.Signer{cryptoSigner, ourSigner},
			false,
		)
	})

	t.Run("Trying to mint a ticket whit an Event that doesn't exist...", func(t *testing.T) {
		tx, err := mintTicket(b, e, ticketAddress, 99, 1, "123.0", ticketAddress)
		require.NoError(t, err)
		cryptoSigner, err := b.ServiceKey().Signer()

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, ticketAddress},
			[]crypto.Signer{cryptoSigner, ourSigner},
			true,
		)
	})

	t.Run("Trying to mint a ticket whit metadata that doesn't exist...", func(t *testing.T) {
		tx, err := mintTicket(b, e, ticketAddress, 1, 99, "123.0", ticketAddress)
		require.NoError(t, err)

		cryptoSigner, err := b.ServiceKey().Signer()
		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, ticketAddress},
			[]crypto.Signer{cryptoSigner, ourSigner},
			true,
		)
	})

	t.Run("Trying to mint mote tickets than the maximum allowed for VIP...", func(t *testing.T) {

		// First ticket
		tx, err := mintTicket(b, e, ticketAddress, 1, 1, "123.0", ticketAddress)
		require.NoError(t, err)
		cryptoSigner, err := b.ServiceKey().Signer()

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, ticketAddress},
			[]crypto.Signer{cryptoSigner, ourSigner},
			false,
		)

		// Second ticket
		tx2, err2 := mintTicket(b, e, ticketAddress, 1, 1, "123.0", ticketAddress)
		require.NoError(t, err2)

		signAndSubmit(
			t, b, tx2,
			[]flow.Address{b.ServiceKey().Address, ticketAddress},
			[]crypto.Signer{cryptoSigner, ourSigner},
			false,
		)

		// Third ticket
		tx3, err3 := mintTicket(b, e, ticketAddress, 1, 1, "123.0", ticketAddress)
		require.NoError(t, err3)

		signAndSubmit(
			t, b, tx3,
			[]flow.Address{b.ServiceKey().Address, ticketAddress},
			[]crypto.Signer{cryptoSigner, ourSigner},
			true,
		)

	})

	t.Run("Closing an inexistent event.. ", func(t *testing.T) {

		tx, err := closeEvent(b, e, ticketAddress, 99)
		require.NoError(t, err)

		cryptoSigner, err := b.ServiceKey().Signer()
		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, ticketAddress},
			[]crypto.Signer{cryptoSigner, ourSigner},
			true,
		)
	})

	t.Run("Trying transferring a ticket to an account that hasn’t been setup ...", func(t *testing.T) {

		tx, err := transferTicket(b, e, ticketAddress, anotherAddress, 0)
		require.NoError(t, err)

		cryptoSigner, err := b.ServiceKey().Signer()
		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, ticketAddress},
			[]crypto.Signer{cryptoSigner, ourSigner},
			true,
		)
	})

	t.Run("Setup another account...", func(t *testing.T) {

		tx := setupAccount(b, e, anotherAddress)
		cryptoSigner, _ := b.ServiceKey().Signer()
		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, anotherAddress},
			[]crypto.Signer{cryptoSigner, anotherSigner},
			false,
		)
	})

	t.Run("Trying to purchase ticket that doesnt exist...", func(t *testing.T) {

		tx, err := purchaseTicket(b, e, anotherAddress, ticketAddress, 9, "250.0")
		require.NoError(t, err)
		cryptoSigner, err := b.ServiceKey().Signer()

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, anotherAddress},
			[]crypto.Signer{cryptoSigner, anotherSigner},
			true,
		)
	})

	t.Run("Trying to create new event whit a non Admin account.. ", func(t *testing.T) {

		ourMap := make(map[string]int)
		ourMap["VIP"] = 2
		ourMap["GA"] = 2
		tx, err := createEvent(b, e, anotherAddress, 2, "FirstBlockchainTestingIujuuu", 24, ourMap)
		require.NoError(t, err)
		cryptoSigner, err := b.ServiceKey().Signer()

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, anotherAddress},
			[]crypto.Signer{cryptoSigner, anotherSigner},
			true,
		)
	})

	t.Run("Trying to create new metadata whit a non Admin account ...", func(t *testing.T) {

		ourMetadataMap := make(map[string]string)
		ourMetadataMap["SEAT"] = "12"
		ourMetadataMap["ROW"] = "1"
		tx, err := createMetadata(b, e, anotherAddress, 1, 8, "FirstMetadata", "FirstImage", "VIP", "EventType", "Venue", "EventDescription", "Date", ourMetadataMap)
		require.NoError(t, err)
		cryptoSigner, err := b.ServiceKey().Signer()

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, anotherAddress},
			[]crypto.Signer{cryptoSigner, anotherSigner},
			true,
		)
	})

	t.Run("Trying to purchase ticket from an account whitout privileges...", func(t *testing.T) {

		tx, err := purchaseTicket(b, e, ticketAddress, anotherAddress, 0, "250.0")
		require.NoError(t, err)
		cryptoSigner, err := b.ServiceKey().Signer()

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, anotherAddress},
			[]crypto.Signer{cryptoSigner, anotherSigner},
			true,
		)
	})

	t.Run("Trying to transfer ticket from an account whitout privileges...", func(t *testing.T) {

		tx, err := transferTicket(b, e, ticketAddress, anotherAddress, 0)
		require.NoError(t, err)
		cryptoSigner, err := b.ServiceKey().Signer()

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, anotherAddress},
			[]crypto.Signer{cryptoSigner, anotherSigner},
			true,
		)
	})

}
*/
