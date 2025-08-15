"use client";

import { useOverlay } from "@/hooks/useOverlay";
import { useMemo } from "react";
import clsx from "clsx";

export default function OverlayRight() {
  const { isRightOpen, rightNode, rightContentKey, rightWidth } = useOverlay();
  const targetWidth = Number.isFinite(rightWidth) ? Math.max(0, Math.round(rightWidth)) : 0;
  const widthPx = useMemo(() => isRightOpen ? targetWidth : 0, [isRightOpen, targetWidth]);

  return (
    <div
      className="absolute top-0 right-0 bottom-0 p-3 pt-[80px] h-screen transition-[width] duration-300 ease-in-out"
      style={{ width: `${widthPx}px`, pointerEvents: isRightOpen ? 'auto' : 'none' }}
      aria-hidden={!isRightOpen}
      role="complementary"
      aria-label="Right overlay"
    >
      <div className={clsx(
        "w-full h-full flex flex-col",
        isRightOpen ? "overflow-visible" : "overflow-hidden",
      )}>
        <div key={rightContentKey} className="flex-1 animate-fade-in">
          {rightNode}
        </div>
      </div>
    </div>
  )
}
