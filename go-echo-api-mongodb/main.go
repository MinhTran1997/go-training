package main

import (
	"fmt"
	"go-echo-api-mongodb/router"
)

// The main.go is our server. It will start a server on xxxx port and serve all the Router.

func main() {
	e := router.Router()

	fmt.Println("Starting server on the port 1325...")

	e.Logger.Fatal(e.Start(":1325"))
}
