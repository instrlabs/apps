"use client";

import React from "react";
import Image from "next/image";

export default function ImagePreviewOverlay({ src }: {
  src: string
}) {
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
