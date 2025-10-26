"use client";

import React from "react";
import Icon from "@/components/icon";
import Chip from "@/components/chip";
import Button from "@/components/button";
import { FileMetadata } from "./types";

interface UploadedStateProps {
  fileMetadata: FileMetadata[];
  onRemoveFile: (id: string) => void;
  onSubmit: () => void;
  onReset: () => void;
  onDownload: (fileName: string) => void;
  formatFileMetadata: (file: File, dimensions?: { width: number; height: number }) => string;
}

export default function UploadedState({
  fileMetadata,
  onRemoveFile,
  onSubmit,
  onReset,
  onDownload,
  formatFileMetadata,
}: UploadedStateProps) {
  return (
    <div className="flex h-full flex-col gap-4">
      {/* File List */}
      {fileMetadata.map((item) => (
        <div
          key={item.id}
          className="flex items-center justify-between rounded border border-white/10 bg-white/4 p-3"
        >
          {/* Left Side - File Info */}
          <div className="flex flex-1 items-center gap-4">
            <Icon name="rectangle" size={40} className="text-white/80" />
            <div className="flex flex-col">
              {/* Filename and Status Chip */}
              <div className="flex items-center gap-2">
                <p className="text-sm leading-5 font-normal text-white">{item.file.name}</p>
                <Chip
                  label={item.status.toUpperCase()}
                  state={
                    item.status === "success"
                      ? "success"
                      : item.status === "processing"
                        ? "processing"
                        : "default"
                  }
                />
              </div>
              {/* Metadata */}
              <p className="text-xs leading-4 font-normal text-white/30">
                {formatFileMetadata(item.file, item.dimensions)}
              </p>
            </div>
          </div>

          {/* Right Side - Status-specific Content */}
          {item.status === "success" && (
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

          {item.status === "processing" && (
            <Icon name="progress" size={24} className="text-white/60" />
          )}

          {item.status === "waiting" && (
            <button
              onClick={() => onRemoveFile(item.id)}
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
