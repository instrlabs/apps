"use client";

import { useOverlay } from "@/hooks/useOverlay";

export default function OverlayLeft() {
  const { isLeftOpen, leftNode, leftContentKey, leftWidth } = useOverlay();

  const widthPx = Number.isFinite(leftWidth) ? Math.max(0, Math.round(leftWidth)) : 0;
  return (
    <div
      className="absolute top-0 left-0 bottom-0 p-3 pt-[80px] h-screen transition-[width] duration-300 ease-in-out"
      style={{ width: `${widthPx}px`, pointerEvents: isLeftOpen ? 'auto' : 'none' }}
      aria-hidden={!isLeftOpen}
      role="complementary"
      aria-label="Left overlay"
    >
      <div className="w-full h-full rounded-3xl flex flex-col overflow-visible">
        <div key={leftContentKey} className="flex-1 overflow-y-auto overflow-x-visible animate-fade-in">
          {leftNode}
        </div>
      </div>
    </div>
  )
}
