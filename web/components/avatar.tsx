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
  const fallbackUrl = `https://ui-avatars.com/api/?name=${encodeURIComponent(name)}`;
  const finalSrc = src ?? fallbackUrl;
  const finalAlt = alt ?? (name ? `${name} avatar` : "Profile avatar");

  return (
    <img
      src={finalSrc}
      alt={finalAlt}
      width={size}
      height={size}
      className={clsx(
        rounded && "rounded-full",
        "object-cover",
        className
      )}
    />
  );
}
