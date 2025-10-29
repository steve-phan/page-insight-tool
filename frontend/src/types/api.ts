export interface HealthResponse {
  status: string;
  timestamp: string;
  version: string;
  build_date: string;
  git_commit: string;
  uptime: string;
}

export interface AnalysisRequest {
  url: string;
}

export interface AnalysisResponse {
  html_version: string;
  page_title: string;
  headings: {
    h1: number;
    h2: number;
    h3: number;
    h4: number;
    h5: number;
    h6: number;
  };
  links: {
    internal: number;
    external: number;
    inaccessible: number;
  };
  has_login_form: boolean;
  analysis_time_ms: number;
}

export interface ApiError {
  message: string;
  code: string;
  details?: string;
}
