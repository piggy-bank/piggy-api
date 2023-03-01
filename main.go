package main

import (
	"github.com/manubidegain/piggy-api/cmd/api"
)

func main() {
	app := &api.App{}
	app.Initialize()
}
