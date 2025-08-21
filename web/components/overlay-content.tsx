"use client";

import { ReactNode, useMemo } from "react";
import { usePathname } from "next/navigation";
import { useOverlay } from "@/hooks/useOverlay";

export type OverlayAnimation = "fadeIn" | "bounceInLeft" | "none";

export default function OverlayContent({
  children,
  animation = "fadeIn",
  contentKey,
  durationMs,
}: {
  children: ReactNode;
  animation?: OverlayAnimation;
  contentKey?: string | number;
  durationMs?: number;
}) {
  const pathname = usePathname();
  const { leftWidth, rightWidth, isLeftOpen, isRightOpen } = useOverlay();

  // Determine a stable key for the content so that when it changes, we replay the animation
  const keyForContent = useMemo(() => {
    if (contentKey !== undefined && contentKey !== null) return String(contentKey);
    return pathname ?? "static";
  }, [contentKey, pathname]);

  const animationClass = animation === "fadeIn"
    ? "animate-fade-in"
    : animation === "bounceInLeft"
    ? "animate-bounce-in-left"
    : undefined;

  const leftTargetPx = Number.isFinite(leftWidth) ? Math.max(0, Math.round(leftWidth)) : 0;
  const rightTargetPx = Number.isFinite(rightWidth) ? Math.max(0, Math.round(rightWidth)) : 0;
  const leftPx = isLeftOpen ? leftTargetPx : 0;
  const rightPx = isRightOpen ? rightTargetPx : 0;
  const styleVars: React.CSSProperties = {
    left: `${leftPx}px`,
    right: `${rightPx}px`,
  } as React.CSSProperties;

  const innerStyle: React.CSSProperties = durationMs
    ? ({ ["--overlay-anim-duration"]: `${durationMs}ms` } as React.CSSProperties)
    : {};

  return (
    <div
      className="absolute top-0 bottom-0 p-3 pt-[80px] transition-[left,right] duration-300 ease-in-out"
      style={styleVars}
    >
      <div className="w-full h-full rounded-xl bg-card overflow-auto">
        <div
          key={keyForContent}
          className={"h-full w-full flex " + (animationClass ? ` ${animationClass}` : "")}
          style={innerStyle}
        >
          {children}
        </div>
      </div>
    </div>
  );
}
