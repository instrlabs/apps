import React, { forwardRef } from "react";
import clsx from "clsx";

interface TextFieldProps extends React.InputHTMLAttributes<HTMLInputElement> {
  xSize?: "sm" | "md" | "lg";
}

const TextField = forwardRef<HTMLInputElement, TextFieldProps>(
  ({ xSize = "md", className, ...rest }, ref) => {
    const sizeClasses =
      xSize === "sm" ? "px-2 py-1.5 text-sm" // 8px x, 6px y, 14px font
      : xSize === "md" ? "px-2 py-2 text-base" // 8px all, 16px font
        : xSize === "lg" ? "px-3 py-3 text-base" // 12px all, 16px font
          : "";

    const baseClasses = `
      bg-[rgba(255,255,255,0.08)]
      border border-[rgba(255,255,255,0.3)]
      text-white placeholder:[color:rgba(255,255,255,0.4)]
      focus:outline-none w-full rounded
    `;

    return (
      <input
        className={clsx(baseClasses, sizeClasses, className)}
        ref={ref}
        {...rest}
      />
    );
  },
);

TextField.displayName = "TextField";

export type { TextFieldProps };
export default TextField;
