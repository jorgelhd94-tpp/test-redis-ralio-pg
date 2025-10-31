export function logInfo(...args) {
  console.log(new Date().toISOString(), "ℹ️", ...args);
}

export function logError(...args) {
  console.error(new Date().toISOString(), "❌", ...args);
}
