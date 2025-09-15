"use client";

import React, { useCallback, useMemo, useRef, useState } from "react";

export type FileDropzoneProps = {
  accept?: string; // comma-separated accept string, e.g. "image/png,image/jpeg"
  onFilesAdded: (files: File[]) => void;
  className?: string;
  children?: React.ReactNode;
  allowMultiple?: boolean; // when false, only one file can be selected/dropped
};

function filterAccepted(files: FileList | File[], accept?: string): File[] {
  const list = Array.from(files);
  if (!accept) return list;
  const acceptedTypes = accept.split(",").map((s) => s.trim()).filter(Boolean);
  if (acceptedTypes.length === 0) return list;
  return list.filter((f) => acceptedTypes.includes(f.type));
}

const FileDropzone: React.FC<FileDropzoneProps> = ({ accept, onFilesAdded, className, children, allowMultiple = false }) => {
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
      let accepted = filterAccepted(e.dataTransfer.files, accept);
      if (!allowMultiple) accepted = accepted.slice(0, 1);
      if (accepted.length > 0) onFilesAdded(accepted);
      e.dataTransfer.clearData();
    }
  }, [accept, allowMultiple, onFilesAdded]);

  const baseClass = useMemo(() => (
    `w-full max-w-2xl aspect-video flex flex-col items-center justify-center gap-3 ` +
    `border-1 border-dashed rounded-xl p-10 cursor-pointer ` +
    (isDragging ? "bg-gray-50" : "")
  ), [isDragging]);

  return (
    <div
      role="button"
      tabIndex={0}
      aria-label="Upload files by dragging and dropping or by browsing"
      onClick={openFileDialog}
      onKeyDown={(e) => {
        if (e.key === "Enter" || e.key === " ") {
          openFileDialog();
        }
      }}
      onDragOver={handleDragOver}
      onDragEnter={handleDragOver}
      onDragLeave={handleDragLeave}
      onDrop={handleDrop}
      className={(className ? className + " " : "") + baseClass}

    >
      {children ?? (
        <div className="text-center">
          <p className="text-base font-light">Maximum file size: 50mb</p>
          <p className="text-base font-light">Supports .PNG, .JPG, .WEBP, .GIF</p>
        </div>
      )}
      <input
        ref={inputRef}
        type="file"
        accept={accept}
        multiple={allowMultiple}
        className="hidden"
        onChange={(e) => {
          if (e.target.files) {
            let accepted = filterAccepted(e.target.files, accept);
            if (!allowMultiple) accepted = accepted.slice(0, 1);
            if (accepted.length > 0) onFilesAdded(accepted);
          }
          // reset so same files can be selected again
          e.currentTarget.value = "";
        }}
      />
    </div>
  );
};

export default FileDropzone;
