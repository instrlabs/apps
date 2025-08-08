import { AUTH_ENDPOINTS } from "@/constants/api";
import { fetchWithErrorHandling } from "@/utils";

interface RegisterResponse {
  message: string;
  data: {
    email: string;
  }
}

interface LoginResponse {
  message: string;
}

interface RefreshTokenResponse {
  message: string;
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

export async function registerUser(email: string, password: string): Promise<WrapperResponse<RegisterResponse>> {
  return fetchWithErrorHandling(AUTH_ENDPOINTS.REGISTER, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
  });
}

export async function loginUser(email: string, password: string): Promise<WrapperResponse<LoginResponse>> {
  return await fetchWithErrorHandling(AUTH_ENDPOINTS.LOGIN, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
  });
}

export async function refreshToken(refreshToken: string): Promise<WrapperResponse<RefreshTokenResponse>> {
  return fetchWithErrorHandling(AUTH_ENDPOINTS.REFRESH_TOKEN, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ refresh_token: refreshToken }),
  });
}

export async function requestPasswordReset(email: string): Promise<WrapperResponse<ForgotPasswordResponse>> {
  return fetchWithErrorHandling(AUTH_ENDPOINTS.FORGOT_PASSWORD, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email }),
  });
}

export async function resetPassword(new_password: string): Promise<WrapperResponse<ResetPasswordResponse>> {
  return fetchWithErrorHandling(AUTH_ENDPOINTS.RESET_PASSWORD, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ new_password }),
  });
}

export async function verifyToken(): Promise<WrapperResponse<VerifyTokenResponse>> {
  return fetchWithErrorHandling(AUTH_ENDPOINTS.VERIFY_TOKEN, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    credentials: "include"
  });
}