"use client";

import React, { useState, useContext, createContext, ReactNode, useEffect } from "react";

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
}

const SSEContext = createContext<SSEContextProps | undefined>(undefined);

export const SSEProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [message, setMessage] = useState<SSEMessageEvent | null>(null);

  useEffect(() => {
    async function start() {
      const url = "/api/sse";
      const res = await fetch(url, {
        cache: "no-store",
      });

      if (!res.ok) {
        console.warn("[useSSE] Connection failed:", res.status, res.statusText);
        return;
      }

      if (!res.body) {
        console.warn("[useSSE] Response has no body");
        return;
      }

      console.log("[useSSE] Connection established, reading stream...");

      const reader = res.body.getReader();
      const decoder = new TextDecoder("utf-8");

      try {
        while (true) {
          const { done, value } = await reader.read();
          if (done) break;

          const text = decoder.decode(value, { stream: true });
          const lines = text.split(/\r?\n/).filter(Boolean);
          if (lines.length === 0) continue;
          // Basic SSE parsing for lines like: "event: <name>" and "data: <json>"
          const eventLine = lines.find((l) => l.startsWith("event:"));
          const dataLine = lines.find((l) => l.startsWith("data:"));
          const eventName = eventLine ? eventLine.slice(6).trimStart() : "message";
          const dataText = dataLine ? dataLine.slice(5).trimStart() : "{}";
          try {
            const data = JSON.parse(dataText);
            console.log("[useSSE] Received message:", { eventName, data });

            // Type guard for instruction notifications
            const isInstructionNotification = (parsed: any): parsed is InstructionNotification => {
              return (
                parsed &&
                typeof parsed.user_id === "string" &&
                typeof parsed.instruction_id === "string" &&
                typeof parsed.instruction_detail_id === "string"
              );
            };

            // Set message with proper typing
            setMessage({
              eventName,
              data: isInstructionNotification(data) ? data : data,
            });
          } catch (err) {
            console.warn("[useSSE] Failed to parse data:", dataText, err);
            // ignore non-JSON data lines
          }
        }
      } catch (err) {
        console.warn("SSE connection error:", err);
      } finally {
        reader.releaseLock();
      }
    }

    start().then();
  }, []);

  return <SSEContext.Provider value={{ message }}>{children}</SSEContext.Provider>;
};

const useSSE = (): SSEContextProps => {
  const context = useContext(SSEContext);

  if (context === undefined) {
    throw new Error("useSSE must be used within a SSEProvider");
  }

  return context;
};

export default useSSE;
