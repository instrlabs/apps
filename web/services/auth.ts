"use server"

import {cookies} from "next/headers";
import { APIs } from "@/constants/api";
import { ResponseCookie, ResponseCookies } from "next/dist/compiled/@edge-runtime/cookies";
import { ApiResponse, EmptyBody, FormErrors, fetchGET, fetchPOST, fetchPUT } from "@/utils/fetch";
import {redirect} from "next/navigation";

export interface ProfileResponse {
  name: string;
  email: string
}
export async function loginUser({ email, pin }: {
  email: string, pin: string
}): Promise<ApiResponse<null>> {
  const res = await fetch(process.env.API_URL + `${APIs.AUTH}/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, pin }),
  });

  const isOK = res.ok;
  const resBody = await res.json();
  if (!isOK) {
    return {
      success: isOK,
      message: resBody.message,
      data: null,
      errors: (resBody.errors as FormErrors) ?? null,
    }
  }

  const reqSetCookie = new ResponseCookies(res.headers);
  const storeCookie = await cookies();
  storeCookie.set(reqSetCookie.get("AccessToken") as ResponseCookie);
  storeCookie.set(reqSetCookie.get("RefreshToken") as ResponseCookie);

  return {
    success: isOK,
    message: resBody.message,
    data: isOK ? resBody.data : null,
    errors: !isOK ? (resBody.errors as FormErrors) : null,
  };
}

export async function refreshToken(): Promise<ApiResponse<EmptyBody>> {
  const res = await fetch(process.env.API_URL + `${APIs.AUTH}/refresh`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ refresh_token: (await cookies()).get("RefreshToken")?.value }),
  });

  const isOK = res.ok;
  const resBody = await res.json();
  const reqSetCookie = new ResponseCookies(res.headers);
  const storeCookie = await cookies();
  storeCookie.set(reqSetCookie.get("AccessToken") as ResponseCookie);
  storeCookie.set(reqSetCookie.get("RefreshToken") as ResponseCookie);

  return {
    success: isOK,
    message: resBody.message,
    data: isOK ? resBody.data : null,
    errors: !isOK ? resBody.errors : null,
  };
}

export async function logoutUser(): Promise<ApiResponse<EmptyBody>> {
  await fetchPOST<EmptyBody>(`${APIs.AUTH}/logout`, {});
  const storeCookie = await cookies();
  storeCookie.delete("AccessToken");
  storeCookie.delete("RefreshToken");
  redirect("/login")
}

export async function sendPIN({ email }: {
  email: string,
}): Promise<ApiResponse<EmptyBody>> {
  return await fetchPOST(`${APIs.AUTH}/send-pin`, { email });
}

export async function getProfile(): Promise<ApiResponse<ProfileResponse>> {
  return await fetchGET(`${APIs.AUTH}/profile`);
}

export async function updateProfile(name: string): Promise<ApiResponse<EmptyBody>> {
  return await fetchPUT(`${APIs.AUTH}/profile`, { name });
}

export async function changePassword(
  current_password: string,
  new_password: string
): Promise<ApiResponse<EmptyBody>> {
  return await fetchPOST(`${APIs.AUTH}/change-password`, { current_password, new_password });
}
