import clsx from "clsx";
import React from "react";
import Image from "next/image";

export type AvatarSize = "sm" | "md" | "lg";

export type AvatarProps = {
  name?: string;
  src?: string;
  xsize: AvatarSize;
  className?: string;
  onClick?: React.MouseEventHandler<HTMLButtonElement>;
};

function getBucket(name: string | undefined) {
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

function getInitial(name: string | undefined) {
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

export default function Avatar({
  name = "",
  xsize = "md",
  onClick,
}: AvatarProps) {
  const bucket = getBucket(name);
  const initials = getInitial(name);

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

  const bgClass = bgPaletteCls[bucket];
  const fgClass = textPaletteCls[bucket];

  // Size Classes
  const smClasses = "size-8 text-sm";
  const mdClasses = "size-10";
  const lgClasses = "size-12";

  const sizeClass =
    xsize === "sm" ? smClasses :
      xsize === "md" ? mdClasses :
        xsize === "lg" ? lgClasses :
          "";

  const baseClasses = `
    flex items-center justify-center rounded-full
    cursor-pointer select-none
  `;

  return (
    <button
      type="button"
      onClick={onClick}
      className={clsx(
        baseClasses,
        sizeClass,
        bgClass,
        fgClass,
      )}
    >
      {initials}
    </button>
  );
}
