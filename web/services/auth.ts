"use server"

import {cookies} from "next/headers";
import { APIs } from "@/constants/api";
import { ResponseCookie, ResponseCookies } from "next/dist/compiled/@edge-runtime/cookies";
import { ApiResponse, EmptyBody, FormErrors, fetchGET, fetchPOST, fetchPUT } from "@/utils/fetch";

interface RegisterResponse {
  email: string
}

interface ProfileResponse {
  user: {
    name: string;
    email: string
  }
}

export async function registerUser({ name, email, password }: {
  name: string,
  email: string,
  password: string
}): Promise<ApiResponse<RegisterResponse>> {
  return await fetchPOST(APIs.AUTH_REGISTER, { name, email, password });
}

export async function loginUser({ email, password }: {
  email: string,
  password: string
}): Promise<ApiResponse<EmptyBody>> {
  const url = process.env.API_URL + APIs.AUTH_LOGIN;
  const res = await fetch(url, {
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
  reqSetCookie.getAll().forEach((cookie: ResponseCookie) => storeCookie.set(cookie));

  return {
    success: isOK,
    message: resBody.message,
    data: isOK ? resBody.data : null,
    errors: !isOK ? (resBody.errors as FormErrors) : null,
  };
}

export async function refreshToken(): Promise<ApiResponse<EmptyBody>> {
  const url = process.env.API_URL + APIs.AUTH_REFRESH;
  const storeCookie = await cookies();
  const res = await fetch(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "Cookie": storeCookie.toString()
    },
  });

  const isOK = res.ok;
  const resBody = await res.json();
  const reqSetCookie = new ResponseCookies(res.headers);
  reqSetCookie.getAll().forEach((cookie: ResponseCookie) => storeCookie.set(cookie));

  return {
    success: isOK,
    message: resBody.message,
    data: isOK ? resBody.data : null,
    errors: !isOK ? resBody.errors : null,
  };
}

export async function requestPasswordReset(email: string): Promise<ApiResponse<EmptyBody>> {
  return await fetchPOST(APIs.AUTH_FORGOT_PASSWORD, { email });
}

export async function resetPassword(token: string, new_password: string): Promise<ApiResponse<EmptyBody>> {
  return await fetchPOST(APIs.AUTH_RESET_PASSWORD, { token, new_password });
}

export async function profile(): Promise<ApiResponse<ProfileResponse>> {
  return await fetchGET(APIs.AUTH_PROFILE);
}

export async function logoutUser(): Promise<ApiResponse<EmptyBody>> {
  const res = await fetchPOST<EmptyBody>(APIs.AUTH_LOGOUT, {});
  const storeCookie = await cookies();
  storeCookie.delete("AccessToken");
  storeCookie.delete("RefreshToken");
  return res;
}

export async function updateProfile(name: string): Promise<ApiResponse<EmptyBody>> {
  return await fetchPUT(APIs.AUTH_PROFILE, { name });
}

export async function changePassword(
  current_password: string,
  new_password: string
): Promise<ApiResponse<EmptyBody>> {
  return await fetchPOST(APIs.AUTH_CHANGE_PASSWORD, { current_password, new_password });
}
