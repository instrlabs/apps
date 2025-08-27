import clsx from "clsx";
import React from "react";

export type AvatarProps = {
  name?: string; // Used for fallback initials/avatar service
  src?: string; // If provided, used as image source instead of generated
  alt?: string; // Accessible alt text; defaults to `${name} avatar` or "Profile avatar"
  size?: number; // Pixel size for width/height
  className?: string; // Extra classes to merge
  rounded?: boolean; // Whether to make it fully rounded
};

/**
 * Reusable Avatar component with sensible defaults.
 * - If `src` is not provided, it falls back to ui-avatars.com using the `name`.
 * - Controls width/height via the `size` prop (default 40px).
 * - Applies `rounded-full` and `object-cover` by default for consistency.
 */
export default function Avatar({
  name = "User",
  src,
  alt,
  size = 40,
  className,
  rounded = true,
}: AvatarProps) {
  // Choose one of 8 background colors based on the first letter (A–Z) bucketed by modulo 8
  const safeName = (name ?? "User").trim() || "User";
  const firstChar = safeName[0]?.toLowerCase() ?? "u";
  const isAlpha = firstChar >= "a" && firstChar <= "z";
  const alphaIndex = isAlpha ? firstChar.charCodeAt(0) - 97 : (firstChar.charCodeAt(0) % 26);
  const bucket = ((alphaIndex % 26) + 26) % 8; // ensure 0-7

  // 8 distinct, readable colors (Tailwind-inspired hex, without the leading '#')
  const bgPalette = [
    "3b82f6", // blue-500
    "22c55e", // green-500
    "ef4444", // red-500
    "eab308", // yellow-500
    "a855f7", // purple-500
    "14b8a6", // teal-500
    "f97316", // orange-500
    "64748b", // slate-500
  ];
  // Pick contrasting text color; for lighter backgrounds prefer black text
  const textPalette = [
    "ffffff", // on blue
    "ffffff", // on green
    "ffffff", // on red
    "000000", // on yellow
    "ffffff", // on purple
    "ffffff", // on teal
    "000000", // on orange
    "ffffff", // on slate
  ];

  const background = bgPalette[bucket] ?? "3b82f6";
  const textColor = textPalette[bucket] ?? "ffffff";

  // Let ui-avatars create initials from the full name; ensure proper encoding
  const initialsName = encodeURIComponent(safeName);
  const fallbackUrl = `https://ui-avatars.com/api/?name=${initialsName}&background=${background}&color=${textColor}`;

  const finalSrc = src ?? fallbackUrl;
  const finalAlt = alt ?? (safeName ? `${safeName} avatar` : "Profile avatar");

  return (
    <img
      src={finalSrc}
      alt={finalAlt}
      width={size}
      height={size}
      className={clsx(
        rounded && "rounded-full",
        "object-cover shadow-primary",
        className
      )}
    />
  );
}
