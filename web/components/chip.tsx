"use client";

import React from "react";

type ChipState = "processing" | "success" | "error" | "info" | "default";

type ChipProps = {
  label?: string;
  state?: ChipState;
  className?: string;
};

export default function Chip({
  label = "TEXT CHIP",
  state = "processing",
  className = "",
}: ChipProps) {
  // Base classes - shared structure and layout
  const baseClasses =
    "inline-flex items-center justify-center rounded border border-solid gap-2.5 px-2 py-1";

  // State configuration - colors and borders
  const stateConfig: Record<ChipState, string> = {
    processing:
      "bg-yellow-400/40 border-yellow-400",
    success:
      "bg-emerald-500/40 border-emerald-400",
    error:
      "bg-red-400/40 border-red-400",
    info:
      "bg-blue-400/40 border-blue-400",
    default:
      "bg-white/8 border-white/10",
  };

  const currentState = stateConfig[state];

  return (
    <div className={[baseClasses, currentState, className].filter(Boolean).join(" ")}>
      <p className="text-xs leading-3 font-medium text-white/80">{label}</p>
    </div>
  );
}
