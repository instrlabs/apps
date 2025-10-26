"use client";

import React from "react";

export type IconButtonProps = React.ButtonHTMLAttributes<HTMLButtonElement> & {
  size?: "sm" | "base" | "lg";
  variant?: "primary" | "secondary";
};

export default function IconButton({
  size = "base",
  variant = "secondary",
  className = "",
  children,
  ...rest
}: IconButtonProps) {
  // Base classes - structure and layout
  const baseClasses = "inline-flex items-center justify-center rounded transition-colors focus:outline-none disabled:cursor-not-allowed disabled:opacity-60";

  // Size configuration - spacing and icon size
  const sizeConfig: Record<"sm" | "base" | "lg", string> = {
    sm: "gap-2 p-2",
    base: "gap-2 p-2",
    lg: "gap-3 p-3",
  };

  // Variant configuration - colors and states
  const variantConfig: Record<"primary" | "secondary", string> = {
    primary: "bg-white text-black hover:bg-white/90",
    secondary: "bg-white/8 border border-white/10 text-white hover:bg-white/12",
  };

  return (
    <button
      className={[
        baseClasses,
        sizeConfig[size],
        variantConfig[variant],
        className,
      ]
        .filter(Boolean)
        .join(" ")}
      {...rest}
    >
      {children}
    </button>
  );
}
