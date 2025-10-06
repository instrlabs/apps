import React from "react";
import { notFound } from "next/navigation";
import { SSEProvider } from "@/hooks/useSSE";
import { NotificationProvider, NotificationWidget } from "@/hooks/useNotification";
import { ProfileProvider } from "@/hooks/useProfile";
import { getProfile } from "@/services/auth";

export const dynamic = "force-dynamic";

export default async function DebugLayout({ children }: { children: React.ReactNode }) {
  if (process.env.NODE_ENV === "production") {
    notFound();
  }

  const res = await getProfile();

  return (
    <ProfileProvider data={res.data?.user || null}>
    <SSEProvider>
    <NotificationProvider>
      {children}
      <NotificationWidget />
    </NotificationProvider>
    </SSEProvider>
    </ProfileProvider>
  );
}
