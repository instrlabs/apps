"use client";

import React, { useEffect, useMemo, useRef, useState } from "react";
import useSSE from "@/hooks/useSSE";

export default function DebugSSEPage() {
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
      className={`
      w-screen h-screen bg-gray-950 text-gray-200
      font-mono p-8 overflow-auto whitespace-pre-wrap`}
    >
      {content}
    </div>
  );
}
