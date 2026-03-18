// API client module to handles all HTTP requests to the backend

const BASE_URL = "/api";

/**
 * Fetch today's assigned words
 * @returns {Promise<{date: string, count: number, words: Array}>}
 */
export async function fetchTodayWords() {
  const response = await fetch(`${BASE_URL}/words/today`);

  if (!response.ok) {
    const body = await response.json().catch(() => ({}));
    throw new Error(body.Error || `HTTP ${response.status}`);
  }

  return response.json();
}

export async function fetchRandomWords() {
  const response = await fetch(`${BASE_URL}/words/random`);

  if (!response.ok) {
    const body = await response.json().catch(() => ({}));
    throw new Error(body.Error || `HTTP ${response.status}`);
  }

  return response.json();
}
