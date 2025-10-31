package main

import (
	"log"

	"asynq-demo/internal/tasks"

	"github.com/hibiken/asynq"
)

func main() {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: "localhost:6379"},
		asynq.Config{
			Concurrency: 5,
			Queues: map[string]int{
				"ralio-queue": 1,
			},
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc("ralio:greet", tasks.HandleGreetTask)

	log.Println("👷 Worker ejecutándose en cola por defecto...")
	if err := srv.Run(mux); err != nil {
		log.Fatalf("Error iniciando worker: %v", err)
	}
}
