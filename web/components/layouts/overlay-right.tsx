"use client";

import { useOverlay } from "@/hooks/useOverlay";

export default function OverlayRight() {
  const { rightNode, rightKey, rightWidth } = useOverlay();

  return (
    <div className="flex flex-col">
      {rightNode}
    </div>
  )
}
