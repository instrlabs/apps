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
  color?: "primary" | "secondary";
  state?: "Default" | "Hover" | "Disabled";
};

export default function Button({
  label,
  leftIconName = null,
  rightIconName = null,
  hasLeftIcon = false,
  hasRightIcon = false,
  size = "base",
  color = "primary",
  state = "Default",
  className = "",
  children,
  disabled = false,
  onMouseEnter,
  onMouseLeave,
  ...rest
}: ButtonProps) {
  const [isHovered, setIsHovered] = React.useState(false);

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

  // Determine current state
  const currentState = disabled ? "Disabled" : isHovered ? "Hover" : "Default";
  const stateStyle = disabled ? currentColor.disabled : isHovered ? currentColor.hover : currentColor.default;

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

  const renderIcon = (icon: React.ReactNode | string | null) => {
    if (!icon) return null;

    const iconOpacity = disabled ? "opacity-60" : currentState === "Hover" ? "" : "opacity-[0.99]";

    if (typeof icon === "string") {
      return (
        <span className={["relative", "shrink-0", iconOpacity].filter(Boolean).join(" ")}>
          <Icon name={icon} size={currentSize.iconWidth} />
        </span>
      );
    }

    return icon;
  };

  const handleMouseEnter = (e: React.MouseEvent<HTMLButtonElement>) => {
    setIsHovered(true);
    onMouseEnter?.(e);
  };

  const handleMouseLeave = (e: React.MouseEvent<HTMLButtonElement>) => {
    setIsHovered(false);
    onMouseLeave?.(e);
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
      onMouseEnter={handleMouseEnter}
      onMouseLeave={handleMouseLeave}
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
