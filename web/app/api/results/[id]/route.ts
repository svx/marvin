import { readFile } from 'fs/promises';
import { join } from 'path';
import { NextResponse } from 'next/server';
import type { Result, ResultWithId } from '@/lib/types';

const RESULTS_DIR = join(process.cwd(), '..', '.marvin', 'results');

export async function GET(
  request: Request,
  { params }: { params: { id: string } }
) {
  try {
    const { id } = params;
    const filePath = join(RESULTS_DIR, `${id}.json`);
    
    const content = await readFile(filePath, 'utf-8');
    const data: Result = JSON.parse(content);
    
    const result: ResultWithId = {
      id,
      ...data
    };
    
    return NextResponse.json(result);
  } catch (error) {
    console.error('Error reading result:', error);
    return NextResponse.json(
      { error: 'Result not found' },
      { status: 404 }
    );
  }
}
