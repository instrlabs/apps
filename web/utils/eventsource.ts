import {fetchEventStream} from "@/utils/fetch";

export type SSEMessageEvent = {
  eventName: string;
  data: string;
};

export class ServerEventSource {
  url: string;
  onmessage: (ev: SSEMessageEvent) => void;

  constructor(url: string) {
    this.url = url;
    this.onmessage = () => {};

    this.start().finally();
  }

  private async start() {
    const res = await fetchEventStream("/sse");

    if (!res.ok) {
      console.warn("SSE connection failed:", res.status, res.statusText);
      return;
    }

    if (!res.body) {
      console.warn("SSE response has no body");
      return;
    }

    const reader = res.body.getReader();
    const decoder = new TextDecoder("utf-8");

    try {
      while (true) {
        const { done, value } = await reader.read();
        if (done) break;

        const text = decoder.decode(value, { stream: true });
        const lines = text.split(/\r?\n/);
        const eventName = lines[0].slice(6).trimStart();
        const dataText = lines[1].slice(5).trimStart();
        const data = JSON.parse(dataText);
        this.onmessage({ eventName, data });
      }

    } catch (err) {
      console.warn("SSE connection error:", err);
    } finally {
      reader.releaseLock();
    }
  }
}

export function EventSource(url: string) {
  return new ServerEventSource(url);
}

export default EventSource;
