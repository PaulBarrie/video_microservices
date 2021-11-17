package main

import (
	"fmt"
	"os"
	"config"
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
	config.API.Router = router.InitializeRouter()
	fmt.Println("App running on port 3000...")
	config.API.Run(fmt.Sprintf(":%s", os.Getenv("API_PORT")))
	config.API.Db.Close()
}
