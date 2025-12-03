import {
  createContext,
  type PropsWithChildren,
  type ReactNode,
  useCallback,
  useContext,
  useEffect,
  useId,
  useMemo,
  useRef,
  useState,
} from 'react';
import { createPortal } from 'react-dom';

export type ToastVariant = 'default' | 'success' | 'error';

export interface ToastOptions {
  title: string;
  description?: string;
  variant?: ToastVariant;
  action?: ReactNode;
  duration?: number;
}

export interface ToastRecord extends ToastOptions {
  id: string;
}

interface ToastContextValue {
  toast: (options: ToastOptions) => void;
  dismiss: (id: string) => void;
  clear: () => void;
}

const ToastContext = createContext<ToastContextValue | null>(null);

const VARIANT_STYLES: Record<ToastVariant, string> = {
  default: 'border-border/70 bg-surface-raised text-text',
  success: 'border-state-success/70 bg-state-success/10 text-state-success',
  error: 'border-state-danger/70 bg-state-danger/10 text-state-danger',
};

export const ToastProvider = ({ children }: PropsWithChildren): JSX.Element => {
  const defaultDuration = 4000;
  const counterRef = useRef(0);
  const [toasts, setToasts] = useState<ToastRecord[]>([]);
  const portalId = useId();

  const ensurePortal = useCallback(() => {
    const portalIdAttribute = `toast-portal-${portalId}`;
    let portal = document.getElementById(portalIdAttribute);
    if (!portal) {
      portal = document.createElement('div');
      portal.id = portalIdAttribute;
      document.body.appendChild(portal);
    }
    return portal;
  }, [portalId]);

  const toast = useCallback(
    (options: ToastOptions) => {
      counterRef.current += 1;
      const id = `toast-${counterRef.current}`;
      const record: ToastRecord = {
        variant: 'default',
        duration: defaultDuration,
        ...options,
        id,
      };
      setToasts((prev) => [...prev, record]);
    },
    [defaultDuration],
  );

  const dismiss = useCallback((id: string) => {
    setToasts((prev) => prev.filter((toastItem) => toastItem.id !== id));
  }, []);

  const clear = useCallback(() => {
    setToasts([]);
  }, []);

  useEffect(() => {
    if (!toasts.length) {
      return;
    }
    const timers = toasts.map((toastItem) => {
      const timeout = window.setTimeout(
        () => dismiss(toastItem.id),
        toastItem.duration ?? defaultDuration,
      );
      return timeout;
    });
    return () => {
      timers.forEach((timeout) => window.clearTimeout(timeout));
    };
  }, [defaultDuration, dismiss, toasts]);

  const value = useMemo<ToastContextValue>(
    () => ({
      toast,
      dismiss,
      clear,
    }),
    [toast, dismiss, clear],
  );

  return (
    <ToastContext.Provider value={value}>
      {children}
      {createPortal(
        <div className="pointer-events-none fixed inset-x-0 top-4 z-50 flex flex-col items-center gap-3 px-4">
          {toasts.map((toastItem) => (
            <ToastItem key={toastItem.id} toast={toastItem} onDismiss={dismiss} />
          ))}
        </div>,
        ensurePortal(),
      )}
    </ToastContext.Provider>
  );
};

// eslint-disable-next-line react-refresh/only-export-components
export const useToast = (): ToastContextValue => {
  const context = useContext(ToastContext);
  if (!context) {
    throw new Error('useToast must be used within a ToastProvider');
  }
  return context;
};

interface ToastItemProps {
  toast: ToastRecord;
  onDismiss: (id: string) => void;
}

const ToastItem = ({ toast, onDismiss }: ToastItemProps): JSX.Element => {
  const { id, title, description, action, variant = 'default' } = toast;
  return (
    <div
      className={`pointer-events-auto flex w-full max-w-md items-start gap-3 rounded-2xl border px-5 py-4 shadow-surface-md transition ${VARIANT_STYLES[variant]}`}
    >
      <div className="flex-1 space-y-1">
        <p className="text-sm font-semibold tracking-tight">{title}</p>
        {description ? <p className="text-sm text-text-muted">{description}</p> : null}
        {action ? <div className="pt-2">{action}</div> : null}
      </div>
      <button
        type="button"
        onClick={() => onDismiss(id)}
        className="rounded-full border border-transparent p-1 text-xs uppercase tracking-[0.2em] text-text-muted transition hover:border-border/70 hover:text-text"
      >
        Close
      </button>
    </div>
  );
};


