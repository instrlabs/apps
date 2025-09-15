"use client";

import React from "react";
import clsx from "clsx";


type ButtonIconProps = React.ButtonHTMLAttributes<HTMLButtonElement> & {
  xSize: "sm" | "md" | "lg";
  xColor: "primary" | "secondary";
};

export default function ButtonIcon({
  className,
  children,
  type = "button",
  xSize,
  xColor,
  ...rest
}: ButtonIconProps) {
  const sizeClass =
    xSize === "sm" ? "p-1"
      : xSize === "md" ? "p-[6px]"
        : "p-2";

  const colorClasses = xColor === "primary"
    ? [
        "bg-[var(--btn-primary-bg)]",
        "text-[var(--btn-primary-text)]",
        "hover:bg-[var(--btn-primary-hover)]",
        "active:bg-[var(--btn-primary-active)]",
        "disabled:bg-[var(--btn-primary-disabled)]",
      ]
    : [
        "bg-[var(--btn-secondary-bg)]",
        "text-[var(--btn-secondary-text)]",
        "hover:bg-[var(--btn-secondary-hover)]",
        "active:bg-[var(--btn-secondary-active)]",
        "disabled:bg-[var(--btn-secondary-disabled)]",
      ];

  return (
    <button
      type={type}
      className={clsx(
        colorClasses,
        "border border-[var(--btn-border)]",
        sizeClass,
        "rounded-full cursor-pointer shadow-primary",
        className
      )}
      {...rest}
    >
      {children}
    </button>
  );
}
