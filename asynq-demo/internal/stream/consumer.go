package stream

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
)

// Payload define la estructura de tus mensajes
type Payload struct {
	Message string `json:"message"`
}

// ConsumeStream procesa mensajes pendientes y nuevos del stream
func ConsumeStream(ctx context.Context, rdb *redis.Client, asynqClient *asynq.Client, stream, group, consumer string) {
	// 1️⃣ Procesar mensajes pendientes primero
	processPending(ctx, rdb, asynqClient, stream, group, consumer)

	// 2️⃣ Leer nuevos mensajes continuamente
	for {
		res, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    group,
			Consumer: consumer,
			Streams:  []string{stream, ">"},
			Count:    10,
			Block:    0, // Bloquea máximo 2s esperando nuevos mensajes
		}).Result()

		if err != nil {
			if err != redis.Nil {
				// Solo log de errores reales
				log.Println("Error leyendo del stream:", err)
			}
			// Si es timeout o sin mensajes, dormir un poco y continuar
			time.Sleep(100 * time.Millisecond)
			continue
		}

		if res == nil {
			// No hay mensajes nuevos
			continue
		}

		for _, s := range res {
			for _, msg := range s.Messages {
				processMessage(ctx, rdb, asynqClient, stream, group, msg)
			}
		}
	}
}

// processPending revisa mensajes pendientes y los reclama
func processPending(ctx context.Context, rdb *redis.Client, asynqClient *asynq.Client, stream, group, consumer string) {
	pending, err := rdb.XPendingExt(ctx, &redis.XPendingExtArgs{
		Stream: stream,
		Group:  group,
		Start:  "-",
		End:    "+",
		Count:  100,
	}).Result()

	if err != nil {
		if err != redis.Nil {
			log.Println("Error obteniendo pendientes:", err)
		}
		return
	}

	for _, item := range pending {
		msgs, err := rdb.XClaim(ctx, &redis.XClaimArgs{
			Stream:   stream,
			Group:    group,
			Consumer: consumer,
			MinIdle:  0,
			Messages: []string{item.ID},
		}).Result()

		if err != nil {
			log.Println("Error reclamando mensaje:", err)
			continue
		}

		for _, msg := range msgs {
			processMessage(ctx, rdb, asynqClient, stream, group, msg)
		}
	}
}

// processMessage parsea y confirma el mensaje
func processMessage(ctx context.Context, rdb *redis.Client, asynqClient *asynq.Client, stream, group string, msg redis.XMessage) {
	var taskName string
	if value, ok := msg.Values["task"].(string); ok {
		taskName = value
	} else {
		log.Println("Error parseando task. Stream: ", stream)
		return
	}

	var payload Payload
	if p, ok := msg.Values["payload"].(string); ok {
		if err := json.Unmarshal([]byte(p), &payload); err != nil {
			log.Println("Error parseando payload:", err)
			return
		}
	}

	fmt.Printf("📩 [ID %s] [TASK %s] Mensaje recibido: %+v\n", msg.ID, taskName, payload)

	// Simula trabajo
	taskPayload, _ := json.Marshal(payload)
	task := asynq.NewTask(taskName, taskPayload)
	if _, err := asynqClient.Enqueue(task, asynq.Queue("ralio-queue")); err != nil {
		log.Println("Error encolando tarea en Asynq:", err)
		return
	}

	fmt.Printf("➡️ [ID %s] Tarea encolada en Asynq: %s\n", msg.ID, taskName)

	// Confirmar mensaje procesado
	if _, err := rdb.XAck(ctx, stream, group, msg.ID).Result(); err != nil {
		log.Println("Error confirmando mensaje:", err)
	} else {
		fmt.Printf("✅ [ID %s] Confirmado\n", msg.ID)
	}
}
