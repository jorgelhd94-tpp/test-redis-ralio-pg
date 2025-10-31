import { Worker } from "bullmq";
import { bullConnection } from "./config/bullmq.js"; // <-- usar esta
import { logInfo, logError } from "./utils/logger.js";

const worker = new Worker(
  "payment-gateway-queue",
  async (job) => {
    logInfo(`👷 Procesando job BullMQ [${job.id}]. Task [${job.name}]`, job.data);

    // Simula trabajo
    await new Promise((resolve) => setTimeout(resolve, 1000));

    logInfo(`✅ Job BullMQ [${job.id}] completado`);
  },
  { connection: bullConnection } // <-- aquí sí
);

worker.on("failed", (job, err) => {
  logError(`❌ Job BullMQ [${job.id}] falló:`, err);
});
