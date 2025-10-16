"use client";

import { useOverlay } from "@/hooks/useOverlay";

export default function OverlayLeft() {
  const { isLeftOpen, rightNode } = useOverlay();

  return isLeftOpen && (
    <div className="h-full flex flex-col">
      {rightNode}
    </div>
  )
}
