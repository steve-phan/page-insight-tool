import { HealthResponse, AnalysisResponse, AnalysisRequest, ApiError } from '@/types/api';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export class ApiClient {
  private async fetch<T>(endpoint: string, options?: RequestInit): Promise<T> {
    const url = `${API_BASE_URL}${endpoint}`;
    const response = await fetch(url, {
      headers: {
        'Content-Type': 'application/json',
        ...options?.headers,
      },
      ...options,
    });

    if (!response.ok) {
      const error: ApiError = await response.json().catch(() => ({
        message: 'Unknown error occurred',
        code: 'UNKNOWN_ERROR',
      }));
      throw new Error(error.message || `HTTP ${response.status}`);
    }

    return response.json();
  }

  // Health check
  async getHealth(): Promise<HealthResponse> {
    return this.fetch<HealthResponse>('/api/health');
  }

  // URL analysis
  async analyzeUrl(request: AnalysisRequest): Promise<AnalysisResponse> {
    return this.fetch<AnalysisResponse>('/api/v1/analyze', {
      method: 'POST',
      body: JSON.stringify(request),
    });
  }

  // Statistics
  async getStats(): Promise<any> {
    return this.fetch('/api/v1/stats');
  }
}

// Singleton instance
export const apiClient = new ApiClient();
