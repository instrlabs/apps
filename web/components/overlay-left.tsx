"use client";

import { useOverlay } from "@/hooks/useOverlay";

export default function OverlayLeft() {
  const { isLeftOpen } = useOverlay();
  return (
    <div
      className="absolute top-0 left-0 bottom-0 p-3 pt-[90px] h-screen transition-[width] duration-300 ease-in-out overflow-hidden"
      style={{ width: 'var(--overlay-left-width, 300px)', pointerEvents: isLeftOpen ? 'auto' : 'none' }}
      aria-hidden={!isLeftOpen}
      role="complementary"
      aria-label="Left overlay"
    >
      <div className="w-full h-full rounded-3xl bg-gray-200 flex flex-col justify-between">
        <div>Left</div>
        <div>Center</div>
        <div>Right</div>
      </div>
    </div>
  )
}
