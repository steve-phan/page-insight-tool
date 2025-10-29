import { Analyze } from '@/types/Analyze';
import { Health } from '@/types/Health';
import { ApiConfig } from '@/types/http-client';
import type { ModelsAnalysisResponse, HealthHealthResponse, ModelsHTTPError } from '@/types/data-contracts';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

// Configure the API client
const apiConfig: ApiConfig = {
  baseUrl: `${API_BASE_URL}/api/v1`,
};

// Create API client instances
export const analyzeApi = new Analyze(apiConfig);
export const healthApi = new Health(apiConfig);

// Re-export types for convenience
export type { ModelsAnalysisResponse as AnalysisResponse };
export type { HealthHealthResponse as HealthResponse };
export type { ModelsHTTPError as HTTPError };

// Convenience wrapper functions
export const apiClient = {
  /**
   * Analyze a web page
   * @param url - URL of the web page to analyze
   */
  async analyzeUrl(url: string) {
    return analyzeApi.analyzeList({ url });
  },

  /**
   * Get health status
   */
  async getHealth() {
    return healthApi.healthList();
  },
};
