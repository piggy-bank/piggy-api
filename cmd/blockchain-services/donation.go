package blockchainservices

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/manubidegain/piggy-api/cmd/api/configuration"
	flowUtils "github.com/manubidegain/piggy-api/flow"
	"github.com/manubidegain/piggy-api/utils"
)

// Make donation
// Get donations

func MintDonation(userAddress string, donationComment string, piggyID uint, config *configuration.FlowConfig,
	profile string, ctx *gin.Context, log *log.Logger, projectConfig *configuration.ProjectConfig) (uint64, error) {
	flowClient, err := utils.ConnectToFlow(profile, config)
	if err != nil {
		msg := "Cannot connect to flow" + err.Error()
		panic(msg)
	}
	env := flowUtils.NewEnv(profile)
	recipient := utils.GetAccount(ctx, flowClient, userAddress)
	serviceAcctAddr, serviceAcctKey, signer := utils.GetServiceAccount(flowClient, config, profile)
	recipientAddress := recipient.Address
	//recipientSigner, _ := crypto.NewInMemorySigner(recipientPrivateKey, recipientAcctKey.HashAlgo)

	mintTicketTx, err := flowUtils.MintDonation(flowClient, env, serviceAcctAddr, int(piggyID), donationComment, recipientAddress, serviceAcctKey, log)
	err = utils.HandleAndLogError(log, err)
	if err != nil {
		return 0, err
	}

	err = mintTicketTx.SignEnvelope(serviceAcctAddr, serviceAcctKey.Index, signer)
	err = utils.HandleAndLogError(log, err)
	if err != nil {
		return 0, err
	}

	err = flowClient.SendTransaction(ctx, *mintTicketTx)
	if err != nil {
		return 0, err
	}

	mintTxResp := utils.WaitForSeal(ctx, flowClient, mintTicketTx.ID())

	var DonationID uint64

	for _, event := range mintTxResp.Events {
		if strings.Contains(event.Type, "Minted") {
			value := event.Value.Fields[0].ToGoValue()
			if uint64value, ok := value.(uint64); ok {
				DonationID = uint64value
			}
		}
	}

	return DonationID, nil

}
