package main

import (
	"api"
	"fmt"
	"net/http"

	"github.com/go-flow/flow/v2"

	_ "api/docs"
)

// @title Core API
// @version 0.1.0
// @description REST API for Core API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /
// @query.collection.format multi
func main() {

	app, err := flow.Bootstrap(new(api.AppModule))

	if err != nil {
		fmt.Printf("unable to bootstrap app. Error: %s \n", err.Error())
	}

	if err := app.Serve(); err != nil && err != http.ErrServerClosed {
		fmt.Printf("error serving App. Error: %v", err)
	}

}
