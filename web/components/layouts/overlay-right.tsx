"use client";

import { useOverlay } from "@/hooks/useOverlay";

export default function OverlayRight() {
  const { isRightOpen, rightNode } = useOverlay();

  return isRightOpen && (
    <div className="flex flex-col">
      {rightNode}
    </div>
  )
}
