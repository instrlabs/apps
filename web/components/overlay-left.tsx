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
      className="absolute top-0 left-0 bottom-0 p-3 pt-[80px] h-screen transition-[width] duration-300 ease-in-out"
      style={{ width: `${widthPx}px`, pointerEvents: isLeftOpen ? 'auto' : 'none' }}
      aria-hidden={!isLeftOpen}
      role="complementary"
      aria-label="Left overlay"
    >
      <div className={clsx(
        "w-full h-full flex flex-col",
        isLeftOpen ? "overflow-visible" : "overflow-hidden",
      )}>
        <div key={leftKey} className="flex-1 animate-fade-in">
          {leftNode}
        </div>
      </div>
    </div>
  )
}
