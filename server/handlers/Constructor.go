package handlers

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	Client        *redis.Client
	Context       context.Context
	KAFKAUSERNAME string
	KAFKAPASSWORD string
	KAFKAADDRESS  string
}

// NewHandler initializes and returns a new Handler instance
func NewHandler(tableStore *redis.Client, user string, pass string,address string) *Handler {
	return &Handler{
		Client:        tableStore,
		Context:       context.Background(),
		KAFKAUSERNAME: user,
		KAFKAPASSWORD: pass,
		KAFKAADDRESS: address,
	}
}
