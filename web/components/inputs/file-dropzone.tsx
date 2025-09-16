"use client";

import React, { useCallback, useMemo, useRef, useState } from "react";
import { bytesToString } from "@/utils/bytesToString";
import { acceptsToExtensions } from "@/utils/acceptsToExtensions";
import useNotification from "@/hooks/useNotification";

export type FileDropzoneProps = {
  accepts: string[];
  multiple: boolean;
  onFilesAdded: (files: File[]) => void;
  maxFileSize: number; // in bytes
};

function validateFiles(files: File[], accepts: string[], maxFileSize: number): boolean {
  const acceptedTypes = accepts.map((s) => s.trim()).filter(Boolean);
  return files.every((file) => acceptedTypes.includes(file.type) && file.size <= maxFileSize);
}

const FileDropzone: React.FC<FileDropzoneProps> = ({ accepts, onFilesAdded, multiple, maxFileSize }) => {
  const { showNotification } = useNotification();
  const [isDragging, setIsDragging] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);

  const openFileDialog = useCallback(() => {
    inputRef.current?.click();
  }, []);

  const handleDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    if (!isDragging) setIsDragging(true);
  }, [isDragging]);

  const handleDragLeave = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(false);
  }, []);

  const handleDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(false);
    if (e.dataTransfer?.files && e.dataTransfer.files.length > 0) {
      const files = Array.from(e.dataTransfer.files);

      if (!multiple && files.length > 1) {
        showNotification({ title: "Error", message: "Only one file can be dropped at a time.", type: "error" });
        return;
      }

      if (!validateFiles(files, accepts, maxFileSize)) {
        showNotification({ title: "Error", message: "Invalid file type or file size.", type: "error" });
        return;
      }

      onFilesAdded(files);
      e.dataTransfer.clearData();
    }
  }, [multiple, accepts, maxFileSize, onFilesAdded, showNotification]);

  const baseClass = useMemo(() => (
    `w-full max-w-2xl aspect-video flex flex-col items-center justify-center gap-3 ` +
    `border-1 border-dashed rounded-xl p-10 cursor-pointer ` +
    (isDragging ? "bg-gray-50" : "")
  ), [isDragging]);

  return (
    <div
      role="button"
      onClick={openFileDialog}
      onDragOver={handleDragOver}
      onDragEnter={handleDragOver}
      onDragLeave={handleDragLeave}
      onDrop={handleDrop}
      className={baseClass}
    >
      <div className="text-center gap-1">
        <p className="text-base font-light">Maximum file size: {bytesToString(maxFileSize)}</p>
        <p className="text-base font-light">Support format: {acceptsToExtensions(accepts).join(", ")}</p>
      </div>
      <input
        ref={inputRef}
        className="hidden"
        type="file"
        accept={accepts.join(",")}
        multiple={multiple}
        onChange={(e) => {
          if (e.target.files) {
            const files = Array.from(e.target.files);

            if (!multiple && files.length > 1) {
              showNotification({
                type: "error",
                title: "Something went wrong",
                message: "Only one file can be dropped at a time."
              });
            } else if (!validateFiles(files, accepts, maxFileSize)) {
              showNotification({
                type: "error",
                title: "Something went wrong",
                message: "Invalid file type or file size."
              });
            } else {
              onFilesAdded(files);
              return;
            }
          }

          e.currentTarget.value = "";
        }}
      />
    </div>
  );
};

export default FileDropzone;
