"use client";

import React from "react";
import CircleErrorSvg from "@/components/svgs/circle-error";
import CircleSuccessSvg from "@/components/svgs/circle-success";
import CircleInfoSvg from "@/components/svgs/circle-info";
import WarningSvg from "@/components/svgs/warning";

type IconType = "error" | "warning" | "info" | "success";

type SnackbarIconProps = React.SVGProps<SVGSVGElement> & {
  type?: IconType;
};

export default function SnackbarIcon({
  type = "success",
  className = "w-5 h-5",
  ...props
}: SnackbarIconProps) {
  const iconProps = { className, ...props };

  switch (type) {
    case "error":
      return <CircleErrorSvg {...iconProps} />;
    case "warning":
      return <WarningSvg {...iconProps} />;
    case "info":
      return <CircleInfoSvg {...iconProps} />;
    case "success":
    default:
      return <CircleSuccessSvg {...iconProps} />;
  }
}
