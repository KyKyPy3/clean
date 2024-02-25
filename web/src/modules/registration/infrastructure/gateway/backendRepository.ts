import { HttpStatusCode, type HttpClient } from "@/src/main/http"
import { UnexpectedErrorException } from "@/src/core/exceptions"

import type { RegistrationEntity } from "@registration/domain/entity/registration"
import { EmailInUsedException } from "@registration/domain/exceptions"
import type { RegistrationRepository } from "@registration/application/ports"

import { mapToRequest } from "./dto"

export class RegistrationRepositoryImpl implements RegistrationRepository {
  constructor(
    private readonly url: string,
    private readonly httpClient: HttpClient,
  ) {}

  async register(registration: RegistrationEntity): Promise<void> {
    const request = mapToRequest(registration)

    const httpResponse = await this.httpClient.request({
      url: this.url,
      method: 'post',
      body: request,
    })

    switch (httpResponse.statusCode) {
      case HttpStatusCode.ok:
        return
      case HttpStatusCode.forbidden:
        throw new EmailInUsedException('provide email in registration already used')
      default:
        throw new UnexpectedErrorException()
    }
  }
}