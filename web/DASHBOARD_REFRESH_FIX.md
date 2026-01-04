# Dashboard Refresh Fix

## Problem
When running a check via the re-run button, the dashboard showed a success message but didn't display the latest run at the top of the list. Users only saw older results from minutes ago.

## Root Cause
1. **Path Mismatch (PRIMARY ISSUE)**: The CLI was writing results to `cli/.marvin/results/` but the web dashboard was reading from the root `.marvin/results/` directory
2. **Timing Issue**: The `onSuccess` callback was triggered with only a 1-second delay, but the file system might not have finished writing the new result file yet
3. **Caching Issue**: Browser was caching API responses, preventing fresh data from being fetched
4. **No Visual Feedback**: Users couldn't tell when the dashboard was refreshing after a check completed

## Solution

### 1. Fixed Output Directory Path (CRITICAL FIX)
**File**: `web/app/api/run-check/route.ts`
- Added explicit `--output-dir` flag when running CLI from web dashboard
- Ensures CLI writes to the root `.marvin/results/` directory
- Command now includes: `--output-dir "${projectRoot}/.marvin/results"`
- This ensures the CLI writes to the same location the web dashboard reads from

### 2. Increased Delay Before Refresh
**File**: `web/components/ui/run-check-button.tsx`
- Increased delay from 1 second to 2 seconds before calling `onSuccess` callback
- This ensures the CLI has finished writing the result file to disk
- Also increased success message display time from 3 to 5 seconds for better UX

### 3. Added Cache-Busting
**File**: `web/lib/api.ts`
- Added timestamp query parameter (`?_t=${Date.now()}`) to all fetch requests
- Added `cache: 'no-store'` option to fetch calls
- Applied to both `fetchResults()` and `fetchResultById()` functions
- Ensures browser always fetches fresh data from the API

### 4. Added Refreshing Indicators
**Files**:
- `web/app/page.tsx` (Dashboard)
- `web/app/checks/page.tsx` (All Checks)

Added visual feedback when results are being refreshed:
- Blue notification banner with spinning icon
- "Refreshing results..." message
- Appears while fetching new data after a check completes

### 5. Improved Detail Page Navigation
**File**: `web/app/checks/[checker]/page.tsx`
- Removed `window.location.reload()` which caused a jarring full page reload
- Now smoothly navigates to the latest result if it's different
- Updates data in place if viewing the same result
- Better user experience with no page flicker

## Technical Details

### API Response Sorting
The API already sorts results by timestamp (newest first) in `web/app/api/results/route.ts`:
```typescript
results.sort((a, b) => 
  new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime()
);
```

### Refresh Flow
1. User clicks "Re-run" button
2. Check runs via API (`/api/run-check`)
3. API explicitly sets output directory to root `.marvin/results/`
4. CLI writes result file to the correct location
5. After 2-second delay, `onSuccess` callback fires
6. Dashboard fetches fresh results with cache-busting
7. Results are sorted by timestamp (newest first)
8. UI updates to show latest run at the top

### Why the Path Fix Was Critical
- **Before**: CLI ran from `cli/` directory, so `.marvin/results` resolved to `cli/.marvin/results/`
- **After**: Explicit `--output-dir` flag points to root `.marvin/results/`
- **Result**: CLI writes to same directory that web dashboard reads from

## Testing
To verify the fix:
1. Navigate to the dashboard or checks page
2. Click "Re-run vale" or "Re-run markdownlint"
3. Wait for the success message
4. Observe the "Refreshing results..." indicator
5. Verify the latest run appears at the top of the list
6. Check that the timestamp shows "just now" or "a few seconds ago"

## Files Modified
- `web/app/api/run-check/route.ts` - **CRITICAL**: Added explicit output directory path
- `web/components/ui/run-check-button.tsx` - Increased delays
- `web/lib/api.ts` - Added cache-busting
- `web/app/page.tsx` - Added refresh indicator
- `web/app/checks/page.tsx` - Added refresh indicator
- `web/app/checks/[checker]/page.tsx` - Improved navigation
