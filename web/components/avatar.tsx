"use client";

import React from "react";

export type AvatarSize = "sm" | "base" | "lg" | "xl";

export type AvatarProps = React.ButtonHTMLAttributes<HTMLButtonElement> & {
  name?: string;
  src?: string;
  size?: AvatarSize;
};

/**
 * Determines color bucket (0-7) based on first character of name
 */
function getBucket(name: string | undefined): number {
  let bucket = 0;
  if (name && name.trim() !== "") {
    const safeName = name.trim();
    const firstChar = safeName[0]?.toLowerCase() ?? "u";
    const isAlpha = firstChar >= "a" && firstChar <= "z";
    const alphaIndex = isAlpha ? firstChar.charCodeAt(0) - 97 : firstChar.charCodeAt(0) % 26;
    bucket = ((alphaIndex % 26) + 26) % 8;
  }
  return bucket;
}

/**
 * Extracts initials from name (max 2 characters)
 */
function getInitials(name: string | undefined): string {
  let initials = "";
  if (name && name.trim() !== "") {
    const safeName = name.trim();
    const words = safeName.split(/\s+/).filter(Boolean);
    if (words.length >= 2) {
      initials = `${words[0][0] ?? ""}${words[1][0] ?? ""}`;
    } else {
      const lettersOnly = safeName.replace(/[^A-Za-z]/g, "");
      initials = lettersOnly.slice(0, 2) || safeName.slice(0, 2);
    }
    initials = initials.toUpperCase();
  }
  return initials;
}

/**
 * Avatar component with size variants and color bucketing
 *
 * Based on Figma design at node 347-1628
 * Sizes: sm (36px), base (40px), lg (48px), xl (80px)
 * Colors: 8 color buckets based on name's first character
 */
export default function Avatar({
  name = "",
  size = "sm",
  className = "",
  ...rest
}: AvatarProps) {
  const bucket = getBucket(name);
  const initials = getInitials(name);

  // Color palette: 8 background colors
  const bgPalette = [
    "bg-blue-500",
    "bg-green-500",
    "bg-red-500",
    "bg-yellow-500",
    "bg-purple-500",
    "bg-teal-500",
    "bg-orange-500",
    "bg-slate-500",
  ];

  // Text colors corresponding to each background
  const textPalette = [
    "text-white",  // on blue
    "text-white",  // on green
    "text-white",  // on red
    "text-black",  // on yellow
    "text-white",  // on purple
    "text-white",  // on teal
    "text-black",  // on orange
    "text-white",  // on slate
  ];

  const bgClass = bgPalette[bucket];
  const textClass = textPalette[bucket];

  const baseClasses = "flex items-center justify-center rounded-full select-none cursor-pointer transition-colors";

  // Size configuration matching Figma specs
  const sizeConfig: Record<AvatarSize, string> = {
    sm: "size-9 text-base leading-6 font-medium",      // 36px, 16px text
    base: "size-10 text-base leading-6 font-medium",   // 40px, 16px text
    lg: "size-12 text-xl leading-7 font-medium",       // 48px, 20px text
    xl: "size-20 text-4xl leading-10 font-medium",     // 80px, 36px text
  };

  return (
    <button
      type="button"
      className={[
        baseClasses,
        sizeConfig[size],
        bgClass,
        textClass,
        className,
      ]
        .filter(Boolean)
        .join(" ")}
      {...rest}
    >
      {initials}
    </button>
  );
}
