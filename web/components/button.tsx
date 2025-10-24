"use client";

import React from "react";

type ButtonProps = React.ButtonHTMLAttributes<HTMLButtonElement> & {
  label?: React.ReactNode;
  leftIcon?: React.ReactNode | null;
  rightIcon?: React.ReactNode | null;
  hasLeftIcon?: boolean;
  hasRightIcon?: boolean;
  size?: "sm" | "base" | "lg";
  color?: "primary" | "secondary" | "transparent";
};

export default function Button({
  label,
  leftIcon = null,
  rightIcon = null,
  hasLeftIcon = false,
  hasRightIcon = false,
  size = "base",
  color = "primary",
  className = "",
  children,
  ...rest
}: ButtonProps) {
  const effectiveSize = size;
  const effectiveColor = color;

  const baseClasses = [
    "box-border",
    "inline-flex",
    "items-center",
    "justify-center",
    "rounded",
    "cursor-pointer",
    "transition-opacity",
    "focus:outline-none",
    "disabled:cursor-not-allowed",
  ].join(" ");

  // spacing/gap and padding per size
  const sizeMap: Record<string, string> = {
    sm: "gap-2 p-2 text-sm",
    base: "gap-2 p-2 text-base",
    lg: "gap-3 p-3 text-base",
  };

  // icon sizing per size
  const iconMap: Record<string, string> = {
    sm: "h-[18px] w-[20px]",
    base: "h-[22px] w-[24px]",
    lg: "h-[24px] w-[24px]",
  };

  // color/variant classes mapped to the Figma visuals
  const variants: Record<string, string> = {
    primary: [
      "bg-[rgba(255,255,255,0.99)]",
      "text-black",
      "hover:bg-white",
      "disabled:bg-[rgba(255,255,255,0.6)]",
    ].join(" "),
    secondary: [
      "bg-secondary",
      "border",
      "border-primary",
      "text-primary",
      "hover:opacity-100",
      "disabled:opacity-50",
    ].join(" "),
    transparent: [
      "bg-transparent",
      "text-primary",
      "hover:opacity-90",
      "disabled:opacity-60",
    ].join(" "),
  };

  const sizeClasses = sizeMap[effectiveSize] ?? sizeMap.base;
  const colorClasses = variants[effectiveColor] ?? variants.primary;

  const renderIcon = (icon: React.ReactNode | null, which: "left" | "right") => {
    if (!icon) return null;
    const common = [
      iconMap[effectiveSize],
      which === "left" ? "mr-0" : "ml-0",
    ].join(" ");
    return <span className={common}>{icon}</span>;
  };

  const content = (
    <>
      {hasLeftIcon && (leftIcon ? renderIcon(leftIcon, "left") : null)}
      <span className="truncate font-medium">
        {children ?? label ?? "Button"}
      </span>
      {hasRightIcon && (rightIcon ? renderIcon(rightIcon, "right") : null)}
    </>
  );

  return (
    <button
      className={[baseClasses, sizeClasses, colorClasses, className].join(" ")}
      {...rest}
    >
      {content}
    </button>
  );
}
