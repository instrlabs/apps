"use client";

import { useOverlay } from "@/hooks/useOverlay";
import { useMemo } from "react";
import clsx from "clsx";

export default function OverlayRight() {
  const { isRightOpen, rightNode, rightWidth } = useOverlay();
  const targetWidth = Math.max(0, rightWidth);
  const widthPx = useMemo(() => isRightOpen ? targetWidth : 0, [isRightOpen, targetWidth]);

  return (
    <div
      className={clsx(
        "absolute top-0 right-0 bottom-0",
        "pt-[80px] h-screen",
        "transition-[width] duration-300 ease-in-out"
      )}
      style={{
        width: `${widthPx}px`,
        pointerEvents: isRightOpen ? 'auto' : 'none'
      }}
    >
      <div className={clsx(
        "w-full h-full flex flex-col",
        isRightOpen ? "overflow-visible" : "overflow-hidden",
      )}>
        <div className="flex-1 animate-fade-in">
          {rightNode}
        </div>
      </div>
    </div>
  )
}
