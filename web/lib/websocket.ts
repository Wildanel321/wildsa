import { getAuthToken } from './api';

type MessageHandler = (data: any) => void;

export class SawitWebSocket {
  private ws: WebSocket | null = null;
  private handlers: Map<string, Set<MessageHandler>> = new Map();
  private reconnectInterval: number = 3000;
  private isConnecting: boolean = false;

  constructor() {
    this.connect();
  }

  private connect() {
    if (typeof window === 'undefined' || this.isConnecting) return;
    this.isConnecting = true;

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const host = window.location.host;
    const token = getAuthToken() || '';
    const wsUrl = `${protocol}//${host}/api/v1/ws?token=${encodeURIComponent(token)}`;

    try {
      this.ws = new WebSocket(wsUrl);

      this.ws.onopen = () => {
        this.isConnecting = false;
      };

      this.ws.onmessage = (event) => {
        try {
          const payload = JSON.parse(event.data);
          if (payload.topic && this.handlers.has(payload.topic)) {
            const topicHandlers = this.handlers.get(payload.topic);
            topicHandlers?.forEach((handler) => handler(payload.data));
          }
        } catch (e) {
          // Ignore invalid JSON frames
        }
      };

      this.ws.onclose = () => {
        this.isConnecting = false;
        setTimeout(() => this.connect(), this.reconnectInterval);
      };

      this.ws.onerror = () => {
        this.ws?.close();
      };
    } catch (e) {
      this.isConnecting = false;
      setTimeout(() => this.connect(), this.reconnectInterval);
    }
  }

  public subscribe(topic: string, handler: MessageHandler): () => void {
    if (!this.handlers.has(topic)) {
      this.handlers.set(topic, new Set());
    }
    this.handlers.get(topic)!.add(handler);

    return () => {
      this.handlers.get(topic)?.delete(handler);
    };
  }
}

let wsInstance: SawitWebSocket | null = null;

export function getWebSocketInstance(): SawitWebSocket {
  if (!wsInstance && typeof window !== 'undefined') {
    wsInstance = new SawitWebSocket();
  }
  return wsInstance!;
}
