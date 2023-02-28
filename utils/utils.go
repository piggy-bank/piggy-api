package utils

import (
	"os"
	"strings"

	"github.com/manubidegain/piggy-api/cmd/api/configuration"
)

const (
	Development              = "dev"
	Production               = "prod"
	Test                     = "test"
	configurationPackagePath = "configfiles"
)

func CalculateProfile() string {
	scope := os.Getenv("SCOPE")

	if scope == "" {
		msg := "can not start application without 'SCOPE' environment variable"
		panic(msg)
	}

	scope = strings.ToLower(scope)
	switch scope {
	case Production:
		return Production
	case Test:
		return Test
	default:
		return Development
	}
}

func BuildConfig(profile string) *configuration.Config {
	path := configurationPackagePath + "/properties-" + profile + ".yml"
	return configuration.GetConfig(path)

}
