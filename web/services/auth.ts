import { AUTH_ENDPOINTS } from "@/constants/api";
import { fetchWithErrorHandling } from "@/utils";

export async function loginUser(email: string, password: string) {
  return fetchWithErrorHandling(AUTH_ENDPOINTS.LOGIN, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
  });
}

export async function handleGoogleCallback(code: string) {
  return fetchWithErrorHandling(AUTH_ENDPOINTS.GOOGLE_CALLBACK, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ code }),
  });
}

export async function registerUser(email: string, password: string) {
  return fetchWithErrorHandling(AUTH_ENDPOINTS.REGISTER, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
  });
}

export async function requestPasswordReset(email: string) {
  return fetchWithErrorHandling(AUTH_ENDPOINTS.FORGOT_PASSWORD, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email }),
  });
}

export async function resetPassword(token: string, password: string) {
  return fetchWithErrorHandling(AUTH_ENDPOINTS.RESET_PASSWORD, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ token, new_password: password }),
  });
}