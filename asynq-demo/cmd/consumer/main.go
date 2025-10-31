package main

import (
	"context"
	"log"

	"asynq-demo/internal/stream"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()

	// Conexión a Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})

	asynqClient := asynq.NewClient(asynq.RedisClientOpt{Addr: "127.0.0.1:6379"})
	defer asynqClient.Close()

	streamName := "ralio-stream"
	groupName := "ralio-consumer-group"
	consumerName := "ralio-consumer-1"

	// Crear grupo si no existe (usar "0" para procesar mensajes viejos la primera vez)
	err := rdb.XGroupCreateMkStream(ctx, streamName, groupName, "0").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		log.Fatal("Error creando grupo:", err)
	}

	log.Println("👂 Consumer Go escuchando stream:", streamName)
	stream.ConsumeStream(ctx, rdb, asynqClient, streamName, groupName, consumerName)
}
