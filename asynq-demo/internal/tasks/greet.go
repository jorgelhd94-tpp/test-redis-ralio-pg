package tasks

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hibiken/asynq"
)

// Payload de la tarea
type GreetPayload struct {
	Message string `json:"message"`
}

// Crear tarea con cola por defecto
func NewGreetTask(message string) (*asynq.Task, error) {
	payload := GreetPayload{Message: message}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	// Cola "default" implícita
	return asynq.NewTask("ralio:greet", data), nil
}

// Handler de la tarea
func HandleGreetTask(ctx context.Context, t *asynq.Task) error {
	var p GreetPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}

	log.Printf(p.Message)
	return nil
}
