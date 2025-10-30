"use client";

import { useState, useContext, createContext, ReactNode, useEffect } from "react";

type InstructionNotification = {
  user_id: string;
  instruction_id: string;
  instruction_detail_id: string;
};

type SSEMessageEvent = {
  eventName: string;
  data: object | InstructionNotification;
};

interface SSEContextProps {
  message: SSEMessageEvent | null;
  isConnected: boolean;
  error: string | null;
}

const SSEContext = createContext<SSEContextProps | undefined>(undefined);

export const SSEProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [message, setMessage] = useState<SSEMessageEvent | null>(null);
  const [isConnected, setIsConnected] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const connect = async () => {
      try {
        setError(null);

        const res = await fetch("/api/sse", {
          cache: "no-store",
          credentials: "include",
        });

        if (!res.ok) {
          throw new Error(`Connection failed: ${res.status}`);
        }

        if (!res.body) {
          throw new Error("No response body");
        }

        setIsConnected(true);

        const reader = res.body.getReader();
        const decoder = new TextDecoder("utf-8");

        while (true) {
          const { done, value } = await reader.read();
          if (done) break;

          const text = decoder.decode(value, { stream: true });
          const lines = text.split(/\r?\n/).filter(Boolean);
          if (lines.length === 0) continue;

          const eventLine = lines.find((l) => l.startsWith("event:"));
          const dataLine = lines.find((l) => l.startsWith("data:"));

          const eventName = eventLine ? eventLine.slice(6).trimStart() : "message";
          const dataText = dataLine ? dataLine.slice(5).trimStart() : "{}";

          try {
            const data = JSON.parse(dataText);

            const isInstructionNotif = (obj: any): obj is InstructionNotification => {
              return (
                obj &&
                typeof obj.user_id === "string" &&
                typeof obj.instruction_id === "string" &&
                typeof obj.instruction_detail_id === "string"
              );
            };

            setMessage({
              eventName,
              data: isInstructionNotif(data) ? data : data,
            });
          } catch {
            continue;
          }
        }

        reader.releaseLock();
      } catch (err) {
        const msg = err instanceof Error ? err.message : "Unknown error";
        setError(msg);
      } finally {
        setIsConnected(false);
      }
    };

    connect();
  }, []);

  return (
    <SSEContext.Provider value={{ message, isConnected, error }}>
      {children}
    </SSEContext.Provider>
  );
};

const useSSE = (): SSEContextProps => {
  const context = useContext(SSEContext);
  if (!context) {
    throw new Error("useSSE must be used within SSEProvider");
  }
  return context;
};

export default useSSE;
