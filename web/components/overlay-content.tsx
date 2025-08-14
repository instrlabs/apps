"use client";

import { ReactNode, useMemo } from "react";
import { usePathname } from "next/navigation";

export type OverlayAnimation = "fadeIn" | "bounceInLeft" | "none";

export default function OverlayContent({
  children,
  animation = "fadeIn",
  contentKey,
  durationMs,
}: {
  children: ReactNode;
  animation?: OverlayAnimation;
  /**
   * Optional key to indicate when "new content" is set. When this value changes,
   * the inner content is re-mounted and the animation replays.
   * If not provided, we fall back to the current pathname (client-side).
   */
  contentKey?: string | number;
  /** Optional override for animation duration in ms. */
  durationMs?: number;
}) {
  const pathname = usePathname();

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

  const styleVars: React.CSSProperties = {
    left: "var(--overlay-left-width, 300px)",
    right: "var(--overlay-right-width, 300px)",
  } as React.CSSProperties;

  // Allow overriding animation duration via CSS variable
  const innerStyle: React.CSSProperties = durationMs
    ? ({ ["--overlay-anim-duration"]: `${durationMs}ms` } as React.CSSProperties)
    : {};

  return (
    <div
      className="absolute top-0 bottom-0 p-3 pt-[80px] transition-[left,right] duration-300 ease-in-out"
      style={styleVars}
    >
      <div className="w-full h-full rounded-3xl bg-neutral-50">
        <div
          key={keyForContent}
          className={"h-full w-full flex items-center justify-center text-gray-700" + (animationClass ? ` ${animationClass}` : "")}
          style={innerStyle}
        >
          {children}
        </div>
      </div>
    </div>
  );
}
