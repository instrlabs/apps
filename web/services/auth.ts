"use server"

import {cookies} from "next/headers";
import { ResponseCookie, ResponseCookies } from "next/dist/compiled/@edge-runtime/cookies";
import { API_AUTH } from "@/constants/api";
import { ApiResponse, EmptyBody, FormErrors, fetchGET, fetchPOST } from "@/utils/fetch";
import {redirect} from "next/navigation";

export interface User {
  username: string
  email: string
}

export async function login({ email, pin }: {
  email: string,
  pin: string
}) {
  const res = await fetch(process.env.API_URL + `${API_AUTH}/login`, {
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

export async function refresh() {
  const refresh_token = (await cookies()).get("RefreshToken")?.value
  const res = await fetch(process.env.API_URL + `${API_AUTH}/refresh`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ refresh_token })
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
  return redirect(`${process.env.API_URL}${API_AUTH}/google`);
}

