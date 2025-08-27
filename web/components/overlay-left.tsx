"use client";

import { useOverlay } from "@/hooks/useOverlay";
import { useMemo } from "react";
import clsx from "clsx";

export default function OverlayLeft() {
  const { isLeftOpen, leftNode, leftKey, leftWidth } = useOverlay();
  const targetWidth = Math.max(0, leftWidth);
  const widthPx = useMemo(() => (isLeftOpen ? targetWidth : 0), [isLeftOpen, targetWidth]);

  return (
    <div
      className={clsx(
        "absolute top-0 left-0 bottom-0",
        "pt-[80px] h-screen",
        "transition-[width] duration-300 ease-in-out"
      )}
      style={{
        width: `${widthPx}px`,
        pointerEvents: isLeftOpen ? 'auto' : 'none'
      }}
    >
      <div
        key={leftKey}
        className={clsx(
          "flex-1 animate-fade-in",
          isLeftOpen ? "overflow-visible" : "overflow-hidden",
        )}
      >
        {leftNode}
      </div>
    </div>
  )
}
