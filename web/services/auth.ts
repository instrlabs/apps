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

interface FieldError {
  fieldName: string;
  errorMessage: string;
}

interface WrapperResponse<T> {
  data: T | null;
  error: string | null;
  errors?: FieldError[] | null;
}

export async function registerUser(email: string, password: string): Promise<WrapperResponse<RegisterResponse>> {
  return fetchWithErrorHandling(AUTH_ENDPOINTS.REGISTER, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
    credentials: "include"
  });
}

export async function loginUser(email: string, password: string): Promise<WrapperResponse<LoginResponse>> {
  const response = await fetchWithErrorHandling(AUTH_ENDPOINTS.LOGIN, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
    credentials: "include"
  });

  // If login is successful, store the token in localStorage
  if (response.data && !response.error) {
    // For simplicity, we'll use the email as the token
    // In a real application, you would get the actual token from the response
    storeAuthToken(email);
  }

  return response;
}

// Store the authentication token in localStorage
export function storeAuthToken(token: string): void {
  if (typeof window !== 'undefined') {
    localStorage.setItem('auth_token', token);
  }
}

// Clear the authentication token from localStorage
export function clearAuthToken(): void {
  if (typeof window !== 'undefined') {
    localStorage.removeItem('auth_token');
  }
}

export async function refreshToken(): Promise<WrapperResponse<RefreshTokenResponse>> {
  return fetchWithErrorHandling(AUTH_ENDPOINTS.REFRESH_TOKEN, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    credentials: "include"
  });
}

export async function requestPasswordReset(email: string): Promise<WrapperResponse<ForgotPasswordResponse>> {
  return fetchWithErrorHandling(AUTH_ENDPOINTS.FORGOT_PASSWORD, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email }),
    credentials: "include"
  });
}

export async function resetPassword(token: string, new_password: string): Promise<WrapperResponse<ResetPasswordResponse>> {
  return fetchWithErrorHandling(AUTH_ENDPOINTS.RESET_PASSWORD, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ token, new_password }),
    credentials: "include"
  });
}

export async function verifyToken(): Promise<WrapperResponse<VerifyTokenResponse>> {
  return fetchWithErrorHandling(AUTH_ENDPOINTS.VERIFY_TOKEN, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    credentials: "include"
  });
}

// Logout user and clear token
export async function logoutUser(): Promise<void> {
  // In a real application, you would call a logout endpoint
  // For now, we'll just clear the token from localStorage
  clearAuthToken();

  // Disconnect from SSE
  if (typeof window !== 'undefined') {
    const sseService = await import('../services/sse').then(module => module.default);
    sseService.disconnect();
  }
}
