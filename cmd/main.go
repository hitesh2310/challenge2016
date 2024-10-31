package main

import (
	"fmt"
	"main/config"
	"main/pkg/router"
)

func init() {
	config.SetUpApplication()

}

func main() {
	fmt.Println("Hello World")
	server := router.NewServer()
	server.Run(":8090")
}
