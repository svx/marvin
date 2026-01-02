export const CHECKER_TYPES = {
  vale: {
    id: 'vale',
    name: 'Vale',
    description: 'Prose linting for documentation',
    icon: '📝',
  },
  markdownlint: {
    id: 'markdownlint',
    name: 'Markdownlint',
    description: 'Markdown linting and style checking',
    icon: '📄',
  },
} as const;

export const SEVERITY_LABELS = {
  error: 'Error',
  warning: 'Warning',
  info: 'Info',
} as const;

export const RESULTS_PER_PAGE = 20;

export const DEFAULT_RESULTS_DIR = '../.marvin/results';
