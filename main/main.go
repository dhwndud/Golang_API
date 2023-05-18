package main

import (
	"chatbot/mappings"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	mappings.CreateURLMappings()
	//Listen and server on 0.0.0.0:0727
	mappings.Router.Run(":0727")
}
