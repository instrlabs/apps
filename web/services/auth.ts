"use server"

import { headers } from "next/headers";
import { APIs } from "@/constants/api";



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


export async function registerUser(name: string, email: string, password: string): Promise<RegisterResponse> {
  const protocol = process.env.NODE_ENV === "development" ? "http" : "https";
  const h = await headers();
  const url = `${protocol}://${h.get("host")}${APIs.AUTH_REGISTER}`;
  const res = await fetch(url, {
    method: "POST",
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

export async function loginUser(email: string, password: string): Promise<LoginResponse> {
  const protocol = process.env.NODE_ENV === "development" ? "http" : "https";
  const h = await headers();
  const url = `${protocol}://${h.get("host")}${APIs.AUTH_LOGIN}`;
  const res = await fetch(url, {
    method: "POST",
    body: JSON.stringify({ email, password }),
  });

  const isOK = res.ok;
  const resBody = await res.json();

  return {
    success: isOK,
    message: resBody.message,
    data: isOK ? resBody.data : null,
    errors: !isOK ? (resBody.errors as FormErrors) : null,
  };
}

export async function refreshToken(): Promise<RefreshTokenResponse> {
  const protocol = process.env.NODE_ENV === "development" ? "http" : "https";
  const h = await headers();
  const url = `${protocol}://${h.get("host")}${APIs.AUTH_REFRESH}`;
  const res = await fetch(url, {
    method: "POST",
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

export async function requestPasswordReset(email: string): Promise<ForgotPasswordResponse> {
  const protocol = process.env.NODE_ENV === "development" ? "http" : "https";
  const h = await headers();
  const url = `${protocol}://${h.get("host")}${APIs.AUTH_FORGOT_PASSWORD}`;
  const res = await fetch(url, {
    method: "POST",
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
  const protocol = process.env.NODE_ENV === "development" ? "http" : "https";
  const h = await headers();
  const url = `${protocol}://${h.get("host")}${APIs.AUTH_RESET_PASSWORD}`;
  const res = await fetch(url, {
    method: "POST",
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
  const protocol = process.env.NODE_ENV === "development" ? "http" : "https";
  const h = await headers();
  const url = `${protocol}://${h.get("host")}${APIs.AUTH_PROFILE}`;
  const res = await fetch(url, {
    method: "GET",
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
  const protocol = process.env.NODE_ENV === "development" ? "http" : "https";
  const h = await headers();
  const url = `${protocol}://${h.get("host")}${APIs.AUTH_LOGOUT}`;
  const res = await fetch(url, {
    method: "POST",
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

export async function updateProfile(name: string): Promise<UpdateProfileResponse> {
  const protocol = process.env.NODE_ENV === "development" ? "http" : "https";
  const h = await headers();
  const url = `${protocol}://${h.get("host")}${APIs.AUTH_PROFILE}`;
  const res = await fetch(url, {
    method: "PUT",
    body: JSON.stringify({ name }),
  });

  const isOK = res.ok;
  const resBody = await res.json();
  return {
    success: isOK,
    message: resBody.message,
    data: isOK ? resBody : null,
    errors: !isOK ? resBody.errors : null,
  };
}

export async function changePassword(current_password: string, new_password: string): Promise<LoginResponse> {
  const protocol = process.env.NODE_ENV === "development" ? "http" : "https";
  const h = await headers();
  const url = `${protocol}://${h.get("host")}${APIs.AUTH_CHANGE_PASSWORD}`;
  const res = await fetch(url, {
    method: "POST",
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
