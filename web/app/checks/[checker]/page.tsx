'use client';

import { useEffect, useState } from 'react';
import { useSearchParams, useRouter } from 'next/navigation';
import Link from 'next/link';
import Card, { CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import Badge from '@/components/ui/badge';
import RunCheckButton from '@/components/ui/run-check-button';
import { fetchResultById, fetchResults } from '@/lib/api';
import { formatDate, groupIssuesByFile } from '@/lib/utils';
import type { ResultWithId, Issue } from '@/lib/types';

export default function CheckerDetailPage({ params }: { params: { checker: string } }) {
  const searchParams = useSearchParams();
  const router = useRouter();
  const id = searchParams.get('id');
  
  const [result, setResult] = useState<ResultWithId | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedFile, setSelectedFile] = useState<string | null>(null);

  useEffect(() => {
    async function loadResult() {
      if (!id) {
        setError('No result ID provided');
        setLoading(false);
        return;
      }

      try {
        const data = await fetchResultById(id);
        setResult(data);
      } catch (err) {
        setError('Failed to load result');
        console.error(err);
      } finally {
        setLoading(false);
      }
    }
    loadResult();
  }, [id]);

  if (loading) {
    return (
      <div className="container mx-auto px-6 py-8">
        <div className="text-center py-12">
          <p className="text-muted">Loading...</p>
        </div>
      </div>
    );
  }

  if (error || !result) {
    return (
      <div className="container mx-auto px-6 py-8">
        <div className="text-center py-12">
          <p className="text-error-600">{error || 'Result not found'}</p>
          <Link href="/checks" className="text-primary-600 hover:text-primary-700 mt-4 inline-block">
            ← Back to all checks
          </Link>
        </div>
      </div>
    );
  }

  const issuesByFile = groupIssuesByFile(result.issues);
  const files = Object.keys(issuesByFile).sort();
  const displayIssues = selectedFile ? issuesByFile[selectedFile] : result.issues;

  const handleCheckComplete = async () => {
    // Reload to get the latest result for this checker
    try {
      const allResults = await fetchResults();
      const latestResult = allResults
        .filter(r => r.checker === result.checker)
        .sort((a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime())[0];
      
      if (latestResult && latestResult.id !== result.id) {
        // Navigate to the latest result
        router.push(`/checks/${latestResult.checker}?id=${latestResult.id}`);
      } else if (latestResult) {
        // Same result, just reload the data
        const updatedResult = await fetchResultById(latestResult.id);
        setResult(updatedResult);
      }
    } catch (err) {
      console.error('Failed to reload results:', err);
    }
  };

  // Use the exact path from the result for re-running
  // The CLI runs from cli/ directory, so we need to prepend ../ if the path doesn't already have it
  const storedPath = result.path;
  const rerunPath = storedPath.startsWith('../') || storedPath.startsWith('/') ? storedPath : `../${storedPath}`;
  
  // Extract a display name from the path
  const getDisplayName = (path: string): string => {
    // Remove ../ or / prefix for display
    const cleanPath = path.replace(/^\.\.\//, '').replace(/^\//, '');
    // If it's a file, show just the filename
    if (cleanPath.endsWith('.md')) {
      return cleanPath.split('/').pop() || cleanPath;
    }
    // If it's a directory, show the last directory name
    return cleanPath.split('/').pop() || cleanPath;
  };
  
  const displayName = getDisplayName(storedPath);
  
  console.log('Checker detail page:', { storedPath, rerunPath, displayName });

  return (
    <div className="container mx-auto px-6 py-8">
      {/* Header */}
      <div className="mb-6">
        <Link href="/checks" className="text-primary-600 hover:text-primary-700 text-sm mb-2 inline-block">
          ← Back to all checks
        </Link>
        <div className="flex items-start justify-between">
          <div>
            <h1 className="text-3xl font-bold text-foreground mb-2 capitalize">
              {result.checker} Check Results
            </h1>
            <p className="text-muted">{formatDate(result.timestamp)}</p>
            <p className="text-sm text-muted mt-1">
              Checked path: <code className="text-xs bg-gray-100 px-2 py-0.5 rounded">{result.path}</code>
            </p>
          </div>
          <div className="flex gap-2">
            <RunCheckButton
              checker={result.checker as 'vale' | 'markdownlint'}
              onSuccess={handleCheckComplete}
              variant="secondary"
              size="sm"
              path={rerunPath}
              label={`Re-run (${displayName})`}
            />
            <RunCheckButton
              checker={result.checker as 'vale' | 'markdownlint'}
              onSuccess={handleCheckComplete}
              variant="primary"
              size="sm"
              label="Check all docs"
            />
          </div>
        </div>
      </div>

      {/* Summary Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-8">
        <Card>
          <CardContent>
            <p className="text-sm text-muted mb-1">Total Files</p>
            <p className="text-2xl font-bold text-foreground">{result.summary.total_files}</p>
          </CardContent>
        </Card>
        <Card>
          <CardContent>
            <p className="text-sm text-muted mb-1">Files with Issues</p>
            <p className="text-2xl font-bold text-foreground">{result.summary.files_with_issues}</p>
          </CardContent>
        </Card>
        <Card>
          <CardContent>
            <p className="text-sm text-muted mb-1">Total Issues</p>
            <p className="text-2xl font-bold text-foreground">{result.summary.total_issues}</p>
          </CardContent>
        </Card>
        <Card>
          <CardContent>
            <p className="text-sm text-muted mb-1">By Severity</p>
            <div className="flex gap-2 mt-2">
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
          </CardContent>
        </Card>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
        {/* File List Sidebar */}
        <div className="lg:col-span-1">
          <Card>
            <CardHeader>
              <CardTitle className="text-base">Files ({files.length})</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-1">
                <button
                  onClick={() => setSelectedFile(null)}
                  className={`w-full text-left px-3 py-2 rounded text-sm transition-colors ${
                    selectedFile === null
                      ? 'bg-primary-50 text-primary-700 font-medium'
                      : 'text-muted hover:bg-gray-50'
                  }`}
                >
                  All Files ({result.issues.length})
                </button>
                {files.map(file => (
                  <button
                    key={file}
                    onClick={() => setSelectedFile(file)}
                    className={`w-full text-left px-3 py-2 rounded text-sm transition-colors ${
                      selectedFile === file
                        ? 'bg-primary-50 text-primary-700 font-medium'
                        : 'text-muted hover:bg-gray-50'
                    }`}
                  >
                    <div className="truncate">{file.split('/').pop()}</div>
                    <div className="text-xs text-muted mt-1">
                      {issuesByFile[file].length} issues
                    </div>
                  </button>
                ))}
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Issues List */}
        <div className="lg:col-span-3">
          <Card>
            <CardHeader>
              <CardTitle>
                {selectedFile ? `Issues in ${selectedFile}` : 'All Issues'} ({displayIssues.length})
              </CardTitle>
            </CardHeader>
            <CardContent>
              {displayIssues.length === 0 ? (
                <p className="text-muted text-center py-8">No issues found</p>
              ) : (
                <div className="space-y-4">
                  {displayIssues.map((issue, index) => (
                    <div
                      key={index}
                      className="border border-border rounded-stripe p-4 hover:border-primary-300 transition-colors"
                    >
                      <div className="flex items-start justify-between mb-2">
                        <div className="flex-1">
                          <div className="flex items-center gap-2 mb-1">
                            <Badge severity={issue.severity}>{issue.severity}</Badge>
                            <code className="text-xs bg-gray-100 px-2 py-1 rounded">
                              {issue.rule}
                            </code>
                          </div>
                          <p className="text-sm font-medium text-foreground mb-1">
                            {issue.message}
                          </p>
                          <p className="text-xs text-muted">
                            {issue.file} • Line {issue.line}, Column {issue.column}
                          </p>
                        </div>
                      </div>
                      {issue.context && (
                        <div className="mt-3 p-3 bg-gray-50 rounded border border-gray-200">
                          <p className="text-xs text-muted mb-1">Context:</p>
                          <code className="text-sm text-foreground">{issue.context}</code>
                        </div>
                      )}
                    </div>
                  ))}
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
