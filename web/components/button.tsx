"use client";

import React from "react";
import Icon from "./icon";

type ButtonProps = React.ButtonHTMLAttributes<HTMLButtonElement> & {
  label?: React.ReactNode;
  leftIconName?: string | null;
  rightIconName?: string | null;
  hasLeftIcon?: boolean;
  hasRightIcon?: boolean;
  size?: "sm" | "base" | "lg";
  color?: "primary" | "secondary" | "transparent";
};

export default function Button({
  label,
  leftIconName = null,
  rightIconName = null,
  hasLeftIcon = false,
  hasRightIcon = false,
  size = "base",
  color = "primary",
  className = "",
  children,
  disabled = false,
  ...rest
}: ButtonProps) {
  const baseClasses = [
    "box-border",
    "inline-flex",
    "items-center",
    "justify-center",
    "rounded-border",
    "cursor-pointer",
    "transition-opacity",
    "focus:outline-none",
    "disabled:cursor-not-allowed",
  ].join(" ");

  // Spacing and typography per size (matching Figma design)
  const sizeConfig = {
    sm: {
      spacing: "gap-spacing-2 p-spacing-2",
      font: "text-sm leading-5",
      weight: "font-medium",
      iconSize: 20
    },
    base: {
      spacing: "gap-spacing-2 p-spacing-2",
      font: "text-base leading-6",
      weight: "font-medium",
      iconSize: 24
    },
    lg: {
      spacing: "gap-spacing-3 p-spacing-3",
      font: "text-base leading-6",
      weight: "font-semibold",
      iconSize: 24
    }
  };

  // Color variants (matching Figma design system)
  const colorConfig = {
    primary: [
      "btn-primary",
      "hover:btn-primary-hover",
      "disabled:btn-primary-disabled"
    ].join(" "),
    secondary: [
      "btn-secondary",
      "border",
      "opacity-90",
      "hover:opacity-100",
      "hover:btn-secondary-hover",
      "disabled:btn-secondary-disabled"
    ].join(" "),
    transparent: [
      "bg-transparent",
      "text-primary",
      "hover:opacity-90",
      "disabled:opacity-60"
    ].join(" ")
  };

  const currentSize = sizeConfig[size] || sizeConfig.base;
  const currentColor = colorConfig[color] || colorConfig.primary;

  const renderIcon = (icon: React.ReactNode | string | null) => {
    if (!icon) return null;

    const iconOpacity = disabled ? "opacity-60" : "opacity-[0.99]";

    if (typeof icon === "string") {
      return (
        <span className={["relative", "shrink-0", iconOpacity].join(" ")}>
          <Icon name={icon} size={currentSize.iconSize} />
        </span>
      );
    }
  };

  return (
    <button
      className={[
        baseClasses,
        currentSize.spacing,
        currentSize.font,
        currentSize.weight,
        currentColor,
        className
      ].filter(Boolean).join(" ")}
      disabled={disabled}
      {...rest}
    >
      {hasLeftIcon && renderIcon(leftIconName)}
      <span className="relative shrink-0">
        {children ?? label ?? "Button Text"}
      </span>
      {hasRightIcon && renderIcon(rightIconName)}
    </button>
  );
}
