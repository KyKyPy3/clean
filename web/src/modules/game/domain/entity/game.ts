import { AggregateRoot } from "@base";
import { randomID } from "@core";

export interface GameProps {
  name: string;
}

export interface CreateGameProps {
  name: string;
}

export class GameEntity extends AggregateRoot<GameProps> {
  static create(props: CreateGameProps): GameEntity {
    const id = randomID();

    const game = new GameEntity({ id, props });

    return game;
  }

  validate(): void {}
}