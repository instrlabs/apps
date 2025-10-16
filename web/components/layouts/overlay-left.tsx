"use client";

import { useOverlay } from "@/hooks/useOverlay";

export default function OverlayLeft() {
  const { isLeftOpen, rightNode } = useOverlay();

  return isLeftOpen && (
    <div className="flex flex-col">
      {rightNode}
    </div>
  )
}
