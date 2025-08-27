"use client";

import { useOverlay } from "@/hooks/useOverlay";

export default function OverlayModal() {
  const { isModalOpen, modalNode, modalKey, closeAll, modalWidth } = useOverlay();

  if (!isModalOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      <div
        className="absolute inset-0 bg-black/30"
        aria-hidden="true"
        onClick={closeAll}
      />
      <div
        role="dialog"
        aria-modal="true"
        aria-label="Modal dialog"
        className="relative z-10 w-full mx-4"
        style={{ maxWidth: modalWidth ? modalWidth : 0 }}
      >
        <div
          key={modalKey}
          className="rounded-2xl bg-card shadow-xl overflow-auto animate-fade-in"
        >
          {modalNode}
        </div>
      </div>
    </div>
  );
}
