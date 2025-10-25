"use client";

import React from "react";
import Icon from "./icon";

type ButtonProps = React.ButtonHTMLAttributes<HTMLButtonElement> & {
  label?: React.ReactNode;
  leftIconName?: string;
  rightIconName?: string;
  hasLeftIcon?: boolean;
  hasRightIcon?: boolean;
  size?: "sm" | "base" | "lg";
  color?: "primary" | "secondary";
};

export default function Button({
  label,
  leftIconName,
  rightIconName,
  hasLeftIcon = false,
  hasRightIcon = false,
  size = "base",
  color = "primary",
  className = "",
  children,
  disabled,
  ...rest
}: ButtonProps) {

  const sizeConfig = {
    sm: {
      spacing: "gap-2 p-2",
      font: "text-sm",
      lineHeight: "leading-5",
      weight: "font-medium",
      iconWidth: 20,
      iconHeight: 17.889,
    },
    base: {
      spacing: "gap-2 p-2",
      font: "text-base",
      lineHeight: "leading-6",
      weight: "font-medium",
      iconWidth: 24,
      iconHeight: 21.909,
    },
    lg: {
      spacing: "gap-3 p-3",
      font: "text-base",
      lineHeight: "leading-6",
      weight: "font-semibold",
      iconWidth: 24,
      iconHeight: 24,
    },
  };

  const colorConfig = {
    primary: {
      default: "btn-primary",
      hover: "btn-primary-hover",
      disabled: "btn-primary-disabled",
    },
    secondary: {
      default: "btn-secondary border",
      hover: "btn-secondary-hover border",
      disabled: "btn-secondary-disabled border",
    },
  };

  const currentSize = sizeConfig[size] || sizeConfig.base;
  const currentColor = colorConfig[color] || colorConfig.primary;

  // Use default color - CSS pseudo-classes (:hover, :disabled) handle state styling
  const stateStyle = currentColor.default;

  const baseClasses = [
    "box-border",
    "inline-flex",
    "items-center",
    "justify-center",
    "rounded",
    "cursor-pointer",
    "transition-colors",
    "focus:outline-none",
    "disabled:cursor-not-allowed",
  ].join(" ");

  const renderIcon = (icon: string) => {
    if (!icon) return null;

    const iconOpacity = disabled ? "opacity-60" : "opacity-[0.99]";

    return (
      <span className={["relative", "shrink-0", iconOpacity].filter(Boolean).join(" ")}>
        <Icon name={icon} size={currentSize.iconWidth} />
      </span>
    );
  };


  return (
    <button
      className={[
        baseClasses,
        currentSize.spacing,
        currentSize.font,
        currentSize.lineHeight,
        currentSize.weight,
        stateStyle,
        className,
      ].filter(Boolean).join(" ")}
      disabled={disabled}
      {...rest}
    >
      {hasLeftIcon && leftIconName && renderIcon(leftIconName)}
      <span className="relative shrink-0">
        {children ?? label ?? "Button Text"}
      </span>
      {hasRightIcon && rightIconName && renderIcon(rightIconName)}
    </button>
  );
}
