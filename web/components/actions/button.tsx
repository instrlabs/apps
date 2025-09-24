"use client";

import React from "react";

type ButtonProps = React.ButtonHTMLAttributes<HTMLButtonElement> & {
  xSize?: "sm" | "md" | "lg";
  xVariant?: "solid" | "outline" | "transparent";
};

export default function Button({
  xSize = "md",
  xVariant = "outline",
  children,
  ...rest
}: ButtonProps) {
  // Container Classes
  const outlineClasses = [
    "bg-primary-black",
    "shadow-primary",
    "hover:bg-white/8"
  ].join(" ");

  const solidClasses = [
    "bg-white",
    "text-black",
    "hover:bg-white/85",
    "disabled:bg-white/85",
  ].join(" ");

  const transparentClasses = [
    "bg-transparent",
    "shadow-none",
  ].join(" ");

  const variantClasses =
    xVariant === "solid" ? solidClasses :
      xVariant === "outline" ? outlineClasses :
        xVariant === "transparent" ? transparentClasses :
          "";

  // Size Classes
  const smClasses = "py-2 px-3 text-sm min-w-24";
  const mdClasses = "py-3 px-4 min-w-28";
  const lgClasses = "py-4 px-6 text-base min-w-32";

  const sizeClasses =
    xSize === "sm" ? smClasses :
      xSize === "md" ? mdClasses :
        xSize === "lg" ? lgClasses :
          "";

  const baseClasses = "cursor-pointer font-medium rounded-lg transition-colors disabled:cursor-not-allowed";

  return (
    <button
      className={[
        baseClasses,
        sizeClasses,
        variantClasses,
      ].join(" ")}
      {...rest}
    >
      {children}
    </button>
  );
}
