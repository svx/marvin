import { readdir, readFile, access } from 'fs/promises';
import { join } from 'path';
import { NextResponse } from 'next/server';
import type { Result, ResultWithId, ResultsResponse } from '@/lib/types';

const RESULTS_DIR = join(process.cwd(), '..', '.marvin', 'results');

export async function GET(request: Request) {
  try {
    // Check if directory exists
    try {
      await access(RESULTS_DIR);
    } catch {
      // Directory doesn't exist, return empty results
      const response: ResultsResponse = {
        results: [],
        total: 0,
        page: 1,
        pageSize: 20,
      };
      return NextResponse.json(response);
    }

    // Read all files from the results directory
    const files = await readdir(RESULTS_DIR);
    const jsonFiles = files.filter((f: string) => f.endsWith('.json'));
    
    // Read and parse each result file
    const results: ResultWithId[] = await Promise.all(
      jsonFiles.map(async (file: string) => {
        const content = await readFile(join(RESULTS_DIR, file), 'utf-8');
        const data: Result = JSON.parse(content);
        return {
          id: file.replace('.json', ''),
          ...data
        };
      })
    );
    
    // Sort by timestamp, newest first
    results.sort((a, b) => 
      new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime()
    );
    
    // Parse query parameters for filtering
    const { searchParams } = new URL(request.url);
    const checker = searchParams.get('checker');
    const limit = parseInt(searchParams.get('limit') || '20');
    const offset = parseInt(searchParams.get('offset') || '0');
    
    // Filter by checker if specified
    let filteredResults = results;
    if (checker) {
      filteredResults = results.filter(r => r.checker === checker);
    }
    
    // Paginate
    const paginatedResults = filteredResults.slice(offset, offset + limit);
    
    const response: ResultsResponse = {
      results: paginatedResults,
      total: filteredResults.length,
      page: Math.floor(offset / limit) + 1,
      pageSize: limit,
    };
    
    return NextResponse.json(response);
  } catch (error) {
    console.error('Error reading results:', error);
    return NextResponse.json(
      { error: 'Failed to read results', results: [], total: 0, page: 1, pageSize: 20 },
      { status: 500 }
    );
  }
}
