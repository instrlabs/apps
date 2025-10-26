"use client";

import React, { useState, useEffect } from "react";
import Image from "next/image";
import useModal from "@/hooks/useModal";

export type ImagePreviewProps = {
  src: string;
  alt: string;
  width?: number;
  height?: number;
  size?: number;
};

function ImagePreviewOverlay({ src, alt }: { src: string; alt: string }) {
  const [dimensions, setDimensions] = useState({ width: 0, height: 0 });

  useEffect(() => {
    const img = new window.Image();
    img.onload = () => {
      setDimensions({ width: img.naturalWidth, height: img.naturalHeight });
    };
    img.onerror = () => {
      setDimensions({ width: 0, height: 0 });
    };

    img.src = src;
  }, [src]);

  return (
    <div className="relative flex h-full w-full items-center justify-center rounded-lg bg-black/80 p-4">
      <Image
        src={src}
        alt={alt}
        width={dimensions.width}
        height={dimensions.height}
        quality={100}
        unoptimized
        className="h-auto max-h-[90vh] w-auto max-w-[90vw] object-contain"
      />
    </div>
  );
}

export default function ImagePreview({
  src,
  alt,
  width = 40,
  height = 40,
  size = 40,
}: ImagePreviewProps) {
  const { openModal } = useModal();

  return (
    <div
      role="button"
      tabIndex={0}
      className="relative flex cursor-zoom-in items-center justify-center overflow-hidden rounded border border-white/10 bg-white/4 transition-transform hover:scale-110"
      onClick={(e) => {
        e.stopPropagation();
        openModal(<ImagePreviewOverlay src={src} alt={alt} />);
      }}
      onKeyDown={(e) => {
        if (e.key === "Enter" || e.key === " ") {
          e.stopPropagation();
          openModal(<ImagePreviewOverlay src={src} alt={alt} />);
        }
      }}
      style={{ width: `${size}px`, height: `${size}px` }}
    >
      <Image
        src={src}
        alt={alt}
        width={width}
        height={height}
        quality={30}
        className="object-contain object-center"
      />
      <div className="absolute inset-0 bg-black/20 opacity-0 transition-opacity hover:opacity-100" />
    </div>
  );
}
