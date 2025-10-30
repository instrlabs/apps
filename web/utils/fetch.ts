"use server";

import { cookies, headers } from "next/headers";
import { ResponseCookie, ResponseCookies } from "next/dist/compiled/@edge-runtime/cookies";

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

const ACCESS_TOKEN = "access_token";
const REFRESH_TOKEN = "refresh_token";

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

async function updateCookiesFromResponse(res: Response): Promise<void> {
  const cookieStore = await cookies();

  const setCookieHeaders = res.headers.getSetCookie?.() || [];
  console.log("[updateCookiesFromResponse] Set-Cookie headers count:", setCookieHeaders.length);
  console.log("[updateCookiesFromResponse] Set-Cookie headers:", setCookieHeaders);

  if (setCookieHeaders.length === 0) {
    console.log("[updateCookiesFromResponse] No Set-Cookie headers found");
    return;
  }

  for (const setCookieHeader of setCookieHeaders) {
    console.log("[updateCookiesFromResponse] Processing Set-Cookie:", setCookieHeader.substring(0, 50));

    const cookieName = setCookieHeader.split("=")[0].trim();
    console.log("[updateCookiesFromResponse] Cookie name:", cookieName);

    if (cookieName === ACCESS_TOKEN || cookieName === REFRESH_TOKEN) {
      try {
        const resCookies = new ResponseCookies(new Headers({ "set-cookie": setCookieHeader }));
        const cookie = resCookies.get(cookieName);

        if (cookie) {
          console.log("[updateCookiesFromResponse] Setting cookie:", cookieName);
          cookieStore.set(cookie as ResponseCookie);
        }
      } catch (error) {
        console.error("[updateCookiesFromResponse] Error setting cookie:", cookieName, error);
      }
    }
  }
}

export async function fetchGET<T>(
  path: string,
  queries: Record<string, string> = {},
): Promise<ApiResponse<T>> {
  let url = process.env.API_URL + path;
  const params = new URLSearchParams(queries);
  if (queries) url += "?" + params.toString();

  const headers = await getHeaders();
  headers.set("content-type", "application/json");
  const res = await fetch(url, {
    method: "GET",
    headers: headers,
    credentials: "include",
  });

  console.log("[fetchGET] URL:", url, "Response status:", res.status);
  console.log("[fetchGET] Response headers:", Array.from(res.headers.entries()));

  await updateCookiesFromResponse(res);

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
  const url = process.env.API_URL + path;

  const headers = await getHeaders();
  headers.set("content-type", "application/json");
  const res = await fetch(url, {
    method: "POST",
    headers: headers,
    body: JSON.stringify(body),
    credentials: "include",
  });

  console.log("[fetchPOST] URL:", url, "Response status:", res.status);

  await updateCookiesFromResponse(res);

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
  const url = process.env.API_URL + path;

  const headers = await getHeaders();
  headers.set("content-type", "application/json");
  const res = await fetch(url, {
    method: "PUT",
    headers: headers,
    body: JSON.stringify(body),
    credentials: "include",
  });

  await updateCookiesFromResponse(res);

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
  const url = process.env.API_URL + path;

  const headers = await getHeaders();
  headers.set("content-type", "application/json");
  const res = await fetch(url, {
    method: "PATCH",
    headers: headers,
    body: JSON.stringify(body),
    credentials: "include",
  });

  await updateCookiesFromResponse(res);

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
  const url = process.env.API_URL + path;

  const res = await fetch(url, {
    method: "GET",
    headers: await getHeaders(),
    credentials: "include",
  });

  await updateCookiesFromResponse(res);

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
    } catch {
      // ignore parse error
    }
    return {
      success: false,
      message,
      data: null,
      errors,
    };
  }
}

export async function fetchPOSTFormData<T>(
  path: string,
  formData: FormData,
): Promise<ApiResponse<T>> {
  const url = process.env.API_URL + path;

  const res = await fetch(url, {
    method: "POST",
    headers: await getHeaders(),
    body: formData,
    credentials: "include",
  });

  await updateCookiesFromResponse(res);

  const isOK = res.ok;
  const resBody = await res.json();
  return {
    success: isOK,
    message: resBody.message,
    data: isOK ? resBody.data : null,
    errors: !isOK ? resBody.errors : null,
  };
}
