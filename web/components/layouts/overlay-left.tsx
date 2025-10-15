"use client";

import { useOverlay } from "@/hooks/useOverlay";

export default function OverlayLeft() {
  const { leftNode, leftKey } = useOverlay();

  return (
    <div
      className={
        "absolute top-0 left-0 bottom-0 " +
        "pt-[80px] h-screen"
      }
    >
      <div
        key={leftKey}
        className="flex-1 animate-fade-in"
      >
        {leftNode}
      </div>
    </div>
  )
}
