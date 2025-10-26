"use client";

import React from "react";
import CircleSvg from "./svgs/circle";
import RectangleSvg from "./svgs/rectangle";
import RectangleOutlineSvg from "./svgs/rectangle-outline";
import GoogleSvg from "./svgs/google";
import SearchSvg from "./svgs/search";
import VisibleSvg from "./svgs/visible";
import CircleSuccessSvg from "./svgs/circle-success";
import CircleErrorSvg from "./svgs/circle-error";
import WarningSvg from "./svgs/warning";
import CircleInfoSvg from "./svgs/circle-info";
import LogoSvg from "./svgs/logo";
import CloseSvg from "./svgs/close";
import NotificationSvg from "./svgs/notification";
import ProgressSvg from "./svgs/progress";

type IconProps = {
  name: string;
  size?: number;
  className?: string;
};

const iconMap: Record<string, React.ComponentType<React.SVGProps<SVGSVGElement>>> = {
  "logo": LogoSvg,
  "circle": CircleSvg,
  "rectangle": RectangleSvg,
  "rectangle-outline": RectangleOutlineSvg,
  "google": GoogleSvg,
  "search": SearchSvg,
  "visible": VisibleSvg,
  "circle-success": CircleSuccessSvg,
  "circle-error": CircleErrorSvg,
  "warning": WarningSvg,
  "circle-info": CircleInfoSvg,
  "close": CloseSvg,
  "notification": NotificationSvg,
  "progress": ProgressSvg,
};

export default function Icon({
  name,
  size = 24,
  className = "",
}: IconProps) {
  const IconComponent = iconMap[name];

  return (
    <IconComponent width={size} height={size} className={className}>
    </IconComponent>
  );
}
