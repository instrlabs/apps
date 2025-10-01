"use client";

import React, {
  useState,
  useContext,
  createContext,
  ReactNode, useEffect,
} from "react";

type SSEMessageEvent = {
  eventName: string;
  data: object;
}

interface SSEContextProps {
  message: SSEMessageEvent | null;
}

const SSEContext = createContext<SSEContextProps | undefined>(undefined);

export const SSEProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [message, setMessage] = useState<SSEMessageEvent | null>(null);

  useEffect(() => {
    async function start() {
      const url = process.env.NOTIFICATION_URL + "/sse";
      const res = await fetch(url, { credentials: "include" });

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
          setMessage({ eventName, data });
        }

      } catch (err) {
        console.warn("SSE connection error:", err);
      } finally {
        reader.releaseLock();
      }
    }

    start().then()
  }, []);

  return (
    <SSEContext.Provider value={{ message }}>
      {children}
    </SSEContext.Provider>
  );
};

const useSSE = (): SSEContextProps => {
  const context = useContext(SSEContext);

  if (context === undefined) {
    throw new Error("useSSE must be used within a SSEProvider");
  }

  return context;
};

export default useSSE;
