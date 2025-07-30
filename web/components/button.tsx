import React from "react";
import clsx from "clsx";

interface SubmitButtonProps {
  type?: "button" | "submit";
  isLoading?: boolean;
  loadingText?: string;
  children: React.ReactNode;
  className?: string;
  onClick?: () => void;
}

const Button: React.FC<SubmitButtonProps> = ({
  type = "button",
  isLoading = false,
  loadingText = "Loading...",
  children,
  className,
  onClick,
}) => {
  return (
    <button
      type={type}
      className={clsx(
        "py-3 mt-3 rounded cursor-pointer",
        "text-sm font-medium bg-black hover:bg-gray-800 text-white",
        isLoading && "opacity-70 cursor-not-allowed",
        className
      )}
      disabled={isLoading}
      onClick={onClick}
    >
      {isLoading ? loadingText : children}
    </button>
  );
};

export default Button;
