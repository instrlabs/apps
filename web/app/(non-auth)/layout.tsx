import React from "react";
import { NotificationProvider, NotificationWidget } from "@/hooks/useNotification";
import NonAuthGuard from "@/components/non-auth-guard";

export default function LoginLayout({
  children,
}: Readonly<{ children: React.ReactNode }>) {
  return (
    <NotificationProvider>
      <NonAuthGuard />
      {children}
      <NotificationWidget />
    </NotificationProvider>
  );
}
