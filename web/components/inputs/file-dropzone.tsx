"use client";

import React, { useCallback, useMemo, useRef, useState, useId } from "react";
import { bytesToString } from "@/utils/bytesToString";
import { acceptsToExtensions } from "@/utils/acceptsToExtensions";
import useNotification from "@/hooks/useNotification";
import CloudUploadIcon from "@/components/icons/CloudUploadIcon";
import Text from "@/components/text";

export type FileDropzoneProps = {
  accepts: string[];
  multiple: boolean;
  onFilesAdded: (files: File[]) => void;
  maxSize: number;
  className?: string;
};

function validateFiles(files: File[], accepts: string[], maxFileSize: number): boolean {
  const acceptedTypes = accepts.map((s) => s.trim()).filter(Boolean);
  return files.every((file) => acceptedTypes.includes(file.type) && file.size <= maxFileSize);
}

const FileDropzone: React.FC<FileDropzoneProps> = ({ accepts, onFilesAdded, multiple, maxSize, className }) => {
  const { showNotification } = useNotification();
  const [isDragging, setIsDragging] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);
  const helperId = useId();

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
        showNotification({ message: "Only one file can be dropped at a time.", type: "error" });
        return;
      }

      if (!validateFiles(files, accepts, maxSize)) {
        showNotification({ message: "Invalid file type or file size.", type: "error" });
        return;
      }

      onFilesAdded(files);
      e.dataTransfer.clearData();
    }
  }, [multiple, accepts, maxSize, onFilesAdded, showNotification]);

  const baseClass = useMemo(() => (
    `
    group cursor-pointer outline-none
    flex w-full flex-col items-center justify-center
    gap-2 p-6
    rounded-lg border border-dashed border-white/10
    bg-transparent

    transition-colors focus-visible:ring-2 focus-visible:ring-white/20
    ${isDragging ? "bg-white/8" : "hover:bg-white/5"}
    ${className || ""}
    `
  ), [isDragging, className]);

  return (
    <div
      role="button"
      tabIndex={0}
      aria-label={`Upload files. Allowed: ${acceptsToExtensions(accepts).join(", ")}. Max size: ${bytesToString(maxSize)}.`}
      aria-describedby={helperId}
      onKeyDown={(e) => {
        if (e.key === "Enter" || e.key === " ") {
          e.preventDefault();
          openFileDialog();
        }
      }}
      onClick={openFileDialog}
      onDragOver={handleDragOver}
      onDragEnter={handleDragOver}
      onDragLeave={handleDragLeave}
      onDrop={handleDrop}
      className={baseClass}
    >
      <div className={`
          flex flex-col items-center justify-center gap-2
          rounded-lg
          group-hover:bg-white/1
        `}>
        <CloudUploadIcon className="size-10" />
        <Text as="span" xSize="sm" className="font-light text-white/50 group-hover:text-white transition-colors">Upload Files</Text>
      </div>
      <Text
        as="span"
        id={helperId}
        xSize="xs"
        xColor="secondary"
        className="max-w-xs text-center"
      >
        Total file size allowed is {bytesToString(maxSize)}, and the supported formats are {acceptsToExtensions(accepts).join(", ")}.
      </Text>
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
                message: "Only one file can be dropped at a time."
              });
            } else if (!validateFiles(files, accepts, maxSize)) {
              showNotification({
                type: "error",
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
