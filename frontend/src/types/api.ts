// Re-export generated types with cleaner names
import type {
  HealthHealthResponse,
  ModelsAnalysisResponse,
  ModelsHeadings,
  ModelsLinks,
  ModelsHTTPError,
} from "./data-contracts";

export type {
  HealthHealthResponse as HealthResponse,
  ModelsAnalysisResponse as AnalysisResponse,
  ModelsHeadings as Headings,
  ModelsLinks as Links,
  ModelsHTTPError as HTTPError,
};

// For backward compatibility, export the same types
export type AnalysisRequest = {
  url: string;
};

export type ApiError = ModelsHTTPError;
