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

/**
 * Fetch random words for extra practice
 * @returns {Promise<{count: number, words: Array}>}
 */
export async function fetchRandomWords() {
  const response = await fetch(`${BASE_URL}/words/random`);

  if (!response.ok) {
    const body = await response.json().catch(() => ({}));
    throw new Error(body.Error || `HTTP ${response.status}`);
  }

  return response.json();
}

/**
 * Generate new words via Gemini and store them in the database
 * @param {number} count How many words to generate (1-20)
 * @param {string} difficulty "easy", "medium", "hard"
 * @returns {Promise<{message: string, generated: number, inserted: number, duplicates: number, words: Array}>}
 */
export async function generateWords(count, difficulty) {
  const response = await fetch(`${BASE_URL}/generate`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ count, difficulty }),
  });

  if (!response.ok) {
    const body = await response.json().catch(() => ({}));
    throw new Error(body.Error || `HTTP ${response.status}`);
  }

  return response.json();
}
