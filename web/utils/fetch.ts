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
