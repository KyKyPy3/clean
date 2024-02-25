import { HttpStatusCode, type HttpClient } from "@/src/main/http";

import { UnexpectedErrorException, type Email } from "@core";

import type { SessionRepository } from "../../application/session.port";


export class SessionRepositoryImpl implements SessionRepository {
  constructor(
    private readonly url: string,
    private readonly httpClient: HttpClient,
  ) {}

  async login(email: Email, password: string): Promise<{ email: Email, token: string }> {
    const httpResponse = await this.httpClient.request<{ data: { access_token: string } }>({
      url: this.url,
      method: 'post',
      body: {
        email: email.value,
        password: password,
      },
    })

    switch (httpResponse.statusCode) {
      case HttpStatusCode.ok:
        return { email, token: httpResponse.body?.data.access_token || '' }
      // case HttpStatusCode.forbidden:
      //   throw new EmailInUsedException('provide email in registration already used')
      default:
        throw new UnexpectedErrorException()
    }
  }
}
