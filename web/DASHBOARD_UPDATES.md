# Dashboard Updates Summary

This document summarizes the updates made to the Next.js dashboard.

## Changes Implemented

### 1. Favicon Integration ✅
- **File**: `web/public/favicon.svg`
  - Copied from `docs/public/favicon.svg`
  - Blue (#5469d4) icon with white document/dashboard symbols

- **File**: `web/app/layout.tsx`
  - Added favicon metadata to the layout
  - Icon will now appear in browser tabs

### 2. Logo in Header ✅
- **File**: `web/components/layout/header.tsx`
  - Added logo image using the favicon
  - Logo appears next to "Marvin" text
  - Includes hover animation (scale effect)
  - Uses Next.js Image component for optimization

### 3. GitHub Repository Link ✅
- **File**: `web/components/layout/header.tsx`
  - Added GitHub icon and link in the header
  - Links to: `https://github.com/svx/marvin`
  - Opens in new tab with proper security attributes
  - Responsive: shows icon + text on larger screens, icon only on mobile

### 4. Re-run Check Functionality ✅

#### API Endpoint
- **File**: `web/app/api/run-check/route.ts`
  - POST endpoint to trigger check re-runs
  - Accepts `checker` parameter ('vale' or 'markdownlint')
  - Executes CLI commands via Node.js child_process
  - Returns success/error responses with proper error handling
  - 60-second timeout for long-running checks

#### API Helper Function
- **File**: `web/lib/api.ts`
  - Added `runCheck()` function
  - Handles API calls to the run-check endpoint
  - Proper error handling and type safety

#### Reusable Button Component
- **File**: `web/components/ui/run-check-button.tsx`
  - Reusable component for running checks
  - Features:
    - Loading state with spinner animation
    - Success/error feedback messages
    - Auto-dismiss messages (3s for success, 5s for errors)
    - Configurable variants (primary/secondary)
    - Configurable sizes (sm/md/lg)
    - Optional onSuccess callback for refreshing data
  - Accessible with proper ARIA attributes

#### Integration in Pages

**Dashboard Page** (`web/app/page.tsx`):
- Added "Run Quality Checks" card at the top
- Both Vale and Markdownlint buttons available
- Auto-refreshes results after check completion
- Helpful description text for users

**All Checks Page** (`web/app/checks/page.tsx`):
- Added "Run Quality Checks" card
- Same functionality as dashboard
- Results refresh automatically after completion

**Checker Detail Page** (`web/app/checks/[checker]/page.tsx`):
- Added re-run button in the header (top-right)
- Uses secondary variant for subtle appearance
- Redirects to latest result after completion
- Smaller size (sm) to fit in header

## Technical Details

### Dependencies
No new dependencies were added. The implementation uses:
- Next.js built-in features (Image, routing, API routes)
- Node.js built-in modules (child_process, util, path)
- Existing UI components and utilities

### File Structure
```
web/
├── public/
│   └── favicon.svg (new)
├── app/
│   ├── layout.tsx (modified)
│   ├── page.tsx (modified)
│   ├── checks/
│   │   ├── page.tsx (modified)
│   │   └── [checker]/
│   │       └── page.tsx (modified)
│   └── api/
│       └── run-check/
│           └── route.ts (new)
├── components/
│   ├── layout/
│   │   └── header.tsx (modified)
│   └── ui/
│       └── run-check-button.tsx (new)
└── lib/
    └── api.ts (modified)
```

## Testing Recommendations

1. **Favicon**: Check browser tab to see the icon
2. **Logo**: Verify logo appears in header with hover effect
3. **GitHub Link**: Click to ensure it opens the correct repository
4. **Re-run Checks**:
   - Click "Re-run vale" button on dashboard
   - Verify loading state appears
   - Check that results refresh after completion
   - Test error handling by running without CLI available
   - Verify success/error messages display correctly

## Usage

### Running the Dashboard
```bash
cd web
bun install
bun run dev
```

### Running Checks from Dashboard
1. Navigate to the dashboard (http://localhost:3000)
2. Click "Re-run vale" or "Re-run markdownlint"
3. Wait for the check to complete (loading spinner will show)
4. Results will automatically refresh and appear in the list

## Notes

- The API endpoint executes CLI commands, so the Go CLI must be built and available
- Check execution happens server-side for security
- Results are stored in `.marvin/results/` directory
- The dashboard reads from the same directory the CLI writes to
