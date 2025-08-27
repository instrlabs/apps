"use client";

import React from "react";
import clsx from "clsx";
import {useOverlay} from "@/hooks/useOverlay";

export default function OverlayContentWrapper({ children }: {
  children: React.ReactNode;
}) {
  const { isLeftOpen, isRightOpen, leftWidth, rightWidth } = useOverlay();

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
        "absolute top-0 bottom-0 py-3 pt-[80px]",
        "transition-[left,right,padding] duration-300 ease-in-out",
        isLeftOpen ? "pl-0" : "pl-3"
      )}
      style={{ left: `${leftPx}px`, right: `${rightPx}px` }}
    >
      <div className="w-full h-full rounded-xl bg-card overflow-auto">
        <div className="h-full w-full flex animate-fade-in">{children}</div>
      </div>
    </div>
  );
}
