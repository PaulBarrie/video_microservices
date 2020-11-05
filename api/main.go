package main

import (
	"config"
	"fmt"
	"router"
)

// @title Youtube API
// @version 1.0
// @description This is a sample service for managing orders
// @termsOfService http://swagger.io/terms/
// @contact.name barrie_p
// @contact.email barrie_p@etna-alternance.net
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:3000
// @BasePath
func main() {
	config.Api.Router = router.InitializeRouter()
	fmt.Println("App running at 127.0.0.1:3000...")
	config.Api.Run(":3000")
	config.Api.Db.Close()
}
