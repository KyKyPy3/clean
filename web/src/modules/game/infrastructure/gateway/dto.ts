import { GameEntity } from "../../domain/entity/game";

export type GameDTO = {
  id: string;
  name: string;
}

export function mapToDomain(game: GameDTO): GameEntity {
  const entity = new GameEntity({
    id: game.id,
    props: {
      name: game.name,
    }
  });

  return entity;
}