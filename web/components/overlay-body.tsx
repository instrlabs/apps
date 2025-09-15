"use client";

import React, { useEffect } from "react";
import clsx from "clsx";
import {useOverlay} from "@/hooks/useOverlay";

export default function OverlayBody({ children }: {
  children: React.ReactNode;
}) {
  const { isLeftOpen, isRightOpen, leftWidth, rightWidth, isModalOpen, closeAll } = useOverlay();

  useEffect(() => {
    const onKeyDown = (e: KeyboardEvent) => {
      if (e.key === "Escape") {
        if (isLeftOpen || isRightOpen || isModalOpen) {
          e.stopPropagation();
          closeAll();
        }
      }
    };
    document.addEventListener("keydown", onKeyDown);
    return () => document.removeEventListener("keydown", onKeyDown);
  }, [isLeftOpen, isRightOpen, isModalOpen, closeAll]);

  const leftTargetPx = Number.isFinite(leftWidth as number)
    ? Math.max(0, Math.round(leftWidth as number))
    : 0;
  const rightTargetPx = Number.isFinite(rightWidth as number)
    ? Math.max(0, Math.round(rightWidth as number))
    : 0;
  const leftPx = isLeftOpen ? leftTargetPx : 0;
  const rightPx = isRightOpen ? rightTargetPx : 0;

  return (
    <div
      className={clsx(
        "absolute top-0 bottom-0 py-3 pt-[80px] px-3",
        "transition-[left,right,padding] duration-300 ease-in-out",
      )}
      style={{ left: `${leftPx}px`, right: `${rightPx}px` }}
    >
      <div className="w-full h-full rounded-xl bg-card overflow-auto">
        <div className="h-full w-full flex animate-fade-in">{children}</div>
      </div>
    </div>
  );
}
