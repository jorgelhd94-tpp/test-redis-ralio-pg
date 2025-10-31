# Test Redis demos

Este repositorio reúne dos implementaciones que conectan Redis Streams con dos ecosistemas distintos:

1. **asynq-demo** (Go) procesa mensajes del stream y los transforma en tareas de Asynq para ejecutarlas en un worker especializado @asynq-demo/internal/stream/consumer.go#19-129 @asynq-demo/cmd/worker/worker.go#12-27.
2. **bullmq-demo** (Node.js) toma mensajes del stream y los encola en BullMQ para gestionarlos con workers basados en JavaScript @bullmq-demo/src/stream/consumer.js#12-78 @bullmq-demo/src/worker.js#1-20.

Ambos proyectos comparten Redis como capa de transporte y demuestran cómo coordinar servicios escritos en lenguajes diferentes utilizando el mismo stream.

## Flujo general

1. Un productor publica eventos en un Redis Stream con un campo `task` y un `payload` serializado en JSON @asynq-demo/internal/stream/producer.go#17-37 @bullmq-demo/src/stream/producer.js#10-27.
2. Cada consumidor crea (o asegura la existencia de) su `consumer group` antes de empezar a leer @asynq-demo/cmd/consumer/main.go#24-35 @bullmq-demo/src/stream/consumer.js#14-24.
3. Al recibir un mensaje, el consumidor envía el trabajo a su motor de tareas (Asynq o BullMQ) y confirma el mensaje en el stream para evitar reprocesos @asynq-demo/internal/stream/consumer.go#94-129 @bullmq-demo/src/stream/consumer.js#65-71.

## Requisitos

- Redis 7.x o superior ejecutándose localmente (puedes usar `docker run -p 6379:6379 redis:7`).
- Go 1.23.x para compilar `asynq-demo` @asynq-demo/go.mod#1-16.
- Node.js 18+ y npm para ejecutar `bullmq-demo` @bullmq-demo/package.json#1-16.

## Preparación

1. Clona este repositorio y sitúate en la raíz.
2. Asegúrate de tener Redis encendido en `127.0.0.1:6379`.

## asynq-demo (Go)

Estructura clave:

- `cmd/producer`: envía mensajes al stream compartido @asynq-demo/cmd/producer/main.go#13-27.
- `cmd/consumer`: lee del stream, crea tareas Asynq y las envía a la cola `ralio-queue` @asynq-demo/cmd/consumer/main.go#21-35 @asynq-demo/internal/stream/consumer.go#20-129.
- `cmd/worker`: procesa las tareas `ralio:greet` @asynq-demo/cmd/worker/worker.go#11-27 @asynq-demo/internal/tasks/greet.go#27-35.

Pasos recomendados:

```bash
cd asynq-demo
go mod tidy

# Terminal 1: worker
go run ./cmd/worker

# Terminal 2: consumer del stream
go run ./cmd/consumer

# Terminal 3: productor
go run ./cmd/producer
```

Observa cómo el consumer imprime el mensaje recibido y lo encola en Asynq antes de confirmar el ID en el stream @asynq-demo/internal/stream/consumer.go#95-129.

## bullmq-demo (Node.js)

Estructura clave:

- `src/stream/producer.js`: publica mensajes en el stream con una longitud máxima acotada @bullmq-demo/src/stream/producer.js#12-24.
- `src/stream/consumer.js`: consume mensajes mediante `XREADGROUP`, agrega jobs a BullMQ y confirma el stream @bullmq-demo/src/stream/consumer.js#26-78.
- `src/worker.js`: procesa los jobs de la cola `payment-gateway-queue` con un worker de BullMQ @bullmq-demo/src/worker.js#1-20.

Pasos recomendados:

```bash
cd bullmq-demo
npm install

# Terminal 1: worker de BullMQ
node src/worker.js

# Terminal 2: consumer del stream
node src/consumer.js

# Terminal 3: productor
node src/producer.js
```

La cola `payment-gateway-queue` se define con `Queue` de BullMQ y reutiliza una conexión `ioredis` con `maxRetriesPerRequest` deshabilitado, requisito habitual en BullMQ @bullmq-demo/src/config/bullmq.js#1-12.

## Ejecutar ambos flujos combinados

Puedes ejecutar los dos ecosistemas en paralelo para observar el enrutamiento según la `task`:

1. Inicia los workers de Asynq y BullMQ.
2. Lanza los consumidores de cada demo; cada uno permanecerá bloqueado esperando mensajes nuevos.
3. Envía mensajes desde cualquiera de los productores o escribe tus propios payloads reutilizando `SendToStream` o `produceToStream`.

Cada demo utiliza su propia cola de tareas (`ralio-queue` en Go, `payment-gateway-queue` en Node.js), lo que te permite comparar estrategias de ejecución y monitoreo basadas en Redis Streams sin interferir entre sí.

## Pasos siguientes sugeridos

1. Añadir validaciones personalizadas por `task` en los consumidores antes de encolar trabajos.
2. Incorporar métricas o dashboards (por ejemplo, `asynqmon` o `bull-board`) para observar el estado de las colas.
3. Automatizar el despliegue de Redis y de los procesos mediante Docker Compose para facilitar pruebas locales.
