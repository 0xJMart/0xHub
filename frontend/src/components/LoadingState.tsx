export interface LoadingStateProps {
  message?: string;
}

export const LoadingState = ({ message = 'Loading…' }: LoadingStateProps): JSX.Element => (
  <div className="flex min-h-[40vh] flex-col items-center justify-center gap-4 text-center text-text-muted">
    <div className="h-12 w-12 animate-spin rounded-full border-4 border-border/60 border-t-brand" />
    <p className="text-sm font-medium uppercase tracking-[0.3em] text-text-muted/80">{message}</p>
  </div>
);


