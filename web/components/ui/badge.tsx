import { cn } from '@/lib/utils';
import type { SeverityType } from '@/lib/types';

interface BadgeProps {
  severity: SeverityType;
  children: React.ReactNode;
  className?: string;
}

export default function Badge({ severity, children, className }: BadgeProps) {
  const severityStyles = {
    error: 'bg-error-500 text-white',
    warning: 'bg-warning-500 text-gray-900',
    info: 'bg-primary-500 text-white',
  };

  return (
    <span
      className={cn(
        'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium',
        severityStyles[severity],
        className
      )}
    >
      {children}
    </span>
  );
}
