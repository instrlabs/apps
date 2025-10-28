import React, { Suspense } from "react";
import { NotificationProvider, NotificationWidget } from "@/hooks/useSnackbar";

export default function LoginLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  return (
    <NotificationProvider>
      <Suspense>{children}</Suspense>
      <NotificationWidget />
    </NotificationProvider>
  );
}
