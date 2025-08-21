import React from "react";
import { NotificationProvider, Notification } from "@/components/notification";
import NonAuthGuard from "@/components/non-auth-guard";

export default function LoginLayout({
  children,
}: Readonly<{ children: React.ReactNode }>) {
  return (
    <NotificationProvider>
      <NonAuthGuard />
      {children}
      <Notification />
    </NotificationProvider>
  );
}
