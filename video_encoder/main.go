package main

import (
	"config"
	"fmt"
	"router"
)

func main() {
	config.API.Router = router.InitializeRouter()
	fmt.Println("App running at 127.0.0.1:3001 ...")
	config.API.Run(":3001")
}
