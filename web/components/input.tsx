"use client";

import React from "react";
import Icon from "./icon";

type InputProps = React.InputHTMLAttributes<HTMLInputElement> & {
  leftIconName?: string | null;
  rightIconName?: string | null;
  hasLeftIcon?: boolean;
  hasRightIcon?: boolean;
  size?: "sm" | "base" | "lg";
  color?: "primary" | "secondary";
};

export default function Input({
  leftIconName = null,
  rightIconName = null,
  hasLeftIcon = false,
  hasRightIcon = false,
  size = "base",
  color = "primary",
  className = "",
  disabled,
  ...rest
}: InputProps) {
  const sizeConfig: Record<
    "sm" | "base" | "lg",
    {
      spacing: string;
      font: string;
      lineHeight: string;
      height: string;
      iconSize: number;
    }
  > = {
    sm: {
      spacing: "gap-2 p-2",
      font: "text-sm",
      lineHeight: "leading-5",
      height: "",
      iconSize: 20,
    },
    base: {
      spacing: "gap-2 p-2",
      font: "text-base",
      lineHeight: "leading-6",
      height: "h-10",
      iconSize: 24,
    },
    lg: {
      spacing: "gap-3 p-3",
      font: "text-base",
      lineHeight: "leading-6",
      height: "h-12",
      iconSize: 24,
    },
  };

  const colorConfig: Record<"primary" | "secondary", string> = {
    primary: "input-primary border",
    secondary: "input-secondary",
  };

  const currentSize = sizeConfig[size];
  const colorStyle = colorConfig[color];

  const baseClasses = [
    "box-border",
    "flex",
    "items-center",
    "rounded",
    "min-w-[200px]",
    "transition-colors",
    "focus-within:outline-none",
  ].join(" ");

  const inputClasses = [
    "flex-1",
    "min-w-0",
    "bg-transparent",
    "border-none",
    "outline-none",
    currentSize.font,
    currentSize.lineHeight,
    "font-normal",
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
        currentSize.spacing,
        currentSize.height,
        colorStyle,
        className,
      ]
        .filter(Boolean)
        .join(" ")}
    >
      {hasLeftIcon && renderIcon(leftIconName)}
      <input
        className={inputClasses}
        disabled={disabled}
        {...rest}
      />
      {hasRightIcon && renderIcon(rightIconName)}
    </div>
  );
}
