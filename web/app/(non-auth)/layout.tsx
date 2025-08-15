import React from "react";
import { NotificationProvider, Notification } from "@/components/notification";

export default function LoginLayout({
  children,
}: Readonly<{ children: React.ReactNode }>) {
  return (
    <NotificationProvider>
      {children}
      <Notification />
    </NotificationProvider>
  );
}
