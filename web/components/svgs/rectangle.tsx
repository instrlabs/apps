import React from "react";

export default function RectangleSvg(props: React.SVGProps<SVGSVGElement>) {
  return (
    <svg
      viewBox="0 0 24 24"
      xmlns="http://www.w3.org/2000/svg"
      {...props}
    >
      <rect x="0" y="0" width="24" height="24" rx="4" fill="currentColor" />
    </svg>
  );
}
