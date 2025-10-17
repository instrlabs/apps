"use client";

import React from "react";

type ButtonProps = React.ButtonHTMLAttributes<HTMLButtonElement> & {
  xSize?: "sm" | "md" | "lg";
  xColor?: "primary" | "secondary" | "transparent";
};

export default function Button({
  xSize = "md",
  xColor = "primary",
  className,
  children,
  ...rest
}: ButtonProps) {
  const primaryClasses = [
    "bg-primary",
    "text-black",
    "hover:opacity-90",
    "disabled:opacity-60"
  ].join(" ");

  const secondaryClasses = [
    "bg-secondary",
    "text-secondary",
    "border",
    "border-primary",
    "hover:opacity-90",
    "disabled:opacity-60"
  ].join(" ");

  const transparentClasses = [
    "bg-transparent",
    "text-primary",
    "hover:opacity-90",
    "disabled:opacity-60"
  ].join(" ");

  const colorClasses =
    xColor === "primary" ? primaryClasses :
    xColor === "secondary" ? secondaryClasses :
    transparentClasses;

  const smClasses = "py-1.5 px-4 text-sm";
  const mdClasses = "py-2 px-4";
  const lgClasses = "py-3 px-6 text-base min-w-32";

  const sizeClasses =
    xSize === "sm" ? smClasses :
    xSize === "md" ? mdClasses :
    xSize === "lg" ? lgClasses :
    "";

  const baseClasses = `
inline-flex items-center justify-center
rounded-md
cursor-pointer
hover:opacity-90
transition-opacity transition-colors
focus:outline-none
disabled:cursor-not-allowed
  `;

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
