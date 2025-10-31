import { redis } from "../config/redis.js";
import { logInfo, logError } from "../utils/logger.js";

export async function initGroup(stream, group) {
  try {
    await redis.xgroup("CREATE", stream, group, "0", "MKSTREAM");
    logInfo(`Grupo "${group}" creado`);
  } catch (err) {
    if (err.message.includes("BUSYGROUP")) {
      logInfo(`Grupo "${group}" ya existe`);
    } else {
      logError(`Error creando grupo: ${err.message}`);
      throw err;
    }
  }
}
