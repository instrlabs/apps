"use client";

import React from "react";
import FileDropzone from "@/components/file-dropzone";

interface DefaultStateProps {
  onFilesAdded: (files: File[]) => void;
}

export default function DefaultState({ onFilesAdded }: DefaultStateProps) {
  return (
    <FileDropzone
      accepts={["image/jpeg", "image/png", "image/webp"]}
      multiple
      maxSize={5242880} // 5MB
      title="Upload Files"
      onFilesAdded={onFilesAdded}
      className="h-full"
    />
  );
}
