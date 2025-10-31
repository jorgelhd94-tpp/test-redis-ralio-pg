import { redis } from "../config/redis.js";
import { logInfo, logError } from "../utils/logger.js";

/**
 * Envia un payload a un stream Redis
 * @param {string} stream - Nombre del stream
 * @param {string} task - Tipo de tarea (ej: "email:send")
 * @param {object} payload - Payload JSON que quieres enviar
 */
export async function produceToStream(stream, task, payload) {
  try {
    const id = await redis.xadd(
      stream,
      "MAXLEN",
      "~",
      1000, // máximo 1000 mensajes en el stream
      "*", // id automático
      "task",
      task,
      "payload",
      JSON.stringify(payload)
    );

    logInfo(`📤 Mensaje enviado al stream "${stream}" con ID ${id}`);
    return id;
  } catch (err) {
    logError(`❌ Error enviando mensaje al stream "${stream}":`, err);
    throw err;
  }
}
