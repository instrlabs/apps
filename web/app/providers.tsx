"use client";

import React, { useEffect } from "react";

import { NotificationProvider } from "@/hooks/useNotification";
import { OverlayProvider, registerOverlays } from "@/hooks/useOverlay";

export default function Providers({ children }: { children: React.ReactNode }) {
  useEffect(() => {
    registerOverlays();
  }, []);

  return (
    <NotificationProvider>
      <OverlayProvider>
        {children}
      </OverlayProvider>
    </NotificationProvider>
  );
}
