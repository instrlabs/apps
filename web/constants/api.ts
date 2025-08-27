// Base URL for all API requests
export const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://gateway-service.localhost";

// Auth API endpoints
export const AUTH_ENDPOINTS = {
  LOGIN: `${API_BASE_URL}/auth/login`,
  REGISTER: `${API_BASE_URL}/auth/register`,
  GOOGLE: `${API_BASE_URL}/auth/google`,
  GOOGLE_CALLBACK: `${API_BASE_URL}/auth/google/callback`,
  FORGOT_PASSWORD: `${API_BASE_URL}/auth/forgot-password`,
  RESET_PASSWORD: `${API_BASE_URL}/auth/reset-password`,
  REFRESH_TOKEN: `${API_BASE_URL}/auth/refresh`,
  VERIFY_TOKEN: `${API_BASE_URL}/auth/verify-token`,
  PROFILE: `${API_BASE_URL}/auth/profile`,
  CHANGE_PASSWORD: `${API_BASE_URL}/auth/change-password`,
  LOGOUT: `${API_BASE_URL}/auth/logout`,
};
