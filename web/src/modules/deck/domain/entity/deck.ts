import type { AggregateID} from "@base";
import { AggregateRoot } from "@base";
import { randomID } from "@core";

export interface DeckProps {
  name: string;
  cards: string;
}

export interface CreateDeckProps {
  name: string;
  cards: string;
}

export class DeckEntity extends AggregateRoot<DeckProps> {
  protected readonly _id!: AggregateID;

  static create(props: CreateDeckProps): DeckEntity {
    const id = randomID();

    const deck = new DeckEntity({ id, props });

    return deck;
  }

  validate(): void {}
}
