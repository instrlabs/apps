"use client";

import React from "react";
import clsx from "clsx";


type ButtonIconProps = React.ButtonHTMLAttributes<HTMLButtonElement> & {
  xsize: "sm" | "md" | "lg";
  xVariant?: "solid" | "outline" | "transparent";
};

export default function ButtonIcon({
  xsize,
  xVariant = "outline",
  children,
  ...rest
}: ButtonIconProps) {
  // Container Classes
  const outlineClasses = [
    "bg-white/2",
    "text-white/90",
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
  const smClasses = "size-8";
  const mdClasses = "size-10";
  const lgClasses = "size-12";

  const sizeClass =
    xsize === "sm" ? smClasses :
      xsize === "md" ? mdClasses :
        xsize === "lg" ? lgClasses :
          "";

  const baseClasses = `
  flex items-center justify-center
  rounded-full transition-colors aspect-square
  cursor-pointer disabled:cursor-not-allowed
  `;

  return (
    <button
      className={clsx(
        baseClasses,
        variantClasses,
        sizeClass,
      )}
      {...rest}
    >
      {children}
    </button>
  );
}
