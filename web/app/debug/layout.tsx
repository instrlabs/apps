import React from "react";
import { notFound } from "next/navigation";
import { SSEProvider } from "@/hooks/useSSE";

export const dynamic = "force-dynamic";

export default function DebugLayout({ children }: { children: React.ReactNode }) {
  // Prevent access to all /debug pages in production
  if (process.env.NODE_ENV === "production") {
    notFound();
  }

  return (
    <SSEProvider>
      {children}
    </SSEProvider>
  );
}
