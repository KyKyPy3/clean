import { AxiosHttpClient } from "@/src/lib/http/http"

import type { HttpClient } from "./http-client.port"

export const makeHttpClient = (): HttpClient => {
  return new AxiosHttpClient()
}
