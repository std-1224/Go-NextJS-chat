export const API_URL = process.env.API_URL || "http://127.0.0.1:8080"
export const WEBSOCKET_URL =
  process.env.NEXT_PUBLIC_WEBSOCKET_URL || "ws://127.0.0.1:8080"
export const SOCKET_ACTION = {
  LEFT: 0,
  JOIN: 1,
  SEND: 2,
  READ: 3,
} as const
