import { initGroup } from "./stream/group.js";
import { startConsumer } from "./stream/consumer.js";

const stream = "payment-gateway-stream";
const group = "payment-gateway-group";
const consumerName = "payment-gateway-consumer-1";

async function main() {
  await initGroup(stream, group);
  await startConsumer(stream, group, consumerName);
}

main();
