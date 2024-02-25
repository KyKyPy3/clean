import { HttpStatusCode, type HttpClient } from "@/src/main/http"
import { UnexpectedErrorException } from "@/src/core/exceptions"

import type { UserRepository } from "@user/application/ports"
import type { UserEntity } from "@user/domain/entity/user"

import { type UserDTO, mapToDomain } from "./dto"

export class UserRepositoryImpl implements UserRepository {
  constructor(
    private readonly url: string,
    private readonly httpClient: HttpClient
  ) {}

  async me(): Promise<UserEntity> {
    const httpResponse = await this.httpClient.request<{ data: UserDTO }>({
      url: `${this.url}/me`,
      method: 'get',
    })

    switch (httpResponse.statusCode) {
      case HttpStatusCode.ok: {
        const res = httpResponse.body

        return mapToDomain(res!.data)
      // case HttpStatusCode.forbidden:
      //   throw new EmailInUsedException('provide email in registration already used')
      }
      default:
        throw new UnexpectedErrorException()
    }
  }

  async list(): Promise<UserEntity[]> {
    const httpResponse = await this.httpClient.request<{ data: { users: UserDTO[] } }>({
      url: this.url,
      method: 'get',
    })

    switch (httpResponse.statusCode) {
      case HttpStatusCode.ok: {
        const res = httpResponse.body

        return res!.data.users.map(u => mapToDomain(u))
      // case HttpStatusCode.forbidden:
      //   throw new EmailInUsedException('provide email in registration already used')
      }
      default:
        throw new UnexpectedErrorException()
    }
  }
}
