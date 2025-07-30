/**
 * API Constants
 * Contains all API-related constants
 */

// Base URL for all API requests
export const API_BASE_URL = "http://gateway-service.localhost";

// Auth API endpoints
export const AUTH_ENDPOINTS = {
  LOGIN: `${API_BASE_URL}/auth/login`,
  REGISTER: `${API_BASE_URL}/auth/register`,
  GOOGLE_CALLBACK: `${API_BASE_URL}/auth/google/callback`,
  FORGOT_PASSWORD: `${API_BASE_URL}/auth/forgot-password`,
  RESET_PASSWORD: `${API_BASE_URL}/auth/reset-password`,
};