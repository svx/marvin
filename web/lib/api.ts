import type { ResultWithId, ResultsResponse } from './types';

export async function fetchResults(): Promise<ResultWithId[]> {
  const response = await fetch('/api/results');
  if (!response.ok) {
    throw new Error('Failed to fetch results');
  }
  const data: ResultsResponse = await response.json();
  return data.results;
}

export async function fetchResultById(id: string): Promise<ResultWithId> {
  const response = await fetch(`/api/results/${id}`);
  if (!response.ok) {
    throw new Error(`Failed to fetch result ${id}`);
  }
  return response.json();
}
