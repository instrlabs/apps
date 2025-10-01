"use server"

import { cookies, headers } from "next/headers";

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

export type EmptyBody = Record<string, unknown>;

export async function getRequestOrigin(): Promise<string> {
  const h = await headers();
  const host = h.get("x-forwarded-host") ?? h.get("host") ?? "localhost:3000";
  const proto = h.get("x-forwarded-proto") ?? (host.startsWith("localhost") ? "http" : "https");
  return `${proto}://${host}`;
}

export async function fetchGET<T>(
  path: string,
  queries: Record<string, string> = {}
): Promise<ApiResponse<T>> {
  let url = process.env.GATEWAY_URL + path;

  const params = new URLSearchParams(queries);
  if (queries) url += "?" + params.toString();

  const storeCookie = await cookies()

  const res = await fetch(url, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      "Cookie": storeCookie.toString(),
      "Origin": await getRequestOrigin(),
    },
  });

  const isOK = res.ok;
  const resBody = await res.json();
  return {
    success: isOK,
    message: resBody.message,
    data: isOK ? (resBody.data as T) : null,
    errors: !isOK ? resBody.errors : null,
  };
}

export async function fetchPOST<T>(
  path: string,
  body?: unknown
): Promise<ApiResponse<T>> {
  const url = process.env.GATEWAY_URL + path;
  const res = await fetch(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "Cookie": (await cookies()).toString(),
      "Origin": await getRequestOrigin(),
    },
    body: JSON.stringify(body),
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

export async function fetchPUT<T>(
  path: string,
  body: unknown
): Promise<ApiResponse<T>> {
  const url = process.env.GATEWAY_URL + path;

  const res = await fetch(url, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
      "Cookie": (await cookies()).toString(),
      "Origin": await getRequestOrigin(),
    },
    body: JSON.stringify(body),
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

export async function fetchPATCH<T>(
  path: string,
  body: unknown
): Promise<ApiResponse<T>> {
  const url = process.env.GATEWAY_URL + path;

  const res = await fetch(url, {
    method: "PATCH",
    headers: {
      "Content-Type": "application/json",
      "Cookie": (await cookies()).toString(),
      "Origin": await getRequestOrigin(),
    },
    body: JSON.stringify(body),
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

export async function fetchGETBytes(
  path: string,
  queries: Record<string, string> = {}
): Promise<ApiResponse<ArrayBuffer>> {
  let url = process.env.GATEWAY_URL + path;
  const params = new URLSearchParams(queries);
  if (queries && Object.keys(queries).length > 0) url += "?" + params.toString();

  const res = await fetch(url, {
    method: "GET",
    headers: {
      "Cookie": (await cookies()).toString(),
      "Origin": await getRequestOrigin(),
    },
  });

  if (res.ok) {
    const data = await res.arrayBuffer();
    return {
      success: true,
      message: "",
      data,
      errors: null,
    };
  }

  try {
    const errJson = await res.json();
    return {
      success: false,
      message: errJson.message ?? res.statusText,
      data: null,
      errors: errJson.errors ?? null,
    };
  } catch {
    const text = await res.text().catch(() => "");
    return {
      success: false,
      message: text || res.statusText,
      data: null,
      errors: null,
    };
  }
}

export async function fetchPOSTFormData<T>(
  path: string,
  formData: FormData
): Promise<ApiResponse<T>> {
  const url = process.env.GATEWAY_URL + path;

  const res = await fetch(url, {
    method: "POST",
    headers: {
      "Cookie": (await cookies()).toString(),
      "Origin": await getRequestOrigin(),
    },
    body: formData,
  });

  const isOK = res.ok;
  const resBody = await res.json();
  return {
    success: isOK,
    message: resBody.message,
    data: isOK ? (resBody.data as T) : null,
    errors: !isOK ? resBody.errors : null,
  };
}
