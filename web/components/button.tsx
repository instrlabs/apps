import React from "react";
import clsx from "clsx";

interface SubmitButtonProps {
  type?: "button" | "submit";
  isLoading?: boolean;
  loadingText?: string;
  children: React.ReactNode;
  onClick?: () => void;
}

const Button: React.FC<SubmitButtonProps> = ({
  type = "button",
  isLoading = false,
  loadingText = "Loading...",
  children,
  onClick,
}) => {
  return (
    <button
      type={type}
      className={clsx(
        "py-4 rounded-xl cursor-pointer",
        "font-bold text-white bg-blue-500 hover:bg-blue-400 shadow-primary",
        isLoading && "opacity-70 cursor-not-allowed"
      )}
      disabled={isLoading}
      onClick={onClick}
    >
      {isLoading ? loadingText : children}
    </button>
  );
};

export default Button;
