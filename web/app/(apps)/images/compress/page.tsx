"use client";

import React, { useCallback, useEffect, useState } from "react";
import FileDropzone from "@/components/file-dropzone";
import useNotification from "@/hooks/useNotification";
import Button from "@/components/button";
import Icon from "@/components/icon";

type CompressState = "default" | "uploaded" | "success";

interface FileMetadata {
  file: File;
  id: string;
  typeLabel: string;
  dimensions?: { width: number; height: number };
}

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
  const { showNotification } = useNotification();
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
          };
        })
      );

      setFileMetadata((prevMetadata) => [...prevMetadata, ...newMetadata]);
      setState("uploaded");
      showNotification({
        type: "info",
        message: `${newFiles.length} file(s) added successfully!`,
      });
    },
    [showNotification],
  );

  const handleRemoveFile = useCallback((id: string) => {
    setFileMetadata((prevMetadata) => prevMetadata.filter((item) => item.id !== id));
  }, []);

  const handleSubmit = useCallback(async () => {
    if (fileMetadata.length === 0) {
      showNotification({
        type: "error",
        message: "Please select files to compress",
      });
      return;
    }

    // TODO: Implement actual compression logic
    // For now, just transition to success state
    setState("success");
  }, [fileMetadata.length, showNotification]);

  const handleReset = useCallback(() => {
    setFileMetadata([]);
    setState("default");
  }, []);

  // Auto-reset success state after 2 seconds
  useEffect(() => {
    if (state === "success") {
      const timer = setTimeout(() => {
        setFileMetadata([]);
        setState("default");
      }, 2000);
      return () => clearTimeout(timer);
    }
  }, [state]);

  return (
    <div className="flex h-full w-full flex-col gap-2">
      <div className="flex h-full flex-col gap-4 rounded-lg border border-white/10 bg-white/8 p-4">
        {/* Title */}
        <h1 className="text-base leading-6 font-semibold text-white">Image Compress</h1>

        {/* Default State - File Dropzone */}
        {state === "default" && (
          <FileDropzone
            accepts={["image/jpeg", "image/png", "image/webp"]}
            multiple
            maxSize={5242880} // 5MB
            title="Upload Files"
            onFilesAdded={handleFilesAdded}
            className="h-full"
          />
        )}

        {/* Uploaded State - File List */}
        {state === "uploaded" && (
          <div className="flex h-full flex-col gap-4">
            {/* File List */}
            {fileMetadata.map((item) => (
              <div
                key={item.id}
                className="flex items-center justify-between rounded border border-white/10 bg-white/4 p-3"
              >
                <div className="flex items-center gap-4">
                  <Icon name="rectangle" size={40} className="text-white/80" />
                  <div className="flex flex-col gap-1">
                    <p className="text-sm leading-5 font-normal text-white">{item.file.name}</p>
                    <p className="text-xs leading-4 font-normal text-white/30">
                      {formatFileMetadata(item.file, item.dimensions)}
                    </p>
                  </div>
                </div>
                <button
                  onClick={() => handleRemoveFile(item.id)}
                  className="text-white/60 transition-colors hover:text-white"
                  aria-label="Remove file"
                >
                  <Icon name="close" size={24} />
                </button>
              </div>
            ))}

            {/* Action Buttons */}
            <div className="flex gap-2">
              <Button
                variant="primary"
                size="base"
                onClick={handleSubmit}
                className="min-w-[150px]"
              >
                Submit
              </Button>
              <Button
                variant="secondary"
                size="base"
                onClick={handleReset}
                className="min-w-[150px]"
              >
                Reset
              </Button>
            </div>
          </div>
        )}

        {/* Success State - Confirmation */}
        {state === "success" && (
          <div className="flex h-full flex-col items-center justify-center gap-4">
            <Icon name="circle-success" size={80} className="text-green-400" />
            <div className="flex flex-col items-center gap-2">
              <h2 className="text-lg leading-7 font-semibold text-white">
                Compression Complete
              </h2>
              <p className="text-sm leading-5 font-normal text-white/60">
                {fileMetadata.length} file(s) compressed successfully
              </p>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
