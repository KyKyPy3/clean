import type { GameEntity } from "@game/domain/entity/game"

export interface GameRepository {
  list(): Promise<GameEntity[]>
}