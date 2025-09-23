import React, { forwardRef } from "react";
import clsx from "clsx";

interface TextFieldProps extends React.InputHTMLAttributes<HTMLInputElement> {
  xSize?: "sm" | "md" | "lg";
  xIsInvalid?: boolean;
  xErrorMessage?: string;
}

const TextField = forwardRef<HTMLInputElement, TextFieldProps>(
  ({ xSize = "md", xIsInvalid, xErrorMessage, className, ...rest }, ref) => {
    const sizeClasses =
      xSize === "sm" ? "px-3 leading-8 text-sm"
      : xSize === "md" ? "px-3 leading-12"
        : xSize === "lg" ? "px-5"
          : "";

    const baseClasses = `
      bg-white/2
      focus:outline-none
      w-full shadow-primary rounded-lg
      placeholder:[color:var(--text-primary)]/60
      hover:shadow-hover
      focus:shadow-focus`;

    const errorClasses = xIsInvalid
      ? "border-[var(--input-error-border)]"
      : undefined;

    return (
      <div className="relative w-full">
        <input
          className={clsx(baseClasses, sizeClasses, errorClasses, className)}
          ref={ref}
          {...rest}
        />
        {xIsInvalid && xErrorMessage ? (
          <span className="text-error absolute top-full left-5 mt-1 text-xs">{xErrorMessage}</span>
        ) : null}
      </div>
    );
  },
);

TextField.displayName = "TextField";

export type { TextFieldProps };
export default TextField;
