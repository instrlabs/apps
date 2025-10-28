"use client";

import React from "react";

export type IconButtonProps = React.ButtonHTMLAttributes<HTMLButtonElement> & {
  size?: "sm" | "base" | "lg";
  variant?: "primary" | "secondary" | "transparent";
};

export default function IconButton({
  size = "base",
  variant = "secondary",
  className = "",
  children,
  ...rest
}: IconButtonProps) {
  // Base classes - structure and layout
  const baseClasses =
    "inline-flex items-center justify-center rounded transition-colors focus:outline-none disabled:cursor-not-allowed disabled:opacity-60";

  // Size configuration - fixed dimensions and padding
  const sizeConfig: Record<"sm" | "base" | "lg", string> = {
    sm: "size-9 p-2",
    base: "h-10 w-10 p-2",
    lg: "size-12 p-3",
  };

  // Variant configuration - colors and states
  const variantConfig: Record<"primary" | "secondary" | "transparent", string> = {
    primary: "bg-white text-black hover:bg-white/90",
    secondary: "bg-white/8 border border-white/10 text-white hover:bg-white/12",
    transparent: "text-white hover:bg-white/8",
  };

  return (
    <button
      className={[baseClasses, sizeConfig[size], variantConfig[variant], className]
        .filter(Boolean)
        .join(" ")}
      {...rest}
    >
      {children}
    </button>
  );
}
