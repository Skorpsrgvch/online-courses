import { useState, useCallback } from 'react';

interface UseApiState<T> {
  data: T | null;
  loading: boolean;
  error: string | null;
}

interface UseApiReturn<T> extends UseApiState<T> {
  execute: (...args: any[]) => Promise<T | null>;
  reset: () => void;
}

export function useApi<T>(asyncFunction: (...args: any[]) => Promise<T>): UseApiReturn<T> {
  const [state, setState] = useState<UseApiState<T>>({
    data: null,
    loading: false,
    error: null,
  });

  const execute = useCallback(async (...args: any[]): Promise<T | null> => {
    setState((prev) => ({ ...prev, loading: true, error: null }));
    try {
      const result = await asyncFunction(...args);
      setState({ data: result, loading: false, error: null });
      return result;
    } catch (err: any) {
      const errorMessage = err.message || 'Произошла ошибка при загрузке данных';
      setState((prev) => ({ ...prev, loading: false, error: errorMessage }));
      return null;
    }
  }, [asyncFunction]);

  const reset = useCallback(() => {
    setState({ data: null, loading: false, error: null });
  }, []);

  return { ...state, execute, reset };
}