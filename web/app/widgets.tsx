"use client";

import React from "react";
import OverlayLeft from "@/components/overlay-left";
import OverlayRight from "@/components/overlay-right";
import OverlayTop from "@/components/overlay-top";
import OverlayModal from "@/components/overlay-modal";
import { NotificationWidget } from "@/hooks/useNotification";

export default function Widgets() {
  return (
    <>
      <OverlayLeft />
      <OverlayRight />
      <OverlayModal />
      <OverlayTop />
      <NotificationWidget />
    </>
  );
}
