import React from "react";

export default function NotificationSvg(props: React.SVGProps<SVGSVGElement>) {
  return (
    <svg
      viewBox="0 0 24 24"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      {...props}
    >
      <path
        d="M12 2C10.9 2 10 2.9 10 4V8C10 10.5 8.5 12.8 6.5 14V16H17.5V14C15.5 12.8 14 10.5 14 8V4C14 2.9 13.1 2 12 2Z"
        stroke="currentColor"
        strokeWidth="2"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      <path
        d="M10 20H14C14 21.1 13.1 22 12 22C10.9 22 10 21.1 10 20Z"
        stroke="currentColor"
        strokeWidth="2"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  );
}
