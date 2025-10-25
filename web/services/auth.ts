"use server"

import {cookies} from "next/headers";
import { AUTH } from "@/constants/APIs";
import { ApiResponse, EmptyBody, fetchGET, fetchPOST } from "@/utils/fetch";
import {redirect} from "next/navigation";
import { RedirectType } from "next/dist/client/components/redirect-error";

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

  return await fetchPOST(`${AUTH}/login`, { email, pin });
}

export async function refresh() {
  if (typeof window !== "undefined") {
    throw new Error("Function: refresh() must be call on client component");
  }

  return await fetchPOST(`${AUTH}/refresh`);
}

export async function logout() {
  await fetchPOST<EmptyBody>(`${AUTH}/logout`, {});
  const storeCookie = await cookies();
  storeCookie.delete("access_token");
  storeCookie.delete("refresh_token");
}

export async function sendPin({ email }: {
  email: string,
}): Promise<ApiResponse<EmptyBody>> {
  return await fetchPOST(`${AUTH}/send-pin`, { email });
}

export async function getProfile() {
  return await fetchGET<{ user: User }>(`${AUTH}/profile`);
}

export async function loginByGoogle(): Promise<ApiResponse<User>> {
  return redirect(`${process.env.GATEWAY_URL}${AUTH}/google`);
}

