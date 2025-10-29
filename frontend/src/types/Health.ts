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

import { HealthHealthResponse } from "./data-contracts";
import { ContentType, HttpClient, RequestParams } from "./http-client";

export class Health<
  SecurityDataType = unknown,
> extends HttpClient<SecurityDataType> {
  /**
   * @description Returns the health status of the API including version, environment, and server status
   *
   * @tags Health
   * @name HealthList
   * @summary Health check endpoint
   * @request GET:/health
   */
  healthList = (params: RequestParams = {}) =>
    this.request<HealthHealthResponse, any>({
      path: `/health`,
      method: "GET",
      type: ContentType.Json,
      format: "json",
      ...params,
    });
}
