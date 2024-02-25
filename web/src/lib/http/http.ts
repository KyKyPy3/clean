import axios from 'axios';
import { AxiosError, AxiosInstance, AxiosResponse } from 'axios';

import type { HttpClient, HttpRequest, HttpResponse } from '@/src/main/http';

export class AxiosHttpClient implements HttpClient {
  private _instance: AxiosInstance

  constructor() {
    this._instance = axios.create({
      responseType: 'json',
      timeout: 10000,
      withCredentials: true,
      headers: {
        'X-Requested-With': 'XMLHttpRequest',
        'Content-Type': 'application/json'
      }
    });

    this._instance.interceptors.response.use(
      (response) => response,
      async (error) => {
        const originalRequest = error.config;

        if (error.response.status === 401 && !originalRequest._retry) {
          originalRequest._retry = true;

          try {
            // Try to refresh token
            await axios.post('/api/v1/auth/refresh');

            return axios(originalRequest);
          } catch (e) {
            const axiosError = e as AxiosError;
            if (axiosError.status === 404) {
              window.location.href = '/login'
            }
            return Promise.reject(error);
          }
        }

        return Promise.reject(error);
      }
    )
  }

  public async request(data: HttpRequest): Promise<HttpResponse> {
    let axiosResponse: AxiosResponse

    try {
      axiosResponse = await this._instance.request({
        url: data.url,
        method: data.method,
        data: data.body,
        headers: data.headers,
      })
    } catch (error) {
      axiosResponse = (error as AxiosError).response as AxiosResponse
    }

    return this.adapt(axiosResponse);
  }

  private adapt(axiosResponse: AxiosResponse): HttpResponse {
    return {
      body: axiosResponse.data,
      statusCode: axiosResponse.status,
    };
  }
}