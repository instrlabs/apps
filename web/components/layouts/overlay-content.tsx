"use client";

import React from "react";
import {useOverlay} from "@/hooks/useOverlay";
import OverlayLeft from "@/components/layouts/overlay-left";
import OverlayBody from "@/components/layouts/overlay-body";
import OverlayRight from "@/components/layouts/overlay-right";
import { useMobile } from "@/hooks/useMediaQuery";

export default function OverlayContent({ children }: {
  children: React.ReactNode;
}) {
  const { isLeftOpen, isRightOpen } = useOverlay();
  const isMobile = useMobile();

  if (isMobile) {
    return (
      <div className="flex-1 flex flex-col gap-4 p-2">
        {
          isLeftOpen ? <OverlayLeft /> :
          isRightOpen ? <OverlayRight /> :
          <OverlayBody>{children}</OverlayBody>
        }
      </div>
    );
  }

  const gridCols = isLeftOpen && isRightOpen
    ? "auto 1fr auto"
    : isLeftOpen ? "auto 1fr"
      : isRightOpen ? "1fr auto"
        : "1fr";

  return (
    <div className="flex-1 grid p-2 gap-4" style={{ gridTemplateColumns: gridCols }}>
      <OverlayLeft />
      <OverlayBody>
        {children}
      </OverlayBody>
      <OverlayRight />
    </div>
  );
}
