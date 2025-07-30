/**
 * Authentication API service
 * Handles all authentication-related API calls
 */

/**
 * Login user with email and password
 * @param email - User email
 * @param password - User password
 * @returns Promise with login response data
 */
export async function loginUser(email: string, password: string) {
  const response = await fetch("https://api.histweetyy.cc/api/v1/auth/login", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ email, password }),
  });

  const data = await response.json();

  if (!response.ok) {
    throw new Error(data.message || "Login failed");
  }

  return data;
}

/**
 * Handle Google OAuth callback
 * @param code - Authorization code from Google
 * @returns Promise with login response data
 */
export async function handleGoogleCallback(code: string) {
  const response = await fetch("https://api.histweetyy.cc/api/v1/auth/google/callback", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ code }),
  });

  const data = await response.json();

  if (!response.ok) {
    throw new Error(data.message || "Google authentication failed");
  }

  return data;
}

/**
 * Register new user with email and password
 * @param email - User email
 * @param password - User password
 * @returns Promise with registration response data
 */
export async function registerUser(email: string, password: string) {
  const response = await fetch("https://api.histweetyy.cc/api/v1/auth/register", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ email, password }),
  });

  const data = await response.json();

  if (!response.ok) {
    throw new Error(data.message || "Registration failed");
  }

  return data;
}

/**
 * Request password reset email
 * @param email - User email
 * @returns Promise with request response data
 */
export async function requestPasswordReset(email: string) {
  const response = await fetch("https://api.histweetyy.cc/api/v1/auth/forgot-password", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ email }),
  });

  const data = await response.json();

  if (!response.ok) {
    throw new Error(data.message || "Password reset request failed");
  }

  return data;
}

/**
 * Reset password with token
 * @param token - Reset password token
 * @param password - New password
 * @returns Promise with reset response data
 */
export async function resetPassword(token: string, password: string) {
  const response = await fetch("https://api.histweetyy.cc/api/v1/auth/reset-password", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ token, new_password: password }),
  });

  const data = await response.json();

  if (!response.ok) {
    throw new Error(data.message || "Password reset failed");
  }

  return data;
}
