"use client";

import React from "react";
import Icon from "./icon";

type InputProps = Omit<React.InputHTMLAttributes<HTMLInputElement>, "size"> & {
  leftIconName?: string;
  rightIconName?: string;
  hasLeftIcon?: boolean;
  hasRightIcon?: boolean;
  size?: "sm" | "base" | "lg";
  variant?: "primary" | "secondary";
};

export default function Input({
  leftIconName,
  rightIconName,
  hasLeftIcon = false,
  hasRightIcon = false,
  size = "base",
  variant = "primary",
  className = "",
  disabled,
  ...rest
}: InputProps) {
  // Base classes - structure and layout
  const baseClasses = "flex items-center rounded min-w-[200px] transition-colors focus-within:outline-none";

  // Size configuration - spacing, typography, and icon size
  const sizeConfig: Record<"sm" | "base" | "lg", { container: string; input: string; iconSize: number }> = {
    sm: {
      container: "gap-2 p-2",
      input: "text-sm leading-5",
      iconSize: 20,
    },
    base: {
      container: "gap-2 p-2 h-10",
      input: "text-base leading-6",
      iconSize: 24,
    },
    lg: {
      container: "gap-3 p-3 h-12",
      input: "text-base leading-6",
      iconSize: 24,
    },
  };

  // Variant configuration - colors and states
  const variantConfig: Record<"primary" | "secondary", { container: string; input: string }> = {
    primary: {
      container: "bg-white/4 border border-white/10",
      input: "text-white placeholder:text-white/30",
    },
    secondary: {
      container: "bg-white",
      input: "text-black placeholder:text-black/60",
    },
  };

  const currentSize = sizeConfig[size];
  const currentVariant = variantConfig[variant];

  const inputClasses = [
    "flex-1",
    "min-w-0",
    "bg-transparent",
    "border-none",
    "outline-none",
    currentSize.input,
    "font-normal",
    currentVariant.input,
  ].join(" ");

  const renderIcon = (iconName: string | null) => {
    if (!iconName) return null;

    return (
      <span className="relative shrink-0">
        <Icon name={iconName} size={currentSize.iconSize} />
      </span>
    );
  };

  return (
    <div
      className={[
        baseClasses,
        currentSize.container,
        currentVariant.container,
        disabled && "opacity-50",
        className,
      ]
        .filter(Boolean)
        .join(" ")}
    >
      {hasLeftIcon && leftIconName && renderIcon(leftIconName)}
      <input
        className={inputClasses}
        disabled={disabled}
        data-variant={variant}
        {...rest}
      />
      {hasRightIcon && rightIconName && renderIcon(rightIconName)}
    </div>
  );
}
