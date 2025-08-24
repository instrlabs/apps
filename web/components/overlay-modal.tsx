"use client";

import { useEffect, useRef } from "react";
import { useOverlay } from "@/hooks/useOverlay";
import Button from "@/components/button";

export default function OverlayModal() {
  const { isModalOpen, modalNode, modalContentKey } = useOverlay();
  const backdropRef = useRef<HTMLDivElement | null>(null);

  if (!isModalOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      {/* Backdrop */}
      <div
        ref={backdropRef}
        className="absolute inset-0 bg-black/40"
        aria-hidden="true"
        onClick={() => {}}
      />
      {/* Dialog */}
      <div role="dialog" aria-modal="true" aria-label="Modal dialog" className="relative z-10 w-full max-w-2xl mx-4">
        <div className="rounded-2xl bg-card shadow-xl ring-1 ring-foreground/5 overflow-hidden">
          <div key={modalContentKey} className="max-h-[70vh] overflow-auto p-4 animate-fade-in">
            {modalNode}
          </div>
          <div className="px-4 py-3 border-t border-border flex justify-end">
            <Button type="button">Close</Button>
          </div>
        </div>
      </div>
    </div>
  );
}
