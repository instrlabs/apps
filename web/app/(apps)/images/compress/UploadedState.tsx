"use client";

import React, { useState, useEffect } from "react";
import Icon from "@/components/icon";
import Chip from "@/components/chip";
import Button from "@/components/button";
import ImagePreview from "@/components/image-preview";
import { ExtendedInstructionFile } from "./types";

interface UploadedStateProps {
  files: ExtendedInstructionFile[];
  onRemoveFile: (fileIndex: number) => void;
  onSubmit: () => void;
  onReset: () => void;
  onDownload: (fileName: string) => void;
}

export default function UploadedState({
  files,
  onRemoveFile,
  onSubmit,
  onReset,
  onDownload,
}: UploadedStateProps) {
  const [previews, setPreviews] = useState<Record<number, string>>({});

  useEffect(() => {
    const generatePreviews = async () => {
      const newPreviews: Record<number, string> = {};

      for (let index = 0; index < files.length; index++) {
        const item = files[index];
        if (item.file.type.startsWith("image/") && !previews[index]) {
          const reader = new FileReader();
          reader.onload = (e) => {
            newPreviews[index] = e.target?.result as string;
            setPreviews((prev) => ({ ...prev, [index]: e.target?.result as string }));
          };
          reader.readAsDataURL(item.file);
        }
      }
    };

    generatePreviews().then();
  }, [files, previews]);

  return (
    <div className="flex h-full flex-col gap-4">
      {files.map((item, index) => (
        <div
          key={index}
          className="flex items-center justify-between rounded border border-white/10 bg-white/4 p-3"
        >
          <div className="flex flex-1 items-center gap-4">
            {<ImagePreview src={previews[index]} alt={item.file.name} size={40} />}

            <div className="flex flex-col">
              <div className="flex items-center gap-2">
                <p className="text-sm leading-5 font-normal text-white">{item.file.name}</p>
                {item.status && (
                  <Chip
                    label={item.status}
                    state={
                      item.status === "DONE"
                        ? "success"
                        : item.status === "PROCESSING"
                          ? "processing"
                          : "error"
                    }
                  />
                )}
              </div>
              <p className="text-xs leading-4 font-normal text-white/30">
                {(item.file_size / 1024).toFixed(0)}KB
              </p>
            </div>
          </div>

          {item.status === "DONE" && (
            <div className="flex flex-col items-end justify-center gap-2">
              <p className="text-xs leading-4 font-semibold text-green-500">
                {(item.compressedSize || 0) / 1024 < 1
                  ? Math.round(item.compressedSize || 0) + "B"
                  : ((item.compressedSize || 0) / 1024).toFixed(0) + "KB"}{" "}
                ({-Math.abs(item.compressionPercentage || 0)}%)
              </p>
              <button
                onClick={() => onDownload(item.file.name)}
                className="text-xs leading-4 font-semibold text-blue-400 transition-colors hover:text-blue-300"
              >
                Download
              </button>
            </div>
          )}

          {item.status === "PROCESSING" && (
            <Icon name="progress" size={24} className="text-white/60" />
          )}

          {!item.status && (
            <button
              onClick={() => onRemoveFile(index)}
              className="text-white/60 transition-colors hover:text-white"
              aria-label="Remove file"
            >
              <Icon name="close" size={24} />
            </button>
          )}
        </div>
      ))}

      {/* Action Buttons */}
      <div className="flex gap-2">
        <Button variant="primary" size="base" onClick={onSubmit} className="min-w-[150px]">
          Submit
        </Button>
        <Button variant="secondary" size="base" onClick={onReset} className="min-w-[150px]">
          Reset
        </Button>
      </div>
    </div>
  );
}
