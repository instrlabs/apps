import React from "react";
import clsx from "clsx";

interface OutlinedButtonProps {
  icon?: React.ReactNode;
  children: React.ReactNode;
  onClick?: () => void;
  className?: string;
}

const OutlinedButton: React.FC<OutlinedButtonProps> = ({ icon, children, onClick, className }) => {
  return (
    <button
      onClick={onClick}
      className={clsx(
        "flex items-center space-x-1",
        "border border-solid border-black/[.08] rounded-md",
        "px-2.5 py-1",
        "hover:bg-[#f2f2f2] hover:border-transparent transition",
        className
      )}
    >
      {icon}
      <span className="text-black text-sm">{children}</span>
    </button>
  );
};

export default OutlinedButton;
