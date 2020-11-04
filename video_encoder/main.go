package main

import (
	"config"
	"fmt"
	"router"
)

func main() {
	config.Api.Router = router.InitializeRouter()
	fmt.Println("App running at 127.0.0.1:3001 ...")
	config.Api.Run(":3001")
	config.Api.Db.Close()
}
