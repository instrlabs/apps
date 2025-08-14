"use client";

import { useOverlay } from "@/hooks/useOverlay";

export default function OverlayRight() {
  const { isRightOpen, rightNode, rightTitle, rightContentKey } = useOverlay();
  return (
    <div
      className="absolute top-0 right-0 bottom-0 p-3 pt-[80px] h-screen transition-[width] duration-300 ease-in-out"
      style={{ width: 'var(--overlay-right-width, 300px)', pointerEvents: isRightOpen ? 'auto' : 'none' }}
      aria-hidden={!isRightOpen}
      role="complementary"
      aria-label="Right overlay"
    >
      <div className="w-full h-full rounded-3xl bg-neutral-50 flex flex-col overflow-hidden">
        {rightTitle ? (
          <div className="px-4 py-3 border-b text-sm font-semibold text-gray-700">{rightTitle}</div>
        ) : null}
        <div key={rightContentKey} className="flex-1 overflow-auto p-4 animate-fade-in">
          {rightNode}
        </div>
      </div>
    </div>
  )
}
