import React from "react";
import { NotificationProvider, Notification } from "@/components/notification";
import OverlayTop from "@/components/overlay-top";
import OverlayLeft from "@/components/overlay-left";
import OverlayRight from "@/components/overlay-right";
import OverlayContent from "@/components/overlay-content";
import OverlayModal from "@/components/overlay-modal";
import { OverlayProvider } from "@/hooks/useOverlay";

export default function SiteLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <NotificationProvider>
      <OverlayProvider>
        <OverlayContent>
          {children}
        </OverlayContent>
        <OverlayLeft />
        <OverlayRight />
        <OverlayTop />
        <OverlayModal />
        <Notification />
      </OverlayProvider>
    </NotificationProvider>
  );
}
