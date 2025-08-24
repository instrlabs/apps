"use client";

import React, { useEffect } from "react";

import { ProfileProvider } from "@/hooks/useProfile";
import { NotificationProvider } from "@/components/notification";
import { OverlayProvider, registerOverlays } from "@/hooks/useOverlay";

export default function Providers({ children }: { children: React.ReactNode }) {
  useEffect(() => {
    registerOverlays();
  }, []);

  return (
    <ProfileProvider>
      <NotificationProvider>
        <OverlayProvider>
          {children}
        </OverlayProvider>
      </NotificationProvider>
    </ProfileProvider>
  );
}
