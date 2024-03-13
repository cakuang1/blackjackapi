package main

import (
	"blackjackapi/model"
	"blackjackapi/server"
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"os"
)

var ctx = context.Background()

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}
	ConnectString := os.Getenv("REDIS")
	opt, _ := redis.ParseURL(os.Getenv(ConnectString))
	client := redis.NewClient(opt)

	client.Set(ctx, "foo", "bar", 0)
	val := client.Get(ctx, "foo").Val()
	print(val)
}
