"use client";

import React from "react";
import clsx from "clsx";
import {useOverlay} from "@/hooks/useOverlay";

export default function OverlayBody({ children }: {
  children: React.ReactNode;
}) {
  // const { isLeftOpen, isRightOpen, leftWidth, rightWidth } = useOverlay();
  //
  // const leftTargetPx = Number.isFinite(leftWidth as number)
  //   ? Math.max(0, Math.round(leftWidth as number))
  //   : 0;
  // const rightTargetPx = Number.isFinite(rightWidth as number)
  //   ? Math.max(0, Math.round(rightWidth as number))
  //   : 0;
  // const leftPx = isLeftOpen ? leftTargetPx : 0;
  // const rightPx = isRightOpen ? rightTargetPx : 0;

  return (
    // <div
    //   className={clsx(
    //     // "transition-[left,right] duration-300 ease-in-out",
    //   )}
    //   // style={{ left: `${leftPx}px`, right: `${rightPx}px` }}
    // >
      <div className="flex flex-col">
        {children}
      </div>
    // </div>
  );
}
