"use client";

import React from "react";
import { ProfileProvider } from "@/hooks/useProfile";
import {NotificationProvider} from "@/components/notification";
import {OverlayProvider} from "@/hooks/useOverlay";

export default function Providers({ children }: { children: React.ReactNode }) {
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
