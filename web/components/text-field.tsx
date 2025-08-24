import React, { forwardRef } from "react";
import clsx from "clsx";

interface TextFieldProps extends React.InputHTMLAttributes<HTMLInputElement> {
  xSize?: "sm" | "md" | "lg";
  xIsInvalid?: boolean;
  xErrorMessage?: string;
  xIsRounded?: boolean;
}

const TextField = forwardRef<HTMLInputElement, TextFieldProps>(
  (
    {
      xSize = "lg",
      xIsInvalid,
      xErrorMessage,
      className,
      xIsRounded,
      ...rest
    },
    ref
  ) => {
    const sizeClasses =
      xSize === "sm" ? "px-3 py-2 text-sm"
        : xSize === "md" ? "px-4 py-3"
          : "px-5 py-4";

    const shapeClass = xIsRounded === true ? "rounded-full" : "rounded-lg";

    const baseClasses =
      "w-full shadow-primary border bg-[var(--input-bg)] " +
      "border-[var(--input-border)] placeholder:[color:var(--input-placeholder)] " +
      "hover:border-[var(--input-hover-border)] focus:border-[var(--input-focus-border)] " +
      "focus:shadow-[var(--input-focus-shadow)] focus:outline-none " +
      "disabled:bg-[var(--input-disabled-bg)] disabled:text-[color:var(--input-disabled-text)] disabled:cursor-not-allowed";

    // Error state overrides using the provided tokens
    const errorClasses = xIsInvalid
      ? "border-[var(--input-error-border)] focus:shadow-[var(--input-error-shadow)]"
      : undefined;

    return (
      <div className="relative w-full">
        <input
          className={clsx(
            baseClasses,
            shapeClass,
            sizeClasses,
            errorClasses,
            className
          )}
          aria-invalid={xIsInvalid || undefined}
          ref={ref}
          {...rest}
        />
        {xIsInvalid && xErrorMessage ? (
          <span className="absolute left-5 top-full mt-1 text-xs text-error">
            {xErrorMessage}
          </span>
        ) : null}
      </div>
    );
  }
);

TextField.displayName = "TextField";

export type { TextFieldProps };
export default TextField;
