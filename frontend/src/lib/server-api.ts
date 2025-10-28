import { HealthResponse } from '@/types/api';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export class ServerApiClient {
  private static instance: ServerApiClient;
  
  private constructor() {}
  
  public static getInstance(): ServerApiClient {
    if (!ServerApiClient.instance) {
      ServerApiClient.instance = new ServerApiClient();
    }
    return ServerApiClient.instance;
  }

  private async fetch<T>(endpoint: string, options?: RequestInit): Promise<T> {
    const url = `${API_BASE_URL}${endpoint}`;
    
    try {
      const response = await fetch(url, {
        headers: {
          'Content-Type': 'application/json',
          ...options?.headers,
        },
        // Ensure fresh data for SSR
        cache: 'no-store',
        ...options,
      });

      if (!response.ok) {
        throw new Error(`HTTP ${response.status}: ${response.statusText}`);
      }

      return await response.json();
    } catch (error) {
      console.error(`API Error [${endpoint}]:`, error);
      throw error;
    }
  }

  async getHealth(): Promise<HealthResponse> {
    return this.fetch<HealthResponse>('/health');
  }

  async analyzeUrl(url: string): Promise<any> {
    return this.fetch('/api/v1/analyze', {
      method: 'POST',
      body: JSON.stringify({ url }),
    });
  }

  async getStats(): Promise<any> {
    return this.fetch('/api/v1/stats');
  }
}

// Singleton instance for server-side usage
export const serverApi = ServerApiClient.getInstance();
