"use client";

import React, { useEffect } from "react";
import useModal from "@/hooks/useModal";

export default function ImagePreviewOverlay({ src }: { src?: string | null }) {
  const { closeModal } = useModal();

  // Close on Escape
  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") {
        closeModal();
      }
    };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, [closeModal]);

  if (!src) {
    return (
      <div className="p-6">
        <div className="text-sm text-muted">No image selected.</div>
      </div>
    );
  }

  return (
    <div className="relative">
      <button
        className="absolute top-2 right-2 text-2xl px-2 leading-none"
        aria-label="Close preview"
        onClick={() => {
          closeModal();
        }}
      >
        ×
      </button>
      <div className="p-2">
        <img
          src={src}
          alt="Full-size preview"
          className="max-h-[80vh] max-w-[90vw] object-contain rounded bg-white"
        />
      </div>
    </div>
  );
}
