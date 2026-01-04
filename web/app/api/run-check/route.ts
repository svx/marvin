import { NextRequest, NextResponse } from 'next/server';
import { exec } from 'child_process';
import { promisify } from 'util';
import path from 'path';
import fs from 'fs';

const execAsync = promisify(exec);

export async function POST(request: NextRequest) {
  try {
    const { checker, path: checkPath } = await request.json();

    if (!checker || !['vale', 'markdownlint'].includes(checker)) {
      return NextResponse.json(
        { error: 'Invalid checker type. Must be "vale" or "markdownlint"' },
        { status: 400 }
      );
    }

    // Get the project root (assuming web is in the project root)
    const projectRoot = path.resolve(process.cwd(), '..');
    const cliPath = path.join(projectRoot, 'cli');
    const outputDir = path.join(projectRoot, '.marvin', 'results');

    // Use provided path or default to ../docs
    const targetPath = checkPath || '../docs';

    // Note: The CLI auto-detects config files (.markdownlint.yaml and .vale.ini)
    // in the project root, so we don't need to explicitly pass --config flag

    // Log for debugging
    console.log('Running check:', { checker, checkPath, targetPath });

    // Run the CLI command with positional argument, output-dir, and --no-tui flag
    // The CLI will auto-detect .markdownlint.yaml and .vale.ini config files
    const command = `cd ${cliPath} && go run main.go ${checker} ${targetPath} --output-dir "${outputDir}" --no-tui`;
    
    try {
      const { stdout, stderr } = await execAsync(command, {
        timeout: 60000, // 60 second timeout
      });

      return NextResponse.json({
        success: true,
        message: `${checker} check completed successfully`,
        output: stdout,
      });
    } catch (execError: any) {
      // The CLI returns exit code 1 when issues are found, but this is expected behavior
      // Check if we got output (which means the check ran successfully)
      if (execError.stdout && execError.stdout.length > 0) {
        return NextResponse.json({
          success: true,
          message: `${checker} check completed successfully`,
          output: execError.stdout,
        });
      }
      
      // If there's no stdout, it's a real error
      throw execError;
    }
  } catch (error) {
    console.error('Error running check:', error);
    
    const errorMessage = error instanceof Error ? error.message : 'Unknown error';
    
    return NextResponse.json(
      {
        error: 'Failed to run check',
        details: errorMessage,
      },
      { status: 500 }
    );
  }
}
