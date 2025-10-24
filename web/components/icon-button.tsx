"use client";

import React from "react";

export type IconButtonProps = React.ButtonHTMLAttributes<HTMLButtonElement> & {
  xSize?: "sm" | "md" | "lg";
  xColor?: "primary" | "secondary" | "transparent";
};

export default function IconButton({
  xSize = "md",
  xColor = "primary",
  className,
  children,
  ...rest
}: IconButtonProps) {
  const primaryClasses = [
    "bg-primary",
    "text-black",
    "hover:opacity-90",
    "disabled:opacity-60",
  ].join(" ");

  const secondaryClasses = [
    "bg-secondary",
    "text-secondary",
    "border",
    "border-primary",
    "hover:opacity-90",
    "disabled:opacity-60",
  ].join(" ");

  const transparentClasses = [
    "bg-transparent",
    "text-primary",
    "hover:opacity-90",
    "disabled:opacity-60",
  ].join(" ");

  const colorClasses =
    xColor === "primary" ? primaryClasses :
    xColor === "secondary" ? secondaryClasses :
    transparentClasses;

  const smClasses = "w-8 h-8";
  const mdClasses = "w-10 h-10";
  const lgClasses = "w-12 h-12";

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
        className
      ].join(" ")}
      {...rest}
    >
      {children}
    </button>
  );
}
