import * as React from "react";

export default function SearchIcon(props: React.SVGProps<SVGSVGElement>) {
  const { className, ...rest } = props;
  return (
    <svg
      className={className}
      viewBox="0 0 40 40"
      fill="none"
      stroke="currentColor"
      strokeWidth="1.5"
      aria-hidden="true"
      {...rest}
    >
      <g clipPath="url(#clip0_231_66)">
        <path
          d="M5 18C5 19.7072 5.33626 21.3977 5.98957 22.9749C6.64288 24.5521 7.60045 25.9852 8.80761 27.1924C10.0148 28.3995 11.4479 29.3571 13.0251 30.0104C14.6023 30.6637 16.2928 31 18 31C19.7072 31 21.3977 30.6637 22.9749 30.0104C24.5521 29.3571 25.9852 28.3995 27.1924 27.1924C28.3995 25.9852 29.3571 24.5521 30.0104 22.9749C30.6637 21.3977 31 19.7072 31 18C31 16.2928 30.6637 14.6023 30.0104 13.0251C29.3571 11.4479 28.3995 10.0148 27.1924 8.80761C25.9852 7.60045 24.5521 6.64288 22.9749 5.98957C21.3977 5.33625 19.7072 5 18 5C16.2928 5 14.6023 5.33625 13.0251 5.98957C11.4479 6.64288 10.0148 7.60045 8.80761 8.80761C7.60045 10.0148 6.64288 11.4479 5.98957 13.0251C5.33626 14.6023 5 16.2928 5 18Z"
          strokeLinecap="round"
          strokeLinejoin="round"
        />
        <path d="M35 35L28 28" strokeLinecap="round" strokeLinejoin="round" />
      </g>
      <defs>
        <clipPath id="clip0_231_66">
          <rect width="40" height="40" fill="white" />
        </clipPath>
      </defs>
    </svg>
  );
}
