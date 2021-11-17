package main

import (
	"config"
	"fmt"
	"router"
)

func main() {
	config.API.Router = router.InitializeRouter()
	config.API.Run(":3001")
}
