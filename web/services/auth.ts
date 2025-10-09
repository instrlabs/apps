"use server"

import {cookies} from "next/headers";
import { API_AUTH } from "@/constants/api";
import { ApiResponse, EmptyBody, fetchGET, fetchPOST } from "@/utils/fetch";
import {redirect} from "next/navigation";

export interface User {
  username: string
  email: string
}


export async function login({ email, pin }: {
  email: string,
  pin: string
}) {
  if (typeof window !== "undefined") {
    throw new Error("Function: login() must be call on client component");
  }

  return await fetchPOST(`${API_AUTH}/login`, { email, pin });
}

export async function refresh() {
  if (typeof window !== "undefined") {
    throw new Error("Function: refresh() must be call on client component");
  }

  return await fetchPOST(`${API_AUTH}/refresh`);
}

export async function logout() {
  await fetchPOST<EmptyBody>(`${API_AUTH}/logout`, {});
  const storeCookie = await cookies();
  storeCookie.delete("AccessToken");
  storeCookie.delete("RefreshToken");
  redirect("/login")
}

export async function sendPin({ email }: {
  email: string,
}): Promise<ApiResponse<EmptyBody>> {
  return await fetchPOST(`${API_AUTH}/send-pin`, { email });
}

export async function getProfile() {
  return await fetchGET<{ user: User }>(`${API_AUTH}/profile`);
}

export async function loginByGoogle(): Promise<ApiResponse<User>> {
  return redirect(`${process.env.GATEWAY_URL}${API_AUTH}/google`);
}

