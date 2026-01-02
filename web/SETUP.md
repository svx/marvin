# Marvin Web Dashboard - Setup Guide

This guide will help you set up and run the Marvin web dashboard.

## Prerequisites

- **Bun** (latest version) - Already available in your Devbox environment
- **Node.js** 20+ (if not using Bun)

## Installation

### 1. Navigate to the web directory

```bash
cd web
```

### 2. Install dependencies

Using Bun (recommended):
```bash
bun install
```

Or using npm:
```bash
npm install
```

### 3. Set up environment variables (optional)

```bash
cp .env.local.example .env.local
```

Edit `.env.local` if you need to customize the results directory path.

## Development

### Start the development server

Using Bun:
```bash
bun run dev
```

Or using npm:
```bash
npm run dev
```

The dashboard will be available at [http://localhost:3000](http://localhost:3000)

## Project Structure

```
web/
├── app/                      # Next.js App Router
│   ├── layout.tsx           # Root layout
│   ├── page.tsx             # Dashboard home
│   ├── globals.css          # Global styles
│   └── api/                 # API routes
│       └── results/         # Results API endpoints
├── components/              # React components
│   ├── ui/                 # Reusable UI components
│   └── layout/             # Layout components
├── lib/                    # Utilities and types
│   ├── types.ts           # TypeScript types
│   ├── utils.ts           # Utility functions
│   ├── api.ts             # API client
│   └── constants.ts       # Constants
├── public/                # Static assets
├── package.json           # Dependencies
├── tsconfig.json          # TypeScript config
├── tailwind.config.ts     # Tailwind config
└── next.config.js         # Next.js config
```

## Features Implemented

### ✅ Phase 1 - Foundation (MVP)
- [x] Next.js project with TypeScript and Tailwind
- [x] Basic layout (header, main content)
- [x] API routes to read result files
- [x] TypeScript types matching CLI models
- [x] Dashboard overview page with summary stats
- [x] Stripe-inspired design system

### 🚧 Phase 2 - Core Features (To be implemented)
- [ ] All checks view with filtering
- [ ] Individual check result pages
- [ ] Issue list and detail components
- [ ] Severity badges and status indicators
- [ ] Search and filter functionality

### 📋 Phase 3 - Enhanced UX (Future)
- [ ] File tree navigation
- [ ] Sorting and pagination
- [ ] Export functionality (CSV, JSON)
- [ ] Responsive design for mobile
- [ ] Dark mode toggle

### 🚀 Phase 4 - Advanced Features (Future)
- [ ] Trigger checks from web interface
- [ ] Real-time updates via WebSocket
- [ ] Issue trend charts and analytics
- [ ] Rule documentation integration
- [ ] User preferences and settings

## How It Works

### Data Flow

1. **CLI generates results**: The Marvin CLI runs checks and saves JSON files to `.marvin/results/`
2. **API reads files**: Next.js API routes read these JSON files from the filesystem
3. **Dashboard displays data**: React components fetch data from the API and display it

### API Endpoints

- `GET /api/results` - Fetch all check results
  - Query params: `checker`, `limit`, `offset`
- `GET /api/results/[id]` - Fetch a specific result by ID

### File Structure

Results are stored in `.marvin/results/` with the naming pattern:
```
{checker}-{timestamp}.json
```

Example:
```
vale-2026-01-02-123456.json
```

## Development Tips

### TypeScript Errors

The project uses strict TypeScript. If you see errors about missing modules after installation, run:

```bash
bun install
```

This will install all dependencies including type definitions.

### Tailwind CSS

The project uses a custom Tailwind configuration with Stripe-inspired colors:

- **Primary**: Indigo/Purple (`#635BFF`)
- **Success**: Green (`#00D924`)
- **Warning**: Amber (`#FFA500`)
- **Error**: Red (`#DF1B41`)

### Adding New Components

1. Create component in `components/ui/` or `components/dashboard/`
2. Use the `cn()` utility for conditional classes
3. Follow the existing component patterns

Example:
```tsx
import { cn } from '@/lib/utils';

export default function MyComponent({ className }: { className?: string }) {
  return (
    <div className={cn('base-classes', className)}>
      Content
    </div>
  );
}
```

## Testing the Dashboard

### 1. Generate test data

First, run the CLI to generate some results:

```bash
cd ../cli
go run main.go vale ../docs
```

This will create result files in `.marvin/results/`

### 2. Start the dashboard

```bash
cd ../web
bun run dev
```

### 3. View results

Open [http://localhost:3000](http://localhost:3000) to see your results displayed in the dashboard.

## Building for Production

```bash
bun run build
```

This creates an optimized production build in `.next/`

To start the production server:

```bash
bun run start
```

## Troubleshooting

### "Cannot find module" errors

Make sure all dependencies are installed:
```bash
bun install
```

### API returns empty results

Check that:
1. The `.marvin/results/` directory exists
2. There are JSON files in the directory
3. The files are valid JSON

### Styles not loading

Make sure Tailwind is properly configured:
```bash
# Check that tailwind.config.ts exists
# Check that globals.css imports Tailwind directives
```

## Next Steps

To continue development:

1. **Implement the checks page** (`app/checks/page.tsx`)
2. **Create individual checker pages** (`app/checks/[checker]/page.tsx`)
3. **Add filtering and search** components
4. **Implement issue detail views**

See the full plan in [`README.md`](./README.md) for detailed specifications.

## Contributing

When adding new features:

1. Follow the existing code structure
2. Use TypeScript for type safety
3. Follow the Stripe-inspired design patterns
4. Test with actual CLI output
5. Update this guide if needed

## Resources

- [Next.js Documentation](https://nextjs.org/docs)
- [Tailwind CSS](https://tailwindcss.com/docs)
- [TypeScript](https://www.typescriptlang.org/docs)
- [Stripe Design System](https://stripe.com/docs) (for inspiration)
