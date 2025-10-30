import React from "react";
import { NotificationProvider, NotificationWidget } from "@/hooks/useSnackbar";

export default function LoginLayout({ children }: { children: React.ReactNode }) {
  return (
    <NotificationProvider>
      {children}
      <NotificationWidget />
    </NotificationProvider>
  );
}
