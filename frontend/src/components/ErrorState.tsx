import { useToast } from '@/app/providers/ToastProvider';

export interface ErrorStateProps {
  title?: string;
  description?: string;
  actionLabel?: string;
  onRetry?: () => void;
}

export const ErrorState = ({
  title = 'Something went wrong',
  description = 'Please try again in a moment.',
  actionLabel = 'Retry',
  onRetry,
}: ErrorStateProps): JSX.Element => {
  const { toast } = useToast();

  const handleRetry = (): void => {
    if (onRetry) {
      onRetry();
      toast({ title: 'Retrying…', description: 'Request sent to the Hub API.' });
    }
  };

  return (
    <div className="flex min-h-[30vh] flex-col items-center justify-center gap-4 text-center">
      <div className="inline-flex items-center gap-3 rounded-full border border-danger/40 bg-danger/5 px-4 py-2 text-danger">
        <span className="text-xs font-semibold uppercase tracking-[0.28em]">Error</span>
      </div>
      <div className="space-y-2">
        <h3 className="text-2xl font-semibold text-text">{title}</h3>
        <p className="text-sm text-text-muted">{description}</p>
      </div>
      {onRetry ? (
        <button
          type="button"
          onClick={handleRetry}
          className="rounded-full border border-brand/60 bg-brand/10 px-4 py-2 text-sm font-medium text-brand transition hover:bg-brand/20"
        >
          {actionLabel}
        </button>
      ) : null}
    </div>
  );
};


