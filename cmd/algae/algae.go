package main

import (
	"fmt"
	"tnraro/algae/internal/api"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	r := api.SetupRouter()
	r.Run(":41943")
}
