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

import { ModelsAnalysisResponse, ModelsHTTPError } from "./data-contracts";
import { ContentType, HttpClient, RequestParams } from "./http-client";

export class Analyze<
  SecurityDataType = unknown,
> extends HttpClient<SecurityDataType> {
  /**
   * @description Analyzes a web page and extracts HTML version, title, headings, links, login forms, and CSR detection information
   *
   * @tags Analysis
   * @name AnalyzeList
   * @summary Analyze a web page
   * @request GET:/analyze
   */
  analyzeList = (
    query: {
      /**
       * URL of the web page to analyze
       * @example "https://example.com"
       */
      url: string;
    },
    params: RequestParams = {},
  ) =>
    this.request<ModelsAnalysisResponse, ModelsHTTPError>({
      path: `/analyze`,
      method: "GET",
      query: query,
      type: ContentType.Json,
      format: "json",
      ...params,
    });
}
