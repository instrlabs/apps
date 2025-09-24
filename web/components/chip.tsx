import React from "react";
import clsx from "clsx";

export type ChipSize = "sm" | "md" | "lg";
export type ChipVariant = "filled" | "outline" | "outlined"; // 'outlined' kept for backward-compat
export type ChipColor = "default" | "primary" | "error" | "warning" | "info";

export interface ChipProps extends React.HTMLAttributes<HTMLDivElement> {
  children?: React.ReactNode;
  onClick?: (e: React.MouseEvent<HTMLDivElement>) => void;
  xSize?: ChipSize;
  xVariant?: ChipVariant;
  xColor?: ChipColor;
  loading?: boolean; // when true, apply blinking color animation
}

const Chip = React.forwardRef<HTMLDivElement, ChipProps>(
  (
    {
      onClick,
      xSize = "md",
      xVariant = "filled",
      xColor = "default",
      className,
      children,
      loading = false,
      ...rest
    },
    ref
  ) => {
    // Normalize legacy value
    const variant: "filled" | "outline" = xVariant === "outlined" ? "outline" : xVariant;

    // Match Button sizing (padding and text)
    const sizeClasses =
      xSize === "sm"
        ? "py-1 px-2 text-xs gap-2"
        : "";

    // Variant + Color classes
    const filledDefault = "bg-primary-black text-white shadow-primary";
    const filledPrimary = "bg-primary text-primary-foreground";

    const outlineDefault = "border border-border bg-transparent text-[color:var(--text-primary)]";
    const outlinePrimary = "border border-primary/20 bg-primary/10 text-primary";

    // Notification-like semantic colors (outline style)
    const outlineError = "border border-red-500 bg-red-500/10 text-red-500";
    const outlineWarning = "border border-yellow-500 bg-yellow-500/10 text-yellow-500";
    const outlineInfo = "border border-blue-500 bg-blue-500/10 text-blue-500";

    // For filled, provide sensible semantic defaults
    const filledError = "bg-red-500 text-white";
    const filledWarning = "bg-yellow-500 text-black";
    const filledInfo = "bg-blue-500 text-white";

    const variantClasses = (() => {
      if (variant === "filled") {
        if (xColor === "primary") return filledPrimary;
        if (xColor === "error") return filledError;
        if (xColor === "warning") return filledWarning;
        if (xColor === "info") return filledInfo;
        return filledDefault; // default
      }
      // outline
      if (xColor === "primary") return outlinePrimary;
      if (xColor === "error") return outlineError;
      if (xColor === "warning") return outlineWarning;
      if (xColor === "info") return outlineInfo;
      return outlineDefault; // default
    })();

    const containerClasses = clsx(
      "inline-flex items-center justify-center rounded-md select-none",
      sizeClasses,
      variantClasses,
      loading && "animate-blink2 pointer-events-none",
      className
    );

    return (
      <div
        ref={ref}
        className={containerClasses}
        onClick={onClick}
        aria-busy={loading || undefined}
        {...rest}
      >
        {children}
      </div>
    );
  }
);

Chip.displayName = "Chip";

export default Chip;
