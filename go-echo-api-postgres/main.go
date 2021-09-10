package main

import (
	"fmt"
	"go-echo-api-postgres/router"
)

// The main.go is our server. It will start a server on xxxx port (declared bellow in func ListenAndServe) and serve all the Router.

func main() {
	e := router.Router()

	fmt.Println("Starting server on the port 1324...")

	e.Logger.Fatal(e.Start(":1324"))
}
