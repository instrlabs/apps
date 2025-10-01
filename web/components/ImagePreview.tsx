"use client";

import React from "react";
import Image from "next/image";
import useModal from "@/hooks/useModal";

export type ImagePreviewProps = {
  src: string;
  alt: string;
  width: number;
  height: number;
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

export default function ImagePreview({ src, alt, width, height }: ImagePreviewProps) {
  const { openModal } = useModal();

  return (
    <div
      role="button"
      className={
        "relative cursor-zoom-in overflow-hidden flex items-center justify-center shadow-primary rounded-lg"
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
        quality={30}
        className="aspect-square object-contain object-center"
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
