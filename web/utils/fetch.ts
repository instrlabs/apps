"use server";

import { headers } from "next/headers";

export type ApiResponse<TBody> = {
  success: boolean;
  message: string;
  data: TBody | null;
  errors: FormErrors | null;
};

export type FormErrors =
  | {
      errorMessage: string;
      fieldName: string;
    }[]
  | null;

export type EmptyBody = Record<string, unknown>;

async function getHeaders(): Promise<Headers> {
  const h = await headers();
  const customHeaders = new Headers();
  customHeaders.set("x-user-ip", h.get("x-user-ip")!);
  customHeaders.set("x-user-agent", h.get("x-user-agent")!);
  customHeaders.set("x-user-host", h.get("x-user-host")!);
  customHeaders.set("x-user-origin", h.get("x-user-origin")!);
  customHeaders.set("cookie", h.get("cookie")!);

  return customHeaders;
}

export async function fetchGET<T>(path: string, queries: Record<string, string> = {}): Promise<ApiResponse<T>> {
  let url = process.env.GATEWAY_SERVICE + path;
  const params = new URLSearchParams(queries);
  if (queries) url += "?" + params.toString();

  const headers = await getHeaders();
  headers.set("content-type", "application/json");
  const res = await fetch(url, {
    method: "GET",
    headers: headers,
    credentials: "include",
  });

  console.log(`set-cookie: ${res.headers.getSetCookie()}`);

  const isOK = res.ok;
  const resBody = await res.json();

  return {
    success: isOK,
    message: resBody.message,
    data: isOK ? (resBody.data as T) : null,
    errors: !isOK ? resBody.errors : null,
  };
}

export async function fetchPOST<T>(path: string, body?: unknown): Promise<ApiResponse<T>> {
  const url = process.env.GATEWAY_SERVICE + path;

  const headers = await getHeaders();
  headers.set("content-type", "application/json");
  const res = await fetch(url, {
    method: "POST",
    headers: headers,
    body: JSON.stringify(body),
    credentials: "include",
  });

  console.log(`set-cookie: ${res.headers.getSetCookie()}`);

  const isOK = res.ok;
  const resBody = await res.json();
  return {
    success: isOK,
    message: resBody.message,
    data: isOK ? resBody.data : null,
    errors: !isOK ? resBody.errors : null,
  };
}

export async function fetchPUT<T>(path: string, body: unknown): Promise<ApiResponse<T>> {
  const url = process.env.GATEWAY_SERVICE + path;

  const headers = await getHeaders();
  headers.set("content-type", "application/json");
  const res = await fetch(url, {
    method: "PUT",
    headers: headers,
    body: JSON.stringify(body),
    credentials: "include",
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

export async function fetchPATCH<T>(path: string, body: unknown): Promise<ApiResponse<T>> {
  const url = process.env.GATEWAY_SERVICE + path;

  const headers = await getHeaders();
  headers.set("content-type", "application/json");
  const res = await fetch(url, {
    method: "PATCH",
    headers: headers,
    body: JSON.stringify(body),
    credentials: "include",
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

export async function fetchGETBytes(path: string): Promise<ApiResponse<ArrayBuffer>> {
  const url = process.env.GATEWAY_SERVICE + path;

  const res = await fetch(url, {
    method: "GET",
    headers: await getHeaders(),
    credentials: "include",
  });

  const isOK = res.ok;

  if (isOK) {
    const data = await res.arrayBuffer();
    return {
      success: true,
      message: "OK",
      data,
      errors: null,
    };
  } else {
    let message = res.statusText;
    let errors: FormErrors = null;
    try {
      const resBody = await res.json();
      message = resBody.message ?? message;
      errors = resBody.errors ?? null;
    } catch {}
    return {
      success: false,
      message,
      data: null,
      errors,
    };
  }
}

export async function fetchPOSTFormData<T>(path: string, formData: FormData): Promise<ApiResponse<T>> {
  const url = process.env.GATEWAY_SERVICE + path;

  const res = await fetch(url, {
    method: "POST",
    headers: await getHeaders(),
    body: formData,
    credentials: "include",
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
