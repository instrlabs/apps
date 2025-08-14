"use client";

import { useOverlay } from "@/hooks/useOverlay";

export default function OverlayLeft() {
  const { isLeftOpen, leftNode, leftContentKey } = useOverlay();

  return (
    <div
      className="absolute top-0 left-0 bottom-0 p-3 pt-[80px] h-screen transition-[width] duration-300 ease-in-out"
      style={{ width: 'var(--overlay-left-width, 300px)', pointerEvents: isLeftOpen ? 'auto' : 'none' }}
      aria-hidden={!isLeftOpen}
      role="complementary"
      aria-label="Left overlay"
    >
      <div className="w-full h-full rounded-3xl bg-neutral-50 flex flex-col overflow-hidden">
        <div key={leftContentKey} className="flex-1 overflow-auto p-4 animate-fade-in">
          {leftNode}
        </div>
      </div>
    </div>
  )
}
