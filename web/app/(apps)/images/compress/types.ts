import { InstructionFile } from "@/services/images";

export type CompressState = "default" | "uploaded";

export interface ExtendedInstructionFile extends InstructionFile {
  file: File;
  outputFile?: InstructionFile;
  compressedSize?: number;
  compressionPercentage?: number;
}
