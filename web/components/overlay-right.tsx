"use client";

import { useOverlay } from "@/hooks/useOverlay";

export default function OverlayRight() {
  const { isRightOpen } = useOverlay();
  return (
    <div
      className="absolute top-0 right-0 bottom-0 p-3 pt-[90px] h-screen transition-[width] duration-300 ease-in-out overflow-hidden"
      style={{ width: 'var(--overlay-right-width, 300px)', pointerEvents: isRightOpen ? 'auto' : 'none' }}
      aria-hidden={!isRightOpen}
      role="complementary"
      aria-label="Right overlay"
    >
      <div className="w-full h-full rounded-3xl bg-gray-200 flex flex-col justify-between">
        <div>Left</div>
        <div>Center</div>
        <div>Right</div>
      </div>
    </div>
  )
}
