import React from "react";
import clsx from "clsx";

export type IconButtonProps = React.ButtonHTMLAttributes<HTMLButtonElement> & {
  xSize?: "sm" | "md" | "lg";
  xColor?: "primary" | "secondary";
};

export default function IconButton({
  xSize = "md",
  xColor = "secondary",
  className,
  onClick,
  children,
  "aria-label": ariaLabel,
  type = "button",
  ...rest
}: IconButtonProps) {
  const primaryClasses = [
    "bg-white",
    "text-black",
    "hover:bg-white/85",
    "disabled:bg-white/85",
  ].join(" ");

  const secondaryClasses = [
    "bg-transparent",
    "text-white/90",
    "hover:bg-white/10",
    "disabled:opacity-60",
  ].join(" ");

  const colorClasses = xColor === "primary" ? primaryClasses : secondaryClasses;

  const smClasses = "w-8 h-8";
  const mdClasses = "w-10 h-10";
  const lgClasses = "w-12 h-12";

  const sizeClasses =
    xSize === "sm" ? smClasses : xSize === "md" ? mdClasses : xSize === "lg" ? lgClasses : "";

  const baseClasses =
    "inline-flex items-center justify-center rounded-md cursor-pointer hover:opacity-90 transition-opacity transition-colors focus:outline-none disabled:cursor-not-allowed";

  const handleClick = (e: React.MouseEvent<HTMLButtonElement>) => {
    if (onClick) return onClick(e);
    console.log("IconButton clicked");
  };

  return (
    <button
      type={type}
      aria-label={ariaLabel}
      className={clsx(baseClasses, sizeClasses, colorClasses, className)}
      onClick={handleClick}
      {...rest}
    >
      {children}
    </button>
  );
}
