package main

import (
	"log"

	// "asynq-demo/internal/tasks"

	// "github.com/hibiken/asynq"
	"asynq-demo/internal/stream"
	"fmt"
)

func main() {
	rdb := stream.NewRedisClient()

	payload := map[string]any{
		"message": "Hello! This message come from Ralio in Go 🚀",
	}

	msgID, err := stream.SendToStream(rdb, "payment-gateway-stream", payload)
	if err != nil {
		log.Fatalf("❌ Error enviando al stream: %v", err)
	}

	fmt.Println("✅ Mensaje enviado correctamente")
	fmt.Println("🆔 ID:", msgID)
}
