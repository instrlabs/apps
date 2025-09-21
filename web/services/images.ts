"use server"

import { APIs } from "@/constants/api";
import { fetchGET, fetchGETBytes, fetchPOST, fetchPOSTFormData } from "@/utils/fetch";

export type Instruction = {
  id: string;
  user_id: string;
  product_id: string;
  created_at: string;
  updated_at: string;
};

export type InstructionFile = {
  id: string;
  instruction_id: string;
  original_name: string;
  file_name: string;
  size: number;
  status: string;
  output_id?: string;
};

export async function createInstruction(productKey: string) {
  return await fetchPOST<{ instruction: Instruction }>(`${APIs.IMAGES}/instructions/${productKey}`);
}

export async function createInstructionDetails(instructionId: string, file: File) {
  const formData = new FormData();
  formData.append("file", file);
  return await fetchPOSTFormData<{ file: InstructionFile }>(`${APIs.IMAGES}/instructions/${instructionId}/details`, formData);
}

export async function getImageInstructions() {
  return await fetchGET<{ instructions: Instruction[] }>(APIs.IMAGES + "/instructions");
}

export async function getImageInstruction(id: string) {
  return await fetchGET<{ instruction: Instruction }>(`${APIs.IMAGES}/instructions/${id}`);
}

export async function getInstructionDetails(id: string) {
  return await fetchGET<{ files: InstructionFile[] }>(`${APIs.IMAGES}/instructions/${id}/details`);
}

export async function getInstructionFileBytes(id: string, fileId: string) {
  return await fetchGETBytes(`${APIs.IMAGES}/instructions/${id}/details/${fileId}`);
}



