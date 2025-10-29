import { useState, useCallback } from 'react';
import { apiClient } from '@/lib/api';
import { AnalysisResponse } from '@/types/api';

export function useApi<T>() {
  const [data, setData] = useState<T | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const execute = useCallback(async (apiCall: () => Promise<T>) => {
    setLoading(true);
    setError(null);
    
    try {
      const result = await apiCall();
      setData(result);
      return result;
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Unknown error';
      setError(message);
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  const reset = useCallback(() => {
    setData(null);
    setError(null);
    setLoading(false);
  }, []);

  return {
    data,
    loading,
    error,
    execute,
    reset,
  };
}

export function useHealth() {
  const { data, loading, error, execute } = useApi();
  
  const getHealth = useCallback(() => 
    execute(() => apiClient.getHealth())
  , [execute]);

  return {
    health: data,
    loading,
    error,
    getHealth,
  };
}

export function useAnalysis() {
  const { data, loading, error, execute } = useApi<AnalysisResponse>();
  
  const analyze = useCallback((url: string) => 
    execute(() => apiClient.analyzeUrl(url))
  , [execute]);

  return {
    analysis: data,
    loading,
    error,
    analyze,
  };
}
