import React, { forwardRef } from "react";

interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  xSize?: "sm" | "md" | "lg";
}

const Input = forwardRef<HTMLInputElement, InputProps>(
  ({ xSize = "md", className, ...rest }, ref) => {
    const sizeClasses =
      xSize === "sm" ? "px-2 py-1.5 text-sm" :
      xSize === "md" ? "px-2 py-2 text-base" :
      xSize === "lg" ? "px-3 py-3 text-base" :
      "";

    const baseClasses = `
w-full
border border-primary rounded
bg-secondary text-white placeholder:text-muted
focus:outline-none
focus:border-white focus:bg-white/15
    `;

    return (
      <input
        className={[baseClasses, sizeClasses, className].filter(Boolean).join(" ")}
        ref={ref}
        {...rest}
      />
    );
  },
);

Input.displayName = "Input";

export type { InputProps };
export default Input;
