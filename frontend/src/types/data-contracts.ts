/* eslint-disable */
/* tslint:disable */
// @ts-nocheck
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

export interface HealthHealthResponse {
  /** @example "2025-10-29T22:00:00Z" */
  build_date?: string;
  /** @example "development" */
  environment?: string;
  /** @example "healthy" */
  status?: string;
  /** @example "2025-10-29T22:00:00Z" */
  timestamp?: string;
  /** @example "1.0.0" */
  version?: string;
}

export interface ModelsAnalysisResponse {
  /** @example 150 */
  analysis_time_ms?: number;
  /** @example true */
  has_login_form?: boolean;
  headings?: ModelsHeadings;
  /** @example "HTML5" */
  html_version?: string;
  links?: ModelsLinks;
  /** @example "Google" */
  page_title?: string;
}

export interface ModelsHTTPError {
  /** @example 400 */
  code?: number;
  details?: Record<string, any>;
  /** @example "Invalid URL" */
  message?: string;
  request_id?: string;
  timestamp?: string;
  /** @example "INVALID_URL" */
  type?: string;
}

export interface ModelsHeadings {
  /** @example 10 */
  h1?: number;
  /** @example 20 */
  h2?: number;
  /** @example 30 */
  h3?: number;
  /** @example 40 */
  h4?: number;
  /** @example 50 */
  h5?: number;
  /** @example 60 */
  h6?: number;
}

export interface ModelsLinks {
  /** @example 20 */
  external?: number;
  /** @example 30 */
  inaccessible?: number;
  /** @example 10 */
  internal?: number;
}
