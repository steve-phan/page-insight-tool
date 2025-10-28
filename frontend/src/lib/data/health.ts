import { HealthResponse } from '@/types/api';
import { serverApi } from '@/lib/server-api';

export interface HealthDataResult {
  data: HealthResponse | null;
  error: string | null;
  timestamp: string;
}

/**
 * Fetches health data from the backend API
 * Used in server components for SSR
 */
export async function fetchHealthData(): Promise<HealthDataResult> {
  const timestamp = new Date().toISOString();
  
  try {
    const data = await serverApi.getHealth();
    return {
      data,
      error: null,
      timestamp,
    };
  } catch (error) {
    console.error('Failed to fetch health data:', error);
    
    // Return fallback data instead of throwing
    return {
      data: {
        status: 'unhealthy',
        timestamp: new Date().toISOString(),
        version: 'unknown',
        build_date: 'unknown',
        git_commit: 'unknown',
        uptime: 'unknown',
      },
      error: error instanceof Error ? error.message : 'Unknown error',
      timestamp,
    };
  }
}

/**
 * Validates health data structure
 */
export function validateHealthData(data: any): data is HealthResponse {
  return (
    data &&
    typeof data.status === 'string' &&
    typeof data.timestamp === 'string' &&
    typeof data.version === 'string' &&
    typeof data.build_date === 'string' &&
    typeof data.git_commit === 'string' &&
    typeof data.uptime === 'string'
  );
}
