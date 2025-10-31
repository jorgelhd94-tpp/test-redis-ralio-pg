package stream

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// RedisClient interfaz mínima para permitir test/mocks
type RedisClient interface {
	XAdd(ctx context.Context, a *redis.XAddArgs) *redis.StringCmd
}

// SendToStream serializa cualquier payload y lo envía al stream dado
func SendToStream(rdb RedisClient, streamName string, payload any) (string, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("error serializando payload: %w", err)
	}

	msgID, err := rdb.XAdd(Ctx, &redis.XAddArgs{
		Stream: streamName,
		// MaxLen: 1000,// Se puede limitar la cantidad de mensajes en el stream
		Values: map[string]any{
			"task":    "payment-gateway:greet",
			"payload": string(data),
		},
	}).Result()

	if err != nil {
		return "", fmt.Errorf("error enviando mensaje al stream: %w", err)
	}

	return msgID, nil
}
