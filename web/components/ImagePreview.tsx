"use client";

import React from "react";
import Image from "next/image";
import useModal from "@/hooks/useModal";

export type ImagePreviewProps = {
  src: string;
  alt: string;
  width: number;
  height: number;
  className?: string;
};

function ImagePreviewOverlay({ src }: { src: string }) {
  return (
    <Image
      fill
      alt="Preview"
      src={src}
      quality={100}
      objectFit="contain"
      className="max-h-[90vh] max-w-[90vw] relative!"
    />
  );
}

export default function ImagePreview({ src, alt, width, height, className }: ImagePreviewProps) {
  const { openModal } = useModal();

  return (
    <div
      role="button"
      className={
        "relative aspect-square object-cover cursor-zoom-in rounded-lg overflow-hidden " +
        (className || "")
      }
      onClick={(e) => {
        e.stopPropagation();
        openModal(<ImagePreviewOverlay src={src} />);
      }}
    >
      <Image
        src={src}
        alt={alt}
        width={width}
        height={height}
        objectFit="cover"
        quality={30}
      />
      <div
        className={`
          absolute inset-0 bg-black/20
          opacity-0 hover:opacity-100 transition-opacity
        `}
      />
    </div>
  );
}
