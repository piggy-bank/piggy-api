package blockchainservices

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/manubidegain/piggy-api/cmd/api/configuration"
	flowUtils "github.com/manubidegain/piggy-api/flow"
	"github.com/manubidegain/piggy-api/utils"
)

// Create piggy

func CreateBlockchainPiggy(userAddress string, name string, description string, config *configuration.FlowConfig, profile string, ctx *gin.Context, log *log.Logger) (uint32, error) {

	flowClient, err := utils.ConnectToFlow(profile, config)
	if err != nil {
		msg := "Cannot connect to flow" + err.Error()
		panic(msg)
	}
	env := flowUtils.NewEnv(profile)

	serviceAcctAddr, serviceAcctKey, signer := utils.GetServiceAccount(flowClient, config, profile)

	metadata := make(map[string]string)
	metadata["Name"] = name
	metadata["Description"] = description
	metadata["Creator"] = userAddress

	createEventTx, err := flowUtils.CreatePiggy(flowClient, env, serviceAcctAddr, metadata, serviceAcctKey, log)
	if err != nil {
		return 0, err
	}

	err = createEventTx.SignEnvelope(serviceAcctAddr, serviceAcctKey.Index, signer)
	if err != nil {
		return 0, err
	}

	err = flowClient.SendTransaction(ctx, *createEventTx)
	if err != nil {
		return 0, err
	}

	mintTxResp := utils.WaitForSeal(ctx, flowClient, createEventTx.ID())
	if mintTxResp.Error != nil {
		return 0, mintTxResp.Error
	}

	var piggyID uint32

	for _, event := range mintTxResp.Events {
		if strings.Contains(event.Type, "Piggy") {
			value := event.Value.Fields[0].ToGoValue()
			if uint32value, ok := value.(uint32); ok {
				piggyID = uint32value
			}
		}
	}
	return piggyID, nil

}
