import { HttpStatusCode, type HttpClient } from "@/src/main/http"
import { UnexpectedErrorException } from "@/src/core/exceptions"

import type { GameRepository } from "@game/application/ports"
import type { GameEntity } from "@game/domain/entity/game"

import { type GameDTO, mapToDomain } from "./dto"

export class GameRepositoryImpl implements GameRepository {
  constructor(
    private readonly url: string,
    private readonly httpClient: HttpClient
  ) {}

  async list(): Promise<GameEntity[]> {
    const httpResponse = await this.httpClient.request<{ data: { games: GameDTO[] } }>({
      url: this.url,
      method: 'get',
    })

    switch (httpResponse.statusCode) {
      case HttpStatusCode.ok: {
        const res = httpResponse.body

        return res!.data.games.map(u => mapToDomain(u))
      // case HttpStatusCode.forbidden:
      //   throw new EmailInUsedException('provide email in registration already used')
      }
      default:
        throw new UnexpectedErrorException()
    }
  }
}
