import { Component, type ErrorInfo, type ReactNode } from 'react';

interface GlobalErrorBoundaryProps {
  children: ReactNode;
  fallback?: ReactNode;
}

interface GlobalErrorBoundaryState {
  hasError: boolean;
  error?: Error;
}

export class GlobalErrorBoundary extends Component<
  GlobalErrorBoundaryProps,
  GlobalErrorBoundaryState
> {
  constructor(props: GlobalErrorBoundaryProps) {
    super(props);
    this.state = { hasError: false };
  }

  static getDerivedStateFromError(error: Error): GlobalErrorBoundaryState {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo): void {
    console.error('Uncaught error:', error, errorInfo);
  }

  handleReset = (): void => {
    this.setState({ hasError: false, error: undefined });
  };

  render(): ReactNode {
    if (this.state.hasError) {
      if (this.props.fallback) {
        return this.props.fallback;
      }
      return (
        <div className="flex min-h-screen flex-col items-center justify-center gap-4 bg-surface text-center text-text">
          <div className="rounded-full border border-danger/40 bg-danger/10 px-4 py-2 text-xs font-semibold uppercase tracking-[0.28em] text-danger">
            Unhandled Error
          </div>
          <div className="space-y-2 px-6">
            <h1 className="text-3xl font-semibold">We hit an unexpected issue</h1>
            <p className="text-sm text-text-muted">
              {this.state.error?.message ?? 'Please refresh the page and try again.'}
            </p>
          </div>
          <div className="flex items-center gap-3">
            <button
              type="button"
              className="rounded-full border border-brand/70 bg-brand/10 px-4 py-2 text-sm font-medium text-brand transition hover:bg-brand/20"
              onClick={() => window.location.reload()}
            >
              Refresh
            </button>
            <button
              type="button"
              className="rounded-full border border-border/60 px-4 py-2 text-sm font-medium text-text-muted transition hover:text-text"
              onClick={this.handleReset}
            >
              Dismiss
            </button>
          </div>
        </div>
      );
    }

    return this.props.children;
  }
}


