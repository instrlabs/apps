"use server"

import { API_IMAGES } from "@/constants/api";
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
  return await fetchGET<{ products: Product[] }>(API_IMAGES + "/products");
}

export async function createInstruction(productKey: string) {
  return await fetchPOST<{ instruction: Instruction }>(`${API_IMAGES}/instructions`, {
    productKey
  });
}

export async function createInstructionDetails(instructionId: string, file: File) {
  const formData = new FormData();
  formData.append("file", file);
  return await fetchPOSTFormData<{ input: InstructionFile, output: InstructionFile }>(`${API_IMAGES}/instructions/${instructionId}/details`, formData);
}

export async function getImageInstructions() {
  return await fetchGET<{ instructions: Instruction[] }>(API_IMAGES + "/instructions");
}

export async function getImageInstruction(id: string) {
  return await fetchGET<{ instruction: Instruction }>(`${API_IMAGES}/instructions/${id}`);
}

export async function getInstructionDetails(id: string) {
  return await fetchGET<{ files: InstructionFile[] }>(`${API_IMAGES}/instructions/${id}/details`);
}

export async function getInstructionFileBytes(id: string, fileId: string) {
  return await fetchGETBytes(`${API_IMAGES}/instructions/${id}/details/${fileId}`);
}



