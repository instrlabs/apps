import React from "react";
import clsx from "clsx";

interface SubmitButtonProps {
  type?: "button" | "submit";
  /**
   * Deprecated: external loading control. If provided, it overrides internal loading state.
   */
  isLoading?: boolean;
  loadingText?: string;
  children: React.ReactNode;
  onClick?: (event: React.MouseEvent<HTMLButtonElement>) => void | Promise<void>;
}

const Button: React.FC<SubmitButtonProps> = ({
  type = "button",
  children,
  onClick,
}) => {
  const [internalLoading, setInternalLoading] = React.useState(false);

  const handleClick = async (e: React.MouseEvent<HTMLButtonElement>) => {
    if (!onClick) return;
    if (internalLoading) return;

    try {
      const maybePromise = onClick(e);
      const isPromise = !!maybePromise && typeof (maybePromise as PromiseLike<unknown>).then === "function";
      if (isPromise && internalLoading) {
        setInternalLoading(true);
        try {
          await (maybePromise as Promise<void>);
        } finally {
          setInternalLoading(false);
        }
      }
    } catch (err) { throw err; }
  };

  return (
    <button
      type={type}
      className={clsx(
        "bg-[var(--btn-primary-bg)]",
        "text-[var(--btn-primary-text)]",
        "hover:bg-[var(--btn-primary-hover)]",
        "active:bg-[var(--btn-primary-active)]",
        "disabled:bg-[var(--btn-primary-disabled)]",

        "border border-[var(--btn-border)]",
        "py-4 rounded-xl cursor-pointer",
        "font-medium shadow-primary",
        internalLoading && "opacity-70 cursor-not-allowed"
      )}
      disabled={internalLoading}
      aria-busy={internalLoading || undefined}
      onClick={handleClick}
    >
      {children}
    </button>
  );
};

export default Button;
