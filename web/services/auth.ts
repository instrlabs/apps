import { AUTH_ENDPOINTS } from "@/constants/api";

export async function loginUser(email: string, password: string) {
  const response = await fetch(AUTH_ENDPOINTS.LOGIN, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
  });

  const data = await response.json();

  if (!response.ok) {
    throw new Error(data.message || "Login failed");
  }

  return data;
}

export async function handleGoogleCallback(code: string) {
  const response = await fetch(AUTH_ENDPOINTS.GOOGLE_CALLBACK, {
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

export async function registerUser(email: string, password: string) {
  const response = await fetch(AUTH_ENDPOINTS.REGISTER, {
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

export async function requestPasswordReset(email: string) {
  const response = await fetch(AUTH_ENDPOINTS.FORGOT_PASSWORD, {
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

export async function resetPassword(token: string, password: string) {
  const response = await fetch(AUTH_ENDPOINTS.RESET_PASSWORD, {
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