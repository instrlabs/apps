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
  onDownload: (instructionId: string, detailId: string, fileName: string) => void;
}

interface ImageDimensions {
  width: number;
  height: number;
}

export default function UploadedState({ files, onRemoveFile, onSubmit, onReset, onDownload }: UploadedStateProps) {
  const [previews, setPreviews] = useState<Record<number, string>>({});
  const [dimensions, setDimensions] = useState<Record<number, ImageDimensions>>({});

  useEffect(() => {
    const generatePreviews = async () => {
      const newPreviews: Record<number, string> = {};
      const newDimensions: Record<number, ImageDimensions> = {};

      for (let index = 0; index < files.length; index++) {
        const item = files[index];
        if (item.file.type.startsWith("image/") && !previews[index]) {
          const reader = new FileReader();
          reader.onload = (e) => {
            const dataUrl = e.target?.result as string;
            newPreviews[index] = dataUrl;
            setPreviews((prev) => ({ ...prev, [index]: dataUrl }));

            // Get image dimensions
            const img = new window.Image();
            img.onload = () => {
              newDimensions[index] = { width: img.naturalWidth, height: img.naturalHeight };
              setDimensions((prev) => ({ ...prev, [index]: { width: img.naturalWidth, height: img.naturalHeight } }));
            };
            img.src = dataUrl;
          };
          reader.readAsDataURL(item.file);
        }
      }
    };

    generatePreviews().then();
  }, [files, previews]);

  return (
    <div className="flex h-full flex-col gap-2 sm:gap-3 md:gap-4">
      <div className="flex flex-col gap-2">
        {files.map((item, index) => (
          <div key={index} className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-2 sm:gap-3 rounded border border-white/10 bg-white/4 p-2 sm:p-3">
            <div className="flex flex-1 items-center gap-2 sm:gap-3 md:gap-4 w-full sm:w-auto">
              {previews[index] ? (
                <ImagePreview src={previews[index]} alt={item.file.name} size={40} />
              ) : (
                <Icon name="rectangle" size={40} className="text-white/60" />
              )}

              <div className="flex flex-col gap-1">
                <div className="flex items-center gap-2">
                  <p className="truncate text-sm leading-5 font-normal text-white">{item.file.name}</p>
                  {item.status && (
                    <Chip
                      label={item.status}
                      state={
                        item.status === "DONE" ? "success" : item.status === "PROCESSING" ? "processing" : "default"
                      }
                    />
                  )}
                </div>
                <p className="text-xs leading-4 font-normal text-white/30">
                  {(item.file_size / 1024).toFixed(0)}KB · {item.mime_type?.split("/")[1]?.toUpperCase() || "FILE"} ·
                  {dimensions[index] ? `${dimensions[index].width}x${dimensions[index].height}` : "loading..."}
                </p>
              </div>
            </div>

            {/* Right: Actions/Status */}
            <div className="flex items-center justify-end gap-2 sm:gap-3 md:gap-4 w-full sm:w-auto">
              {item.status === "DONE" && (
                <div className="flex flex-col items-end justify-center gap-1">
                  <p className="text-xs leading-4 font-semibold" style={{ color: "#34a853" }}>
                    {(item.compressedSize || 0) / 1024 < 1
                      ? Math.round(item.compressedSize || 0) + "B"
                      : ((item.compressedSize || 0) / 1024).toFixed(0) + "KB"}
                    {item.compressionPercentage !== undefined && ` (${-Math.abs(item.compressionPercentage)}%)`}
                  </p>
                  <button
                    onClick={() => {
                      if (item.instruction_id && item.output_id) {
                        onDownload(item.instruction_id, item.output_id, item.file.name);
                      }
                    }}
                    className="text-xs leading-4 font-semibold transition-colors"
                    style={{ color: "#4285f4" }}
                  >
                    Download
                  </button>
                </div>
              )}

              {item.status === "PROCESSING" && (
                <Icon name="progress" size={24} className="animate-spin text-white/60" />
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
          </div>
        ))}
      </div>

      {/* Action Buttons */}
      <div className="flex flex-col sm:flex-row gap-2">
        {!files[0].status && (
          <Button variant="primary" size="base" onClick={onSubmit} className="min-w-full sm:min-w-[150px]">
            Submit
          </Button>
        )}
        <Button variant="secondary" size="base" onClick={onReset} className="min-w-full sm:min-w-[150px]">
          Reset
        </Button>
      </div>
    </div>
  );
}
