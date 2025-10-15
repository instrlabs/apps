import React from "react";
import clsx from "clsx";

export type MenuIconProps = React.SVGAttributes<SVGSVGElement> & {
  title?: string;
};

export default function MenuIcon({ className, title = "Menu", ...rest }: MenuIconProps) {
  return (
    <svg
      role="img"
      aria-label={title}
      width="28"
      height="28"
      viewBox="0 0 28 28"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      className={clsx("inline-block", className)}
      {...rest}
    >
      <title>{title}</title>
      <path
        d="M0.666748 12.4617V0.666748H12.4617V12.4617H0.666748ZM0.666748 27.3334V15.5384H12.4617V27.3334H0.666748ZM15.5384 12.4617V0.666748H27.3334V12.4617H15.5384ZM15.5384 27.3334V15.5384H27.3334V27.3334H15.5384ZM2.07717 11.0513H11.0513V2.07716H2.07717V11.0513ZM16.9488 11.0513H25.923V2.07716H16.9488V11.0513ZM16.9488 25.923H25.923V16.9488H16.9488V25.923ZM2.07717 25.923H11.0513V16.9488H2.07717V25.923Z"
        fill="currentColor"
      />
    </svg>
  );
}
