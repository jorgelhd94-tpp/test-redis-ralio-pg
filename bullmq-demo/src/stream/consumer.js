import { redis } from "../config/redis.js";
import { logInfo, logError } from "../utils/logger.js";
import { paymentGatewayQueue } from "../config/bullmq.js";

/**
 * Inicia un consumidor de Redis Stream con grupos de consumo (XREADGROUP)
 * @param {string} stream - Nombre del stream (ej: "payment-gateway-stream")
 * @param {string} group - Nombre del grupo (ej: "payment-gateway-group")
 * @param {string} consumerName - Nombre del consumidor (ej: "payment-gateway-consumer-1")
 */
export async function startConsumer(stream, group, consumerName) {
  logInfo(`👂 Worker "${consumerName}" escuchando en stream "${stream}" del grupo "${group}"`);

  // Crear grupo si no existe
  try {
    await redis.call("XGROUP", "CREATE", stream, group, "0", "MKSTREAM");
    logInfo(`✅ Grupo "${group}" creado para stream "${stream}"`);
  } catch (err) {
    if (err.message.includes("BUSYGROUP")) {
      logInfo(`ℹ️ Grupo "${group}" ya existe`);
    } else {
      logError(`❌ Error creando grupo "${group}": ${err.message}`);
    }
  }

  while (true) {
    try {
      const response = await redis.call(
        "XREADGROUP",
        "GROUP",
        group,
        consumerName,
        "BLOCK",
        5000,
        "COUNT",
        1,
        "STREAMS",
        stream,
        ">"
      );

      if (!response) continue;

      // La estructura que devuelve ioredis es:
      // [ [ streamName, [ [ id, [ field1, value1, field2, value2, ... ] ] ] ] ]
      for (const [streamName, messages] of response) {
        for (const [id, rawFields] of messages) {
          const fields = {};

          for (let i = 0; i < rawFields.length; i += 2) {
            fields[rawFields[i]] = rawFields[i + 1];
          }

          logInfo(`📦 [${id}] Campos recibidos:`, fields);

          let data;
          try {
            data = fields?.payload ? JSON.parse(fields.payload) : fields;
          } catch {
            data = { raw: fields.payload };
          }

          logInfo(`📩 [${id}] Mensaje procesado:`, data);

         // Enviar tarea a la cola
         logInfo(`Enviar tarea a la cola payment-gateway-queue: Task[${fields.task}], Payload:`, fields.payload);
         await paymentGatewayQueue.add(fields.task, fields.payload);

          // Confirmar mensaje
          await redis.call("XACK", streamName, group, id);
          logInfo(`✅ [${id}] Confirmado en grupo "${group}"`);
        }
      }
    } catch (err) {
      logError(`💥 Error en consumer "${consumerName}": ${err.message}`);
      await new Promise((r) => setTimeout(r, 2000));
    }
  }
}
