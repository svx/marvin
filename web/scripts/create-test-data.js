#!/usr/bin/env node

/**
 * Create test data for the Marvin dashboard
 * This script creates sample result files in .marvin/results/
 */

const fs = require('fs');
const path = require('path');

const RESULTS_DIR = path.join(__dirname, '..', '..', '.marvin', 'results');

// Sample result data
const createSampleResult = (checker, timestamp, issueCount) => ({
  checker,
  timestamp,
  path: '/docs',
  summary: {
    total_files: 50,
    files_with_issues: Math.floor(issueCount / 3),
    total_issues: issueCount,
    error_count: Math.floor(issueCount * 0.2),
    warning_count: Math.floor(issueCount * 0.5),
    info_count: Math.floor(issueCount * 0.3),
  },
  issues: Array.from({ length: Math.min(issueCount, 10) }, (_, i) => ({
    file: `docs/file-${i + 1}.md`,
    line: Math.floor(Math.random() * 100) + 1,
    column: Math.floor(Math.random() * 80) + 1,
    severity: ['error', 'warning', 'info'][Math.floor(Math.random() * 3)],
    message: `Sample issue message ${i + 1}`,
    rule: `Rule.${checker}.${i + 1}`,
    context: 'This is sample context for the issue.',
  })),
  metadata: {},
});

// Create directory if it doesn't exist
if (!fs.existsSync(RESULTS_DIR)) {
  fs.mkdirSync(RESULTS_DIR, { recursive: true });
  console.log(`✓ Created directory: ${RESULTS_DIR}`);
}

// Generate sample files
const samples = [
  { checker: 'vale', date: new Date(Date.now() - 1000 * 60 * 5), issues: 15 },
  { checker: 'vale', date: new Date(Date.now() - 1000 * 60 * 60), issues: 23 },
  { checker: 'vale', date: new Date(Date.now() - 1000 * 60 * 60 * 24), issues: 8 },
  { checker: 'markdownlint', date: new Date(Date.now() - 1000 * 60 * 30), issues: 12 },
];

samples.forEach(({ checker, date, issues }) => {
  const timestamp = date.toISOString();
  const filename = `${checker}-${date.toISOString().replace(/[:.]/g, '-').slice(0, 19)}.json`;
  const filepath = path.join(RESULTS_DIR, filename);
  
  const data = createSampleResult(checker, timestamp, issues);
  
  fs.writeFileSync(filepath, JSON.stringify(data, null, 2));
  console.log(`✓ Created: ${filename}`);
});

console.log('\n✅ Test data created successfully!');
console.log(`📁 Location: ${RESULTS_DIR}`);
console.log('\nYou can now run: bun run dev');
