/**
 * API Constants
 * Contains all API-related constants
 */

// Base URL for all API requests
export const API_BASE_URL = "http://gateway-service:3000";

// Auth API endpoints
export const AUTH_ENDPOINTS = {
  LOGIN: `${API_BASE_URL}/api/v1/auth/login`,
  REGISTER: `${API_BASE_URL}/api/v1/auth/register`,
  GOOGLE_CALLBACK: `${API_BASE_URL}/api/v1/auth/google/callback`,
  FORGOT_PASSWORD: `${API_BASE_URL}/api/v1/auth/forgot-password`,
  RESET_PASSWORD: `${API_BASE_URL}/api/v1/auth/reset-password`,
};