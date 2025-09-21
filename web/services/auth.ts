"use server"

import {cookies} from "next/headers";
import { APIs } from "@/constants/api";
import { ResponseCookie, ResponseCookies } from "next/dist/compiled/@edge-runtime/cookies";
import { ApiResponse, EmptyBody, FormErrors, fetchGET, fetchPOST, fetchPUT } from "@/utils/fetch";
import {redirect} from "next/navigation";

interface RegisterResponse {
  email: string
}

export interface ProfileResponse {
  name: string;
  email: string
}

export async function registerUser({ name, email, password }: {
  name: string,
  email: string,
  password: string
}): Promise<ApiResponse<RegisterResponse>> {
  return await fetchPOST(`${APIs.AUTH}/register` , { name, email, password });
}

export async function loginUser({ email, password }: {
  email: string,
  password: string
}): Promise<ApiResponse<null>> {
  const res = await fetch(process.env.API_URL + `${APIs.AUTH}/login`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "Cookie": (await cookies()).toString()
    },
    body: JSON.stringify({ email, password }),
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

export async function requestPasswordReset(email: string): Promise<ApiResponse<EmptyBody>> {
  return await fetchPOST(`${APIs.AUTH}/forgot-password`, { email });
}

export async function resetPassword(token: string, new_password: string): Promise<ApiResponse<EmptyBody>> {
  return await fetchPOST(`${APIs.AUTH}/reset-password`, { token, new_password });
}

export async function logoutUser(): Promise<ApiResponse<EmptyBody>> {
  await fetchPOST<EmptyBody>(`${APIs.AUTH}/logout`, {});
  const storeCookie = await cookies();
  storeCookie.delete("AccessToken");
  storeCookie.delete("RefreshToken");
  redirect("/login")
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
