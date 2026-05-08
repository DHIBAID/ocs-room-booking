import { createContext, useCallback, useContext, useMemo, useState } from 'react'

type ToastTone = 'success' | 'error' | 'info'

type Toast = {
  id: string
  tone: ToastTone
  message: string
}

type ToastContextValue = {
  pushToast: (tone: ToastTone, message: string) => void
}

const ToastContext = createContext<ToastContextValue | undefined>(undefined)

const createId = () => `${Date.now()}-${Math.random().toString(16).slice(2)}`

export const ToastProvider = ({ children }: { children: React.ReactNode }) => {
  const [toasts, setToasts] = useState<Toast[]>([])

  const removeToast = useCallback((id: string) => {
    setToasts((prev) => prev.filter((toast) => toast.id !== id))
  }, [])

  const pushToast = useCallback((tone: ToastTone, message: string) => {
    const id = createId()
    setToasts((prev) => [...prev, { id, tone, message }])
    window.setTimeout(() => removeToast(id), 4200)
  }, [removeToast])

  const value = useMemo(() => ({ pushToast }), [pushToast])

  return (
    <ToastContext.Provider value={value}>
      {children}
      <div className="toast-stack" role="status" aria-live="polite">
        {toasts.map((toast) => (
          <div key={toast.id} className={`toast toast--${toast.tone}`}>
            <span>{toast.message}</span>
            <button
              type="button"
              className="toast-dismiss"
              onClick={() => removeToast(toast.id)}
              aria-label="Dismiss notification"
            >
              ×
            </button>
          </div>
        ))}
      </div>
    </ToastContext.Provider>
  )
}

export const useToast = () => {
  const ctx = useContext(ToastContext)
  if (!ctx) {
    throw new Error('useToast must be used within ToastProvider')
  }
  return ctx
}
