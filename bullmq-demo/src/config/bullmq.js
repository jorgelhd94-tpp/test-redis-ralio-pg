import { Queue } from "bullmq";
import IORedis from "ioredis";

export const bullConnection = new IORedis({
  host: "127.0.0.1",
  port: 6379,
  maxRetriesPerRequest: null, // Obligatorio para BullMQ
});

export const paymentGatewayQueue = new Queue("payment-gateway-queue", {
  connection: bullConnection,
});
