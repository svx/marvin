import { type ClassValue, clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function formatDate(dateString: string): string {
  const date = new Date(dateString);
  return new Intl.DateTimeFormat('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  }).format(date);
}

export function formatRelativeTime(dateString: string): string {
  const date = new Date(dateString);
  const now = new Date();
  const diffInSeconds = Math.floor((now.getTime() - date.getTime()) / 1000);

  if (diffInSeconds < 60) {
    return 'just now';
  }

  const diffInMinutes = Math.floor(diffInSeconds / 60);
  if (diffInMinutes < 60) {
    return `${diffInMinutes} minute${diffInMinutes > 1 ? 's' : ''} ago`;
  }

  const diffInHours = Math.floor(diffInMinutes / 60);
  if (diffInHours < 24) {
    return `${diffInHours} hour${diffInHours > 1 ? 's' : ''} ago`;
  }

  const diffInDays = Math.floor(diffInHours / 24);
  if (diffInDays < 7) {
    return `${diffInDays} day${diffInDays > 1 ? 's' : ''} ago`;
  }

  return formatDate(dateString);
}

export function getSeverityColor(severity: 'error' | 'warning' | 'info'): string {
  switch (severity) {
    case 'error':
      return 'text-error-600 bg-error-50 border-error-200';
    case 'warning':
      return 'text-warning-700 bg-warning-50 border-warning-200';
    case 'info':
      return 'text-primary-600 bg-primary-50 border-primary-200';
    default:
      return 'text-muted bg-gray-50 border-gray-200';
  }
}

export function getSeverityBadgeColor(severity: 'error' | 'warning' | 'info'): string {
  switch (severity) {
    case 'error':
      return 'bg-error-500 text-white';
    case 'warning':
      return 'bg-warning-500 text-gray-900';
    case 'info':
      return 'bg-primary-500 text-white';
    default:
      return 'bg-gray-500 text-white';
  }
}

export function groupIssuesByFile<T extends { file: string }>(issues: T[]): Record<string, T[]> {
  return issues.reduce((acc, issue) => {
    if (!acc[issue.file]) {
      acc[issue.file] = [];
    }
    acc[issue.file].push(issue);
    return acc;
  }, {} as Record<string, T[]>);
}

export function calculatePassRate(totalFiles: number, filesWithIssues: number): number {
  if (totalFiles === 0) return 100;
  return Math.round(((totalFiles - filesWithIssues) / totalFiles) * 100);
}
