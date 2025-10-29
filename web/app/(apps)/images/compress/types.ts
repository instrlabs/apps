import { InstructionDetail } from "@/services/images";

export type CompressState = "default" | "uploaded";

export interface ExtendedInstructionFile extends InstructionDetail {
  file: File;
  outputDetail?: InstructionDetail;
  compressedSize?: number;
  compressionPercentage?: number;
}
