'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import Card, { CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import Badge from '@/components/ui/badge';
import RunCheckButton from '@/components/ui/run-check-button';
import { fetchResults } from '@/lib/api';
import { formatRelativeTime, calculatePassRate } from '@/lib/utils';
import type { ResultWithId, Summary } from '@/lib/types';

export default function DashboardPage() {
  const [results, setResults] = useState<ResultWithId[]>([]);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [error, setError] = useState<string | null>(null);

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

  // Calculate aggregate statistics
  const totalChecks = results.length;
  const latestResults = results.slice(0, 5);
  
  const aggregateStats = results.reduce(
    (acc, result) => ({
      totalIssues: acc.totalIssues + result.summary.total_issues,
      errorCount: acc.errorCount + result.summary.error_count,
      warningCount: acc.warningCount + result.summary.warning_count,
      infoCount: acc.infoCount + result.summary.info_count,
      totalFiles: acc.totalFiles + result.summary.total_files,
      filesWithIssues: acc.filesWithIssues + result.summary.files_with_issues,
    }),
    { totalIssues: 0, errorCount: 0, warningCount: 0, infoCount: 0, totalFiles: 0, filesWithIssues: 0 }
  );

  const passRate = totalChecks > 0 
    ? calculatePassRate(aggregateStats.totalFiles, aggregateStats.filesWithIssues)
    : 0;

  const handleCheckComplete = async () => {
    // Reload results after check completes
    setRefreshing(true);
    try {
      const data = await fetchResults();
      setResults(data);
    } catch (err) {
      console.error('Failed to reload results:', err);
    } finally {
      setRefreshing(false);
    }
  };

  return (
    <div className="container mx-auto px-6 py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-foreground mb-2">Dashboard</h1>
        <p className="text-muted">Overview of documentation quality checks</p>
      </div>

      {/* Run Checks Section */}
      <Card className="mb-8">
        <CardHeader>
          <CardTitle>Run Quality Checks</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex flex-wrap gap-4">
            <RunCheckButton
              checker="vale"
              onSuccess={handleCheckComplete}
              variant="primary"
            />
            <RunCheckButton
              checker="markdownlint"
              onSuccess={handleCheckComplete}
              variant="primary"
            />
          </div>
          <p className="text-sm text-muted mt-4">
            Run documentation quality checks directly from the dashboard. Results will appear below once completed.
          </p>
        </CardContent>
      </Card>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <Card>
          <CardContent>
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted mb-1">Total Checks</p>
                <p className="text-3xl font-bold text-foreground">{totalChecks}</p>
              </div>
              <div className="text-4xl">📊</div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent>
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted mb-1">Total Issues</p>
                <p className="text-3xl font-bold text-foreground">{aggregateStats.totalIssues}</p>
              </div>
              <div className="text-4xl">🔍</div>
            </div>
            <div className="mt-3 flex gap-2">
              <Badge severity="error">{aggregateStats.errorCount}</Badge>
              <Badge severity="warning">{aggregateStats.warningCount}</Badge>
              <Badge severity="info">{aggregateStats.infoCount}</Badge>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent>
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted mb-1">Files Checked</p>
                <p className="text-3xl font-bold text-foreground">{aggregateStats.totalFiles}</p>
              </div>
              <div className="text-4xl">📄</div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent>
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted mb-1">Pass Rate</p>
                <p className="text-3xl font-bold text-foreground">{passRate}%</p>
              </div>
              <div className="text-4xl">✅</div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Recent Checks */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <CardTitle>Recent Checks</CardTitle>
            <Link href="/checks" className="text-sm text-primary-600 hover:text-primary-700">
              View all →
            </Link>
          </div>
        </CardHeader>
        <CardContent>
          {refreshing && (
            <div className="mb-4 p-3 bg-blue-50 border border-blue-200 rounded-md flex items-center gap-2 text-blue-700 text-sm">
              <svg
                className="animate-spin h-4 w-4"
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
              >
                <circle
                  className="opacity-25"
                  cx="12"
                  cy="12"
                  r="10"
                  stroke="currentColor"
                  strokeWidth="4"
                />
                <path
                  className="opacity-75"
                  fill="currentColor"
                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                />
              </svg>
              <span>Refreshing results...</span>
            </div>
          )}
          
          {latestResults.length === 0 ? (
            <p className="text-muted text-center py-8">No checks found</p>
          ) : (
            <div className="space-y-4">
              {latestResults.map((result) => (
                <Link
                  key={result.id}
                  href={`/checks/${result.checker}?id=${result.id}`}
                  className="block p-4 border border-border rounded-stripe hover:border-primary-300 hover:bg-primary-50/50 transition-colors"
                >
                  <div className="flex items-center justify-between mb-2">
                    <div className="flex items-center gap-3">
                      <span className="text-2xl">
                        {result.checker === 'vale' ? '📝' : '📄'}
                      </span>
                      <div>
                        <h4 className="font-medium text-foreground capitalize">
                          {result.checker}
                        </h4>
                        <p className="text-sm text-muted">
                          {formatRelativeTime(result.timestamp)}
                        </p>
                      </div>
                    </div>
                    <div className="flex gap-2">
                      {result.summary.error_count > 0 && (
                        <Badge severity="error">{result.summary.error_count}</Badge>
                      )}
                      {result.summary.warning_count > 0 && (
                        <Badge severity="warning">{result.summary.warning_count}</Badge>
                      )}
                      {result.summary.info_count > 0 && (
                        <Badge severity="info">{result.summary.info_count}</Badge>
                      )}
                    </div>
                  </div>
                  <div className="text-sm text-muted">
                    {result.summary.total_files} files checked, {result.summary.files_with_issues} with issues
                  </div>
                </Link>
              ))}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
