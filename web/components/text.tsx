import React from "react";
import clsx from "clsx";

export type TextVariant =
  | "default"
  | "secondary"
  | "muted"
  | "black"
  | "primary"
  | "success"
  | "warning"
  | "danger";

export type TextSize = "xs" | "sm" | "base" | "lg" | "xl";

export type TextProps<T extends React.ElementType = "span"> = {
  as?: T;
  children: React.ReactNode;
  className?: string;
  isBold?: boolean;
  xColor?: TextVariant;
  xSize?: TextSize;
} & Omit<React.ComponentPropsWithoutRef<T>, "as" | "className" | "children">;

export default function Text<T extends React.ElementType = "span">(
  {
    as,
    children,
    className,
    isBold = false,
    xColor = "default",
    xSize = "base",
    ...rest
  }: TextProps<T>
) {
  const Component = (as || "span") as React.ElementType;

  const variantClasses: Record<TextVariant, string> = {
    default: "text-white/90",
    secondary: "text-white/75",
    muted: "text-white/60",
    black: "text-black",
    primary: "text-blue-400",
    success: "text-green-400",
    warning: "text-yellow-400",
    danger: "text-red-400",
  };

  const sizeClasses: Record<TextSize, string> = {
    xs: "text-xs",
    sm: "text-sm",
    base: "text-base",
    lg: "text-lg",
    xl: "text-2xl"
  };

  return (
    <Component
      className={clsx(
        sizeClasses[xSize],
        isBold ? "font-semibold" : "font-normal",
        variantClasses[xColor],
        className,
      )}
      {...rest}
    >
      {children}
    </Component>
  );
}
