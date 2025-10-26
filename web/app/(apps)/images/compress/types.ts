export type CompressState = "default" | "uploaded";
export type FileStatus = "waiting" | "processing" | "success" | "error";

export interface FileMetadata {
  file: File;
  id: string;
  typeLabel: string;
  dimensions?: { width: number; height: number };
  status: FileStatus;
  compressedSize?: number;
  compressionPercentage?: number;
  downloadUrl?: string;
}
