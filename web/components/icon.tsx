"use client";

import React from "react";
import CircleSvg from "./svgs/circle";
import RectangleSvg from "./svgs/rectangle";
import GoogleSvg from "./svgs/google";
import SearchSvg from "./svgs/search";
import VisibleSvg from "./svgs/visible";
import CircleSuccessSvg from "./svgs/circle-success";
import CircleErrorSvg from "./svgs/circle-error";
import WarningSvg from "./svgs/warning";
import CircleInfoSvg from "./svgs/circle-info";
import LogoSvg from "./svgs/logo";

type IconProps = {
  name: string;
  size?: number;
  className?: string;
  title?: string;
};

const iconMap: Record<string, React.ComponentType<React.SVGProps<SVGSVGElement>>> = {
  circle: CircleSvg,
  rectangle: RectangleSvg,
  google: GoogleSvg,
  search: SearchSvg,
  visible: VisibleSvg,
  "circle-success": CircleSuccessSvg,
  "circle-error": CircleErrorSvg,
  warning: WarningSvg,
  "circle-info": CircleInfoSvg,
  logo: LogoSvg,
};

export default function Icon({
  name,
  size = 24,
  className = "",
  title,
}: IconProps) {
  const IconComponent = iconMap[name];

  return (
    <IconComponent
      width={size}
      height={size}
      className={className}
      aria-hidden={title ? undefined : true}
      role={title ? "img" : "presentation"}
    >
      {title ? <title>{title}</title> : null}
    </IconComponent>
  );
}
