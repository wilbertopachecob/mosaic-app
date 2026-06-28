/**
 * Reads a fetch response body as text and parses it as JSON.
 * Returns null when the body is empty or not valid JSON.
 */
export async function parseJson<T>(response: Response): Promise<T | null> {
  const text = await response.text();
  if (!text) {
    return null;
  }
  try {
    return JSON.parse(text) as T;
  } catch {
    return null;
  }
}
