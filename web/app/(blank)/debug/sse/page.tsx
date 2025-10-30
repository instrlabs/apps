"use client";

import { SSEProvider } from "@/hooks/useSSE";
import { SSEConsoleDisplay } from "./ConsoleSSE";

export default function DebugSSEPage() {
  return (
    <SSEProvider>
      <SSEConsoleDisplay />
    </SSEProvider>
  );
}
