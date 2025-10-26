"use client";

import React, { useCallback, useState } from "react";
import useSnackbar from "@/hooks/useSnackbar";
import DefaultState from "./DefaultState";
import UploadedState from "./UploadedState";
import { CompressState, FileStatus, FileMetadata } from "./types";

// Helper to get file type label from MIME type
const getFileTypeLabel = (mimeType: string): string => {
  const typeMap: Record<string, string> = {
    "image/jpeg": "JPEG",
    "image/png": "PNG",
    "image/webp": "WebP",
  };
  return typeMap[mimeType] || mimeType.split("/")[1]?.toUpperCase() || "Unknown";
};

// Helper to get image dimensions
const getImageDimensions = (
  file: File
): Promise<{ width: number; height: number } | null> => {
  return new Promise((resolve) => {
    const reader = new FileReader();
    reader.onload = (e) => {
      const img = new Image();
      img.onload = () => {
        resolve({ width: img.width, height: img.height });
      };
      img.onerror = () => resolve(null);
      img.src = e.target?.result as string;
    };
    reader.onerror = () => resolve(null);
    reader.readAsDataURL(file);
  });
};

// Helper to format file metadata
const formatFileMetadata = (
  file: File,
  dimensions?: { width: number; height: number }
): string => {
  const sizeInKB = (file.size / 1024).toFixed(0);
  const typeLabel = getFileTypeLabel(file.type);
  const dimensionsStr = dimensions ? `${dimensions.width}x${dimensions.height}` : "Unknown";
  return `${sizeInKB}KB · ${typeLabel} · ${dimensionsStr}`;
};


export default function ImageCompress() {
  const { showSnackbar } = useSnackbar();
  const [state, setState] = useState<CompressState>("default");
  const [fileMetadata, setFileMetadata] = useState<FileMetadata[]>([]);

  const handleFilesAdded = useCallback(
    async (newFiles: File[]) => {
      // Create metadata for each file
      const newMetadata = await Promise.all(
        newFiles.map(async (file) => {
          const dimensions = await getImageDimensions(file);
          return {
            file,
            id: `${file.name}-${Date.now()}-${Math.random()}`,
            typeLabel: getFileTypeLabel(file.type),
            dimensions: dimensions || undefined,
            status: "waiting" as FileStatus,
          };
        })
      );

      setFileMetadata((prevMetadata) => [...prevMetadata, ...newMetadata]);
      setState("uploaded");
      showSnackbar({
        type: "info",
        message: `${newFiles.length} file(s) added successfully!`,
      });
    },
    [showSnackbar],
  );

  const handleRemoveFile = useCallback((id: string) => {
    setFileMetadata((prevMetadata) => prevMetadata.filter((item) => item.id !== id));
  }, []);

  const handleSubmit = useCallback(async () => {
    if (fileMetadata.length === 0) {
      showSnackbar({
        type: "error",
        message: "Please select files to compress",
      });
      return;
    }

    // Simulate compression for each file
    for (let i = 0; i < fileMetadata.length; i++) {
      // Set status to processing
      setFileMetadata((prev) =>
        prev.map((item, idx) =>
          idx === i ? { ...item, status: "processing" as FileStatus } : item
        )
      );

      // Simulate compression delay (1-3 seconds per file)
      await new Promise((resolve) => setTimeout(resolve, Math.random() * 2000 + 1000));

      // Set status to success with compression data
      setFileMetadata((prev) =>
        prev.map((item, idx) => {
          if (idx !== i) return item;

          // Simulate 40-60% size reduction
          const compressionPercentage = Math.floor(Math.random() * 20 + 40);
          const compressedSize = Math.floor(
            item.file.size * ((100 - compressionPercentage) / 100)
          );

          return {
            ...item,
            status: "success" as FileStatus,
            compressionPercentage,
            compressedSize,
            downloadUrl: `/api/download/${item.id}`,
          };
        })
      );
    }

    showSnackbar({
      type: "info",
      message: "All files compressed successfully!",
    });
  }, [fileMetadata, showSnackbar]);

  const handleReset = useCallback(() => {
    setFileMetadata([]);
    setState("default");
  }, []);

  return (
    <div className="flex h-full w-full flex-col gap-2">
      <div className="flex h-full flex-col gap-4 rounded-lg border border-white/10 bg-white/8 p-4">
        {/* Title */}
        <h1 className="text-base leading-6 font-semibold text-white">Image Compress</h1>

        {/* Default State */}
        {state === "default" && <DefaultState onFilesAdded={handleFilesAdded} />}

        {/* Uploaded State */}
        {state === "uploaded" && (
          <UploadedState
            fileMetadata={fileMetadata}
            onRemoveFile={handleRemoveFile}
            onSubmit={handleSubmit}
            onReset={handleReset}
            onDownload={(fileName) => {
              showSnackbar({
                type: "info",
                message: `Download started for ${fileName}!`,
              });
            }}
            formatFileMetadata={formatFileMetadata}
          />
        )}
      </div>
    </div>
  );
}
