import React from "react";
import { notFound } from "next/navigation";
import { SSEProvider } from "@/hooks/useSSE";
import { NotificationProvider, NotificationWidget } from "@/hooks/useNotification";

export const dynamic = "force-dynamic";

export default function DebugLayout({ children }: { children: React.ReactNode }) {
  if (process.env.NODE_ENV === "production") {
    notFound();
  }

  return (
    <SSEProvider>
    <NotificationProvider>
      {children}
      <NotificationWidget />
    </NotificationProvider>
    </SSEProvider>
  );
}
