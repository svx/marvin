'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import Card, { CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import Badge from '@/components/ui/badge';
import { fetchResults } from '@/lib/api';
import { formatRelativeTime } from '@/lib/utils';
import type { ResultWithId } from '@/lib/types';

export default function AllChecksPage() {
  const [results, setResults] = useState<ResultWithId[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [filterChecker, setFilterChecker] = useState<string>('all');

  useEffect(() => {
    async function loadResults() {
      try {
        const data = await fetchResults();
        setResults(data);
      } catch (err) {
        setError('Failed to load results');
        console.error(err);
      } finally {
        setLoading(false);
      }
    }
    loadResults();
  }, []);

  if (loading) {
    return (
      <div className="container mx-auto px-6 py-8">
        <div className="text-center py-12">
          <p className="text-muted">Loading...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="container mx-auto px-6 py-8">
        <div className="text-center py-12">
          <p className="text-error-600">{error}</p>
        </div>
      </div>
    );
  }

  // Filter results
  const filteredResults = filterChecker === 'all' 
    ? results 
    : results.filter(r => r.checker === filterChecker);

  // Get unique checkers
  const checkers = Array.from(new Set(results.map(r => r.checker)));

  return (
    <div className="container mx-auto px-6 py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-foreground mb-2">All Checks</h1>
        <p className="text-muted">View and filter all documentation quality checks</p>
      </div>

      {/* Filter Tabs */}
      <div className="mb-6 border-b border-border">
        <div className="flex gap-4">
          <button
            onClick={() => setFilterChecker('all')}
            className={`px-4 py-2 border-b-2 transition-colors ${
              filterChecker === 'all'
                ? 'border-primary-500 text-primary-600 font-medium'
                : 'border-transparent text-muted hover:text-foreground'
            }`}
          >
            All ({results.length})
          </button>
          {checkers.map(checker => {
            const count = results.filter(r => r.checker === checker).length;
            return (
              <button
                key={checker}
                onClick={() => setFilterChecker(checker)}
                className={`px-4 py-2 border-b-2 transition-colors capitalize ${
                  filterChecker === checker
                    ? 'border-primary-500 text-primary-600 font-medium'
                    : 'border-transparent text-muted hover:text-foreground'
                }`}
              >
                {checker} ({count})
              </button>
            );
          })}
        </div>
      </div>

      {/* Results List */}
      {filteredResults.length === 0 ? (
        <Card>
          <CardContent>
            <p className="text-muted text-center py-8">No checks found</p>
          </CardContent>
        </Card>
      ) : (
        <div className="space-y-4">
          {filteredResults.map((result) => (
            <Link
              key={result.id}
              href={`/checks/${result.checker}?id=${result.id}`}
              className="block"
            >
              <Card hover>
                <CardContent>
                  <div className="flex items-center justify-between mb-3">
                    <div className="flex items-center gap-3">
                      <span className="text-3xl">
                        {result.checker === 'vale' ? '📝' : '📄'}
                      </span>
                      <div>
                        <h3 className="font-semibold text-foreground capitalize">
                          {result.checker}
                        </h3>
                        <p className="text-sm text-muted">
                          {formatRelativeTime(result.timestamp)}
                        </p>
                      </div>
                    </div>
                    <div className="flex gap-2">
                      {result.summary.error_count > 0 && (
                        <Badge severity="error">{result.summary.error_count} errors</Badge>
                      )}
                      {result.summary.warning_count > 0 && (
                        <Badge severity="warning">{result.summary.warning_count} warnings</Badge>
                      )}
                      {result.summary.info_count > 0 && (
                        <Badge severity="info">{result.summary.info_count} info</Badge>
                      )}
                    </div>
                  </div>
                  
                  <div className="grid grid-cols-3 gap-4 text-sm">
                    <div>
                      <span className="text-muted">Files Checked:</span>
                      <span className="ml-2 font-medium text-foreground">
                        {result.summary.total_files}
                      </span>
                    </div>
                    <div>
                      <span className="text-muted">Files with Issues:</span>
                      <span className="ml-2 font-medium text-foreground">
                        {result.summary.files_with_issues}
                      </span>
                    </div>
                    <div>
                      <span className="text-muted">Total Issues:</span>
                      <span className="ml-2 font-medium text-foreground">
                        {result.summary.total_issues}
                      </span>
                    </div>
                  </div>

                  {result.path && (
                    <div className="mt-3 text-sm text-muted">
                      Path: <code className="text-xs bg-gray-100 px-2 py-1 rounded">{result.path}</code>
                    </div>
                  )}
                </CardContent>
              </Card>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}
