import { AUTH_ENDPOINTS } from "@/constants/api";
import { fetchWithErrorHandling } from "@/utils";

interface LoginResponse {
  message: string;
  data: {
    access_token: string;
    refresh_token: string;
  }
}

interface RegisterResponse {
  message: string;
  data: {
    email: string;
  }
}

interface GoogleCallbackResponse {
  message: string;
  data: {
    access_token: string;
    refresh_token: string;
  }
}

interface RefreshTokenResponse {
  message: string;
  data: {
    access_token: string;
    refresh_token: string;
  }
}

interface ForgotPasswordResponse {
  message: string;
}

interface ResetPasswordResponse {
  message: string;
}

interface VerifyTokenResponse {
  message: string;
  data: {
    user: { [key: string]: unknown };
  }
}

interface WrapperResponse<T> {
  data: T | null;
  error: string | null;
}

export async function loginUser(email: string, password: string): Promise<WrapperResponse<LoginResponse>> {
  return await fetchWithErrorHandling(AUTH_ENDPOINTS.LOGIN, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
  });
}

export async function handleGoogleCallback(code: string): Promise<WrapperResponse<GoogleCallbackResponse>> {
  return fetchWithErrorHandling(AUTH_ENDPOINTS.GOOGLE_CALLBACK, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ code }),
  });
}

export async function registerUser(email: string, password: string): Promise<WrapperResponse<RegisterResponse>> {
  return fetchWithErrorHandling(AUTH_ENDPOINTS.REGISTER, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
  });
}

export async function requestPasswordReset(email: string): Promise<WrapperResponse<ForgotPasswordResponse>> {
  return fetchWithErrorHandling(AUTH_ENDPOINTS.FORGOT_PASSWORD, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email }),
  });
}

export async function resetPassword(token: string, password: string): Promise<WrapperResponse<ResetPasswordResponse>> {
  return fetchWithErrorHandling(AUTH_ENDPOINTS.RESET_PASSWORD, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ token, new_password: password }),
  });
}

export async function refreshToken(refreshToken: string): Promise<WrapperResponse<RefreshTokenResponse>> {
  return fetchWithErrorHandling(AUTH_ENDPOINTS.REFRESH_TOKEN, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ refresh_token: refreshToken }),
  });
}

export async function verifyToken(token: string): Promise<WrapperResponse<VerifyTokenResponse>> {
  return fetchWithErrorHandling(AUTH_ENDPOINTS.VERIFY_TOKEN, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ token }),
  });
}