package main

import (
	"blackjackapi/server"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"net/http"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}
	ConnectString := os.Getenv("REDIS")
	opt, _ := redis.ParseURL(ConnectString)
	client := redis.NewClient(opt)
	// SET UP HANDLER DEPENDENCIES
	handler := server.NewHandler(client)
	Router := server.NewRouter(handler)
	http.ListenAndServe(":8080", Router)

}
