"use server";

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
};

export type Instruction = {
  id: string;
  user_id: string;
  product_id: string;
  created_at: string;
  updated_at: string;
};

export type InstructionDetail = {
  id?: string;
  instruction_id?: string;
  file_name: string;
  file_size: number;
  mime_type: string;
  status?: "FAILED" | "PENDING" | "PROCESSING" | "DONE";
  output_id?: string;
};

export async function getProducts() {
  return await fetchGET<{ products: Product[] }>(IMAGES + "/products");
}

export async function createInstruction(productId: string) {
  return await fetchPOST<{ instruction: Instruction }>(`${IMAGES}/instructions`, {
    product_id: productId,
  });
}

export async function createInstructionDetail(instructionId: string, file: File) {
  const formData = new FormData();
  formData.append("file", file);
  return await fetchPOSTFormData<{ input: InstructionDetail; output: InstructionDetail }>(
    `${IMAGES}/instructions/${instructionId}/details`,
    formData,
  );
}

export async function getInstructionDetails(id: string) {
  return await fetchGET<{ files: InstructionDetail[] }>(`${IMAGES}/instructions/${id}/details`);
}

export async function getInstructions() {
  return await fetchGET<{ instructions: Instruction[] }>(IMAGES + "/instructions");
}

export async function getFile(id: string, fileId: string) {
  return await fetchGETBytes(`${IMAGES}/instructions/${id}/details/${fileId}`);
}
