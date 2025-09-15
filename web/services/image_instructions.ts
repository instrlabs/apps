"use server"

import { APIs } from "@/constants/api";
import {ApiResponse, fetchGET, fetchGETBytes, fetchPOSTFormData} from "@/utils/fetch";

export type ImageInstruction = {
  id: string;
  user_id: string;
  product_id: string;
  inputs: {
    file_name: string;
    size: number;
  }[];
  outputs: {
    file_name: string;
    size: number;
  }[];
  status: string;
  created_at: string;
  updated_at: string;
};

export async function getImageInstructions(): Promise<ApiResponse<ImageInstruction[]>> {
  return await fetchGET<ImageInstruction[]>(APIs.IMAGE_INSTRUCTIONS);
}

export async function getImageInstructionFile(id: string, fileName: string): Promise<ApiResponse<ArrayBuffer>> {
  return await fetchGETBytes(`${APIs.IMAGE_INSTRUCTIONS}/${id}/${fileName}`);
}

export async function compressImage(files: File[]): Promise<ApiResponse<ImageInstruction[]>> {
  const formData = new FormData();
  files.forEach(file => formData.append("files", file));
  return await fetchPOSTFormData<ImageInstruction[]>(APIs.IMAGE_COMPRESS, formData);
}


