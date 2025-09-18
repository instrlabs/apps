"use client";

import React from "react";
import clsx from "clsx";

type ButtonProps = Omit<React.ButtonHTMLAttributes<HTMLButtonElement>, "onClick"> & {
  xSize?: "sm" | "md" | "lg";
  xColor?: "primary" | "secondary";
  isLoading?: boolean;
  loadingText?: string;
  onClick?: (event: React.MouseEvent<HTMLButtonElement>) => void | Promise<void>;
};

export default function Button({
  type = "button",
  children,
  onClick,
  className,
  disabled,
  xSize = "md",
  xColor = "primary",
  isLoading,
  ...rest
}: ButtonProps) {
  const [internalLoading, setInternalLoading] = React.useState(false);

  const handleClick = async (e: React.MouseEvent<HTMLButtonElement>) => {
    if (!onClick) return;
    if (internalLoading) return;

    try {
      const maybePromise = onClick(e);
      const isPromise = !!maybePromise && typeof (maybePromise as PromiseLike<unknown>).then === "function";
      if (isPromise) {
        setInternalLoading(true);
        try {
          await (maybePromise as Promise<void>);
        } finally {
          setInternalLoading(false);
        }
      }
    } catch (err) {
      throw err;
    }
  };

  // Size classes tailored for text buttons with sensible min-widths for cleaner UI
  const sizeClass =
    xSize === "sm"
      ? "py-2 px-3 text-sm min-w-24"
      : xSize === "md"
      ? "py-3 px-4 min-w-28"
      : "py-4 px-6 text-base min-w-32";

  const colorClasses =
    xColor === "primary"
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

  const isDisabled = disabled || isLoading || internalLoading;

  return (
    <button
      type={type}
      className={clsx(
        colorClasses,
        "border border-[var(--btn-border)]",
        sizeClass,
        "rounded-lg cursor-pointer font-medium shadow-primary",
        isDisabled && "opacity-70 cursor-not-allowed",
        className
      )}
      disabled={isDisabled}
      onClick={handleClick}
      {...rest}
    >
      {children}
    </button>
  );
}
