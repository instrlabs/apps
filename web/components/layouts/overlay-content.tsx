"use client";

import React from "react";
import {useOverlay} from "@/hooks/useOverlay";

export default function OverlayContent({ children }: {
  children: React.ReactNode;
}) {
  const { isLeftOpen, isRightOpen } = useOverlay();

  const gridCols = isLeftOpen && isRightOpen
    ? "auto 1fr auto"
    : isLeftOpen ? "auto 1fr"
      : isRightOpen ? "1fr auto"
        : "1fr";

  return (
    <div className="flex-1 grid p-2 gap-4" style={{ gridTemplateColumns: gridCols }}>
      {children}
    </div>
  );
}
