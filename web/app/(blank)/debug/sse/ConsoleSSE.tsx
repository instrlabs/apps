"use client";

import { useEffect, useRef, useState, useMemo } from "react";
import useSSE from "@/hooks/useSSE";

export function SSEConsoleDisplay() {
  const { message } = useSSE();
  const [logs, setLogs] = useState<string[]>([]);
  const boxRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!message) return;
    const line = `[${new Date().toISOString()}] ${message.eventName}: ${JSON.stringify(message.data)}`;
    setLogs((prev) => [...prev, line]);
  }, [message]);

  useEffect(() => {
    const el = boxRef.current;
    if (!el) return;
    el.scrollTop = el.scrollHeight;
  }, [logs]);

  const content = useMemo(() => logs.join("\n"), [logs]);

  return (
    <div
      ref={boxRef}
      className="h-screen w-screen overflow-auto bg-gray-950 p-8 font-mono whitespace-pre-wrap text-gray-200"
    >
      {content}
    </div>
  );
}
