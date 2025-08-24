"use client";

import React, { useEffect } from "react";

import { ProfileProvider } from "@/hooks/useProfile";
import { NotificationProvider } from "@/components/notification";
import { OverlayProvider } from "@/hooks/useOverlay";
import {registerBuiltInOverlays, resetOverlayRegistry} from "@/hooks/overlayRegistry";

export default function Providers({ children }: { children: React.ReactNode }) {
  useEffect(() => {
    registerBuiltInOverlays();
    return () => resetOverlayRegistry();
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
