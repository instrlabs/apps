"use server"

import { IMAGES } from "@/constants/APIs";
import { fetchGET, fetchGETBytes, fetchPOST, fetchPOSTFormData } from "@/utils/fetch";

export type Product = {
  id: string;
  key: string;
  title: string;
  description: string;
  product_type: string;
  is_active: string;
  is_free: string;
}

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

export async function getProducts() {
  return await fetchGET<{ products: Product[] }>(IMAGES + "/products");
}

export async function createImageInstruction(productKey: string) {
  return await fetchPOST<{ instruction: Instruction }>(`${IMAGES}/instructions`, {
    productKey
  });
}

export async function createImageInstructionDetails(instructionId: string, file: File) {
  const formData = new FormData();
  formData.append("file", file);
  return await fetchPOSTFormData<{ input: InstructionFile, output: InstructionFile }>(`${IMAGES}/instructions/${instructionId}/details`, formData);
}

export async function getImageInstructions() {
  return await fetchGET<{ instructions: Instruction[] }>(IMAGES + "/instructions");
}

export async function getImageInstructionDetails(id: string) {
  return await fetchGET<{ files: InstructionFile[] }>(`${IMAGES}/instructions/${id}/details`);
}

export async function getImageFile(id: string, fileId: string) {
  return await fetchGETBytes(`${IMAGES}/instructions/${id}/details/${fileId}`);
}



