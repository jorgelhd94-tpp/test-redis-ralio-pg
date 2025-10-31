import { produceToStream } from "./stream/producer.js";

(async () => {
  await produceToStream("ralio-stream", "ralio:greet", {
    message: "Hello! This message comes from Payment Gateway in Node.js 🚀"
  });
  process.exit(0);
})();
