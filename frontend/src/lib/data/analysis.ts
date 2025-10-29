import { AnalysisResponse } from "@/types/api";
import { serverApi } from "@/lib/server-api";

export interface AnalysisDataResult {
  data: AnalysisResponse | null;
  error: string | null;
  timestamp: string;
}

/**
 * Fetches analysis data from the backend API
 * Used in server components for SSR
 */
export async function fetchAnalysisData(
  url: string
): Promise<AnalysisDataResult> {
  const timestamp = new Date().toISOString();

  try {
    const data = await serverApi.analyzeUrl(url);
    return {
      data,
      error: null,
      timestamp,
    };
  } catch (error) {
    return {
      data: null,
      error: error instanceof Error ? error.message : "Unknown error",
      timestamp,
    };
  }
}

/**
 * Validates analysis data structure
 */
export function validateAnalysisData(data: any): data is AnalysisResponse {
  return (
    data &&
    typeof data.html_version === "string" &&
    typeof data.page_title === "string" &&
    typeof data.headings === "object" &&
    typeof data.links === "object" &&
    typeof data.has_login_form === "boolean" &&
    typeof data.analysis_time_ms === "number"
  );
}
