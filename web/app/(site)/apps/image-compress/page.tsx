"use client";

import React, {useCallback, useRef, useState, useEffect} from "react";
import HashtagIcon from "@/components/icons/hashtag";

type Preview = { url: string; name: string };

export default function ImageCompressPage() {
  const [previews, setPreviews] = useState<Preview[]>([]);
  const urlsRef = useRef<string[]>([]);

  const handleFiles = useCallback((files: File[]) => {
    if (!files || files.length === 0) return;
    const accepted = files.filter((f) => {
      const typeOk = ["image/png", "image/jpeg", "image/jpg"].includes(f.type);
      const extOk = /\.(png|jpe?g)$/i.test(f.name);
      return typeOk || extOk;
    });

    if (accepted.length === 0) return;

    const newPreviews: Preview[] = accepted.map((f) => {
      const url = URL.createObjectURL(f);
      urlsRef.current.push(url);
      return { url, name: f.name };
    });

    setPreviews((prev) => [...prev, ...newPreviews]);
  }, []);

  const removePreview = useCallback((url: string) => {
    try {
      URL.revokeObjectURL(url);
    } catch {}
    urlsRef.current = urlsRef.current.filter((u) => u !== url);
    setPreviews((prev) => prev.filter((p) => p.url !== url));
  }, []);

  // Cleanup object URLs on unmount
  useEffect(() => {
    return () => {
      urlsRef.current.forEach((u) => URL.revokeObjectURL(u));
    };
  }, []);

  return (
    <div className="w-full h-full flex flex-col">
      <div className="p-6 flex items-center gap-3">
        <HashtagIcon className="w-5 h-5" />
        <h1 className="text-xl font-bold">Image Compress</h1>
      </div>
      <div className="flex-1 p-6">
        {previews.length === 0 && (
          <div
            className="w-full h-80 border-2 border-dashed border-gray-300 rounded-lg flex flex-col items-center justify-center gap-4 cursor-pointer hover:border-gray-400 transition-colors"
            onDragOver={(e) => {
              e.preventDefault();
              e.stopPropagation();
            }}
            onDrop={(e) => {
              e.preventDefault();
              e.stopPropagation();
              const files = Array.from(e.dataTransfer.files || []);
              handleFiles(files);
            }}
            onClick={() => document.getElementById('fileInput')?.click()}
          >
            <p className="text-lg text-gray-600">Drag and drop images here</p>
            <p className="text-sm text-gray-500">or click to select files</p>
            <input
              id="fileInput"
              type="file"
              accept="image/png, image/jpeg, image/jpg"
              multiple
              className="hidden"
              onChange={(e) => {
                const files = Array.from(e.target.files || []);
                handleFiles(files);
              }}
            />
          </div>
        )}

        {previews.length > 0 && (
          <div className="mt-6 grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-4">
            {previews.map((p, idx) => (
              <div key={`${p.url}-${idx}`} className="relative group w-full aspect-square bg-gray-100 rounded overflow-hidden">
                {/* eslint-disable-next-line @next/next/no-img-element */}
                <img src={p.url} alt={p.name} className="w-full h-full object-cover" />
                <button
                  type="button"
                  aria-label="Delete image"
                  onClick={() => removePreview(p.url)}
                  className="absolute top-1.5 right-1.5 rounded-md bg-white/80 text-gray-700 hover:bg-white shadow-sm px-2 py-1 text-xs font-medium transition-colors"
                >
                  ×
                </button>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
