import { HealthResponse } from "@/types/api";

// For server-side rendering, ensure we strip any trailing /api/v1 from the base URL
// since we add it in each endpoint method
let API_BASE_URL =
  process.env.NEXT_PUBLIC_API_URL ||
  process.env.API_URL ||
  "http://localhost:8080";
// Remove trailing /api/v1 if present (common mistake)
API_BASE_URL = API_BASE_URL.replace(/\/api\/v1\/?$/, "");
// Remove trailing slash
API_BASE_URL = API_BASE_URL.replace(/\/$/, "");

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
          "Content-Type": "application/json",
          ...options?.headers,
        },
        // Ensure fresh data for SSR
        cache: "no-store",
        ...options,
      });

      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(`HTTP ${response.status}: ${response.statusText}`);
      }

      const data = await response.json();
      return data;
    } catch (error) {
      throw error;
    }
  }

  async getHealth(): Promise<HealthResponse> {
    return this.fetch<HealthResponse>("/api/v1/health");
  }

  async analyzeUrl(url: string): Promise<any> {
    // URL is already decoded by Next.js searchParams, so we encode it here for the API call
    const encodedUrl = encodeURIComponent(url);
    return this.fetch(`/api/v1/analyze?url=${encodedUrl}`);
  }

  async getStats(): Promise<any> {
    return this.fetch("/api/v1/stats");
  }
}

// Singleton instance for server-side usage
export const serverApi = ServerApiClient.getInstance();
