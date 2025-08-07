/**
 * Utility functions for fetch operations
 */

/**
 * Wrapper function for fetch operations that handles errors internally
 * @param url - The URL to fetch from
 * @param options - The fetch options
 * @returns An object with data and error properties
 */
export async function fetchWithErrorHandling(url: string, options: RequestInit) {
  try {
    const response = await fetch(url, options);
    const data = await response.json();

    if (!response.ok) {
      throw new Error(data.message || "Request failed");
    }

    return { data: data, error: null };
  } catch (err) {
    return { 
      data: null, 
      error: err instanceof Error ? err.message : "An unexpected error occurred" 
    };
  }
}