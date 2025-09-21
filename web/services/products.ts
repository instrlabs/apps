"use server"

import { APIs } from "@/constants/api";
import {fetchGET} from "@/utils";

export type ApiResponse<TBody> = {
  success: boolean;
  message: string;
  data: TBody | null;
  errors: FormErrors | null;
};

export type FormErrors = {
  errorMessage: string;
  fieldName: string;
}[] | null;

export type Product = {
  id?: string;
  key: string;
  name: string;
  price: number;
  description?: string;
  image?: string;
  productType: string;
  userId: string;
  active: boolean;
  isFree: boolean;
  createdAt?: string;
  updatedAt?: string;
};


export async function listProducts(): Promise<ApiResponse<Product[]>> {
  return await fetchGET(`${APIs.PAYMENTS}/products`);
}

// export async function getProduct(id: string): Promise<ApiResponse<Product>> {
//   const endpoint = `${APIs.PRODUCTS}/${encodeURIComponent(id)}`;
//   return callProducts<Product>(endpoint, {
//     method: "GET",
//     withAuth: true,
//   });
// }
//
// export async function updateProduct(id: string, patch: Partial<Product>): Promise<ApiResponse<{ status: string }>> {
//   const endpoint = `${APIs.PRODUCTS}/${encodeURIComponent(id)}`;
//   return callProducts<{ status: string }>(endpoint, {
//     method: "PATCH",
//     body: JSON.stringify(patch),
//     withAuth: true,
//   });
// }
//
// export async function deleteProduct(id: string): Promise<ApiResponse<null>> {
//   const endpoint = `${APIs.PRODUCTS}/${encodeURIComponent(id)}`;
//   const res = await callProducts<null>(endpoint, {
//     method: "DELETE",
//     withAuth: true,
//   });
//   // DELETE returns 204; our normalizer sets data to null already
//   return res;
// }
