import React, { forwardRef } from "react";
import type { InputHTMLAttributes, ReactNode } from "react";
import Input from "./input";

export interface InputWithIconProps extends InputHTMLAttributes<HTMLInputElement> {
  xSize?: "sm" | "md" | "lg";
  icon: ReactNode;
}

const paddingsBySize = {
  sm: { inputPadLeft: "pl-8", iconSize: "w-5 h-5", iconLeft: "left-2" },
  md: { inputPadLeft: "pl-10", iconSize: "w-6 h-6", iconLeft: "left-2" },
  lg: { inputPadLeft: "pl-10", iconSize: "w-6 h-6", iconLeft: "left-2" },
} as const;

const InputWithIcon = forwardRef<HTMLInputElement, InputWithIconProps>(
  ({ xSize = "md", icon, className, ...rest }, ref) => {
    const spec = paddingsBySize[xSize];
    return (
      <div className="relative">
        <div className={`absolute inset-y-0 ${spec.iconLeft} pointer-events-none flex items-center`}>
          <div className={`flex items-center justify-center text-white/80 ${spec.iconSize}`}>
            {icon}
          </div>
        </div>
        <Input
          ref={ref}
          xSize={xSize}
          className={[spec.inputPadLeft, className].filter(Boolean).join(" ")}
          {...rest}
        />
      </div>
    );
  },
);

InputWithIcon.displayName = "InputWithIcon";
export default InputWithIcon;
