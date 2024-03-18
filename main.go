package main

import (
	"blackjackapi/server"
	"blackjackapi/server/handlers"
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
	handler := handlers.NewHandler(client, os.Getenv("UPSTASH_KAFKA_REST_USERNAME"), os.Getenv("UPSTASH_KAFKA_REST_PASSWORD"), os.Getenv("UPSTASH_KAFKA_REST_URL"))
	Router := server.NewRouter(handler)
	http.ListenAndServe(":8080", Router)

}
