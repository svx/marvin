// TypeScript types matching the CLI's Go models

export interface Result {
  checker: string;
  timestamp: string;
  path: string;
  summary: Summary;
  issues: Issue[];
  metadata: Record<string, any>;
}

export interface Summary {
  total_files: number;
  files_with_issues: number;
  total_issues: number;
  error_count: number;
  warning_count: number;
  info_count: number;
}

export interface Issue {
  file: string;
  line: number;
  column: number;
  severity: 'error' | 'warning' | 'info';
  message: string;
  rule: string;
  context?: string;
}

export interface CheckerType {
  id: string;
  name: string;
  description: string;
  icon: string;
}

export interface ResultWithId extends Result {
  id: string;
}

export interface ResultsResponse {
  results: ResultWithId[];
  total: number;
  page: number;
  pageSize: number;
}

export interface IssuesByFile {
  [filename: string]: Issue[];
}

export type SeverityType = 'error' | 'warning' | 'info';

export interface FilterOptions {
  checker?: string;
  severity?: SeverityType;
  dateFrom?: string;
  dateTo?: string;
  searchTerm?: string;
}
