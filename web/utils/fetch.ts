import {AUTH_ENDPOINTS} from "@/constants/api";
import ROUTES from "@/constants/routes";
import {redirect} from "next/navigation";

function redirectToLogin() {
  if (typeof window !== "undefined") {
    window.location.href = ROUTES.LOGIN;
  } else redirect(ROUTES.LOGIN);
}

export async function fetchWithErrorHandling(url: string, options: RequestInit) {
    try {
        const response = await fetch(url, options);

        if (response.status === 401) {
            const refreshResponse = await fetch(AUTH_ENDPOINTS.REFRESH_TOKEN, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                credentials: "include",
            });

            if (!refreshResponse.ok) {
                redirectToLogin();
            }

            const response2 = await fetch(url, options);

            const data = await response2.json();

            return { data, error: null };
        }

        if (!response.ok) {
            const errorBody = await response.json();
            // Pass through standard message and optional field-level errors if present
            return { data: null, error: errorBody?.message ?? "", errors: errorBody?.errors ?? null };
        }

        const data = await response.json();

        return { data, error: null };
    } catch {
        return { data: null, error: "Something went wrong" };
    }
}
