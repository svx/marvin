import type { ResultWithId, ResultsResponse } from './types';

export async function fetchResults(): Promise<ResultWithId[]> {
  // Add cache-busting to ensure we get fresh data
  const response = await fetch(`/api/results?_t=${Date.now()}`, {
    cache: 'no-store',
  });
  if (!response.ok) {
    throw new Error('Failed to fetch results');
  }
  const data: ResultsResponse = await response.json();
  return data.results;
}

export async function fetchResultById(id: string): Promise<ResultWithId> {
  // Add cache-busting to ensure we get fresh data
  const response = await fetch(`/api/results/${id}?_t=${Date.now()}`, {
    cache: 'no-store',
  });
  if (!response.ok) {
    throw new Error(`Failed to fetch result ${id}`);
  }
  return response.json();
}

export async function runCheck(
  checker: 'vale' | 'markdownlint',
  path?: string
): Promise<{ success: boolean; message: string }> {
  const response = await fetch('/api/run-check', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ checker, path }),
  });
  
  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.details || error.error || 'Failed to run check');
  }
  
  return response.json();
}
