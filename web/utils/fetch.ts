"use server"

import {cookies} from "next/headers";

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

export async function fetchGET<T>(
  path: string,
  queries: Record<string, string> = {}
): Promise<ApiResponse<T>> {
  let url = process.env.API_URL + path;

  const params = new URLSearchParams(queries);
  if (queries) url += "?" + params.toString();

  const storeCookie = await cookies()

  const res = await fetch(url, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      "Cookie": storeCookie.toString()
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
  body: unknown
): Promise<ApiResponse<T>> {
  const url = process.env.API_URL + path;

  const res = await fetch(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "Cookie": (await cookies()).toString()
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
  const url = process.env.API_URL + path;

  const res = await fetch(url, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
      "Cookie": (await cookies()).toString()
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
  const url = process.env.API_URL + path;

  const res = await fetch(url, {
    method: "PATCH",
    headers: {
      "Content-Type": "application/json",
      "Cookie": (await cookies()).toString()
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
  let url = process.env.API_URL + path;
  const params = new URLSearchParams(queries);
  if (queries && Object.keys(queries).length > 0) url += "?" + params.toString();

  const res = await fetch(url, {
    method: "GET",
    headers: {
      // Let the server specify Content-Type; we just forward auth cookie
      "Cookie": (await cookies()).toString(),
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
  const url = process.env.API_URL + path;

  const res = await fetch(url, {
    method: "POST",
    headers: {
      "Cookie": (await cookies()).toString(),
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

export async function fetchEventStream(
  path: string
): Promise<Response> {
  const url = "http://localhost:3001" + path;

  return await fetch(url, {
    method: "GET",
    headers: {
      Accept: "text/event-stream",
      Connection: "keep-alive",
      "Cache-Control": "no-cache",
      "Cookie": (await cookies()).toString(),
    },
    cache: "no-store"
  });
}
