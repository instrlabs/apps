"use client";

import { useOverlay } from "@/hooks/useOverlay";

export default function OverlayRight() {
  const { rightNode, rightKey, rightWidth } = useOverlay();

  return (
    <div
      className={
        "absolute top-0 right-0 bottom-0 " +
        "pt-[80px] h-screen " +
        "transition-[width] duration-300 ease-in-out"
      }
      style={{ width: `${rightWidth}px` }}
    >
      <div className="w-full h-full flex flex-col">
        <div key={rightKey} className="flex-1 animate-fade-in">
          {rightNode}
        </div>
      </div>
    </div>
  )
}
