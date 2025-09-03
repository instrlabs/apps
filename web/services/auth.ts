"use server"

import {cookies} from "next/headers";
import { APIs } from "@/constants/api";
import {ResponseCookie, ResponseCookies} from "next/dist/compiled/@edge-runtime/cookies";



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

interface RegisterBody {
  email: string
}

interface ProfileBody {
  user: {
    name: string;
    email: string
  }
}

type EmptyBody = { message?: string } & Record<string, unknown>;

type RegisterResponse = ApiResponse<RegisterBody>;
type LoginResponse = ApiResponse<EmptyBody>;
type RefreshTokenResponse = ApiResponse<EmptyBody>;
type ForgotPasswordResponse = ApiResponse<EmptyBody>;
type ResetPasswordResponse = ApiResponse<EmptyBody>;
type ProfileResponse = ApiResponse<ProfileBody>;
type UpdateProfileResponse = ApiResponse<ProfileBody>;


export async function registerUser({ name, email, password }: {
  name: string,
  email: string,
  password: string
}): Promise<RegisterResponse> {
  const baseUrl = process.env.API_URL;
  const url = `${baseUrl}${APIs.AUTH_REGISTER}`;
  const res = await fetch(url, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ name, email, password })
  });

  const isOK = res.ok;
  const resBody = await res.json();

  return {
    success: isOK,
    message: resBody.message,
    data: isOK ? resBody.data : null,
    errors: !isOK ? resBody.errors : null,
  };
}

export async function loginUser({ email, password }: {
  email: string,
  password: string
}): Promise<LoginResponse> {
  const baseUrl = process.env.API_URL;
  const url = `${baseUrl}${APIs.AUTH_LOGIN}`;
  const res = await fetch(url, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
  });

  const isOK = res.ok;
  const resBody = await res.json();
  const reqSetCookie = new ResponseCookies(res.headers)
  const accessToken = reqSetCookie.get("AccessToken");
  const refreshToken = reqSetCookie.get("RefreshToken");
  const storeCookie = await cookies();
  if (accessToken) storeCookie.set(accessToken);
  if (refreshToken) storeCookie.set(refreshToken);

  return {
    success: isOK,
    message: resBody.message,
    data: isOK ? resBody.data : null,
    errors: !isOK ? (resBody.errors as FormErrors) : null,
  };
}

export async function refreshToken(): Promise<RefreshTokenResponse> {
  const baseUrl = process.env.API_URL;
  const url = `${baseUrl}${APIs.AUTH_REFRESH}`;
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
  reqSetCookie.getAll().forEach(
    (cookie: ResponseCookie) => storeCookie.set(cookie))

  return {
    success: isOK,
    message: resBody.message,
    data: isOK ? resBody.data : null,
    errors: !isOK ? resBody.errors : null,
  };
}

export async function requestPasswordReset(email: string): Promise<ForgotPasswordResponse> {
  const baseUrl = process.env.API_URL;
  const url = `${baseUrl}${APIs.AUTH_FORGOT_PASSWORD}`;
  const res = await fetch(url, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email }),
  });

  const isOK = res.ok;
  const resBody = await res.json();
  return {
    success: isOK,
    message: resBody.message,
    data: isOK ? resBody.data : null,
    errors: !isOK ? resBody.errors : null,
  };
}

export async function resetPassword(token: string, new_password: string): Promise<ResetPasswordResponse> {
  const baseUrl = process.env.API_URL;
  const url = `${baseUrl}${APIs.AUTH_RESET_PASSWORD}`;
  const res = await fetch(url, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ token, new_password }),
  });

  const isOK = res.ok;
  const resBody = await res.json();
  return {
    success: isOK,
    message: resBody.message,
    data: isOK ? resBody.data : null,
    errors: !isOK ? resBody.errors : null,
  };
}

export async function profile(): Promise<ProfileResponse> {
  const baseUrl = process.env.API_URL;
  const url = `${baseUrl}${APIs.AUTH_PROFILE}`;
  const storeCookie = await cookies();
  const res = await fetch(url, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      "Cookie": storeCookie.toString()
    },
  });

  const isOK = res.ok;
  const resBody = await res.json();
  return {
    success: isOK,
    message: resBody.message,
    data: isOK ? resBody.data : null,
    errors: !isOK ? resBody.errors : null,
  };
}

export async function logoutUser(): Promise<LoginResponse> {
  const baseUrl = process.env.API_URL;
  const url = `${baseUrl}${APIs.AUTH_LOGOUT}`;
  const storeCookie = await cookies();
  const res = await fetch(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "Cookie": storeCookie.toString()
    },
  });

  storeCookie.delete("AccessToken");
  storeCookie.delete("RefreshToken");

  const isOK = res.ok;
  const resBody = await res.json();
  return {
    success: isOK,
    message: resBody.message,
    data: isOK ? resBody.data : null,
    errors: !isOK ? resBody.errors : null,
  };
}

export async function updateProfile(name: string): Promise<UpdateProfileResponse> {
  const baseUrl = process.env.API_URL;
  const url = `${baseUrl}${APIs.AUTH_PROFILE}`;
  const storeCookie = await cookies();
  const res = await fetch(url, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
      "Cookie": storeCookie.toString()
    },
    body: JSON.stringify({ name }),
  });

  const isOK = res.ok;
  const resBody = await res.json();
  return {
    success: isOK,
    message: resBody.message,
    data: isOK ? resBody.data : null,
    errors: !isOK ? resBody.errors : null,
  };
}

export async function changePassword(current_password: string, new_password: string): Promise<LoginResponse> {
  const baseUrl = process.env.API_URL;
  const url = `${baseUrl}${APIs.AUTH_CHANGE_PASSWORD}`;
  const storeCookie = await cookies();
  const res = await fetch(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "Cookie": storeCookie.toString()
    },
    body: JSON.stringify({ current_password, new_password }),
  });

  const isOK = res.ok;
  const resBody = await res.json();
  return {
    success: isOK,
    message: resBody.message,
    data: isOK ? resBody.data : null,
    errors: !isOK ? resBody.errors : null,
  };
}
