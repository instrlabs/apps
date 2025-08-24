import React from "react";
import clsx from "clsx";

export type ChipSize = "sm" | "md" | "lg";
export type ChipVariant = "filled" | "outlined";

export interface ChipProps extends React.HTMLAttributes<HTMLDivElement> {
  children?: React.ReactNode;
  onClick?: (e: React.MouseEvent<HTMLDivElement>) => void;
  xSize?: ChipSize;
  xVariant?: ChipVariant;
}

const Chip = React.forwardRef<HTMLDivElement, ChipProps>(
  (
    {
      onClick,
      xSize = "md",
      xVariant = "filled",
      className,
      children,
      ...rest
    },
    ref
  ) => {

    const sizeClasses =
      xSize === "sm"
        ? "h-6 text-xs gap-1.5 px-2"
        : xSize === "lg"
        ? "h-9 text-base gap-2.5 px-4"
        : "h-8 text-sm gap-2 px-3";

    const baseFilled = "bg-muted text-[color:var(--fg,inherit)]";
    const baseOutlined = "border border-border bg-transparent text-[color:var(--fg,inherit)]";

    const containerClasses = clsx(
      "inline-flex items-center rounded-full align-middle select-none",
      sizeClasses,
      xVariant === "outlined" ? baseOutlined : baseFilled,
      className
    );

    return (
      <div ref={ref} className={containerClasses} onClick={onClick} {...rest}>
        {children}
      </div>
    );
  }
);

Chip.displayName = "Chip";

export default Chip;
