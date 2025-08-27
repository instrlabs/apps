import clsx from "clsx";
import React from "react";
import Image from "next/image";

export type AvatarSize = "sm" | "md" | "lg";

export type AvatarProps = {
  name?: string;
  src?: string;
  size: AvatarSize;
  className?: string;
  onClick?: React.MouseEventHandler<HTMLButtonElement>;
};

export default function Avatar({
  name = "User",
  src,
  size = "md",
  className,
  onClick,
}: AvatarProps) {
  const safeName = (name ?? "User").trim() || "User";
  const firstChar = safeName[0]?.toLowerCase() ?? "u";
  const isAlpha = firstChar >= "a" && firstChar <= "z";
  const alphaIndex = isAlpha ? firstChar.charCodeAt(0) - 97 : firstChar.charCodeAt(0) % 26;
  const bucket = ((alphaIndex % 26) + 26) % 8;

  // Tailwind color classes by bucket (no hex values)
  const bgPaletteCls = [
    "bg-blue-500",
    "bg-green-500",
    "bg-red-500",
    "bg-yellow-500",
    "bg-purple-500",
    "bg-teal-500",
    "bg-orange-500",
    "bg-slate-500",
  ];
  const textPaletteCls = [
    "text-white", // on blue
    "text-white", // on green
    "text-white", // on red
    "text-black", // on yellow
    "text-white", // on purple
    "text-white", // on teal
    "text-black", // on orange
    "text-white", // on slate
  ];

  const bgClass = bgPaletteCls[bucket] ?? "bg-blue-500";
  const fgClass = textPaletteCls[bucket] ?? "text-white";

  // Compute two-letter initials (prefer first letters of first two words; fallback to first two letters)
  const words = safeName.split(/\s+/).filter(Boolean);
  let initials = "";
  if (words.length >= 2) {
    initials = `${words[0][0] ?? ""}${words[1][0] ?? ""}`;
  } else {
    const lettersOnly = safeName.replace(/[^A-Za-z]/g, "");
    initials = lettersOnly.slice(0, 2) || safeName.slice(0, 2);
  }
  initials = initials.toUpperCase();

  // Size mappings
  const sizeClass =
    size === "sm" ? "w-8 h-8" : size === "lg" ? "w-14 h-14" : "w-10 h-10"; // md default
  const textSizeClass = size === "sm" ? "text-xs" : size === "lg" ? "text-lg" : "text-sm";

  return (
    <button
      type="button"
      onClick={onClick}
      aria-label={safeName}
      title={safeName}
      className={clsx(
        "relative inline-flex items-center justify-center shadow-primary leading-none rounded-full cursor-pointer overflow-hidden",
        sizeClass,
        bgClass,
        fgClass,
        className
      )}
    >
      {src ? (
        <Image
          src={src}
          alt={safeName}
          fill
          sizes={size === "sm" ? "2rem" : size === "lg" ? "3.5rem" : "2.5rem"}
          className="object-cover rounded-full"
        />
      ) : (
        <div
          className={clsx(
            "w-full h-full flex items-center justify-center font-semibold",
            textSizeClass
          )}
        >
          {initials}
        </div>
      )}
    </button>
  );
}
