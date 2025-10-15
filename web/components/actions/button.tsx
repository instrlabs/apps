"use client";

import React from "react";

type ButtonProps = React.ButtonHTMLAttributes<HTMLButtonElement> & {
  xSize?: "sm" | "md" | "lg";
  xColor?: "primary" | "secondary";
};

export default function Button({
  xSize = "md",
  xColor = "primary",
  className,
  children,
  ...rest
}: ButtonProps) {
  const primaryClasses = [
    "bg-white",
    "text-black",
    "hover:bg-white/85",
    "disabled:bg-white/85"
  ].join(" ");

  const secondaryClasses = [
    "bg-white/8",
    "text-white/90",
    "border",
    "border-white/10",
    "hover:bg-white/10",
    "disabled:opacity-60"
  ].join(" ");

  const colorClasses = xColor === "primary"
    ? primaryClasses
    : secondaryClasses;

  const smClasses = "py-1.5 px-4 text-sm";
  const mdClasses = "py-2 px-4";
  const lgClasses = "py-3 px-6 text-base min-w-32";

  const sizeClasses =
    xSize === "sm" ? smClasses :
      xSize === "md" ? mdClasses :
        xSize === "lg" ? lgClasses :
          "";

  const baseClasses = "cursor-pointer font-medium rounded transition-colors disabled:cursor-not-allowed";

  return (
    <button
      className={[
        baseClasses,
        sizeClasses,
        colorClasses,
        className,
      ].join(" ")}
      {...rest}
    >
      {children}
    </button>
  );
}
