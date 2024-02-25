import type { AggregateID } from "@base";
import { AggregateRoot } from "@base";
import type { Email } from "@core";
import { randomID } from "@core";


export interface SessionProps {
  email: Email
  token: string
}

export interface CreateSessionProps {
  email: Email,
  token: string,
}

export class SessionEntity extends AggregateRoot<SessionProps> {
  protected readonly _id!: AggregateID;

  static create(props: CreateSessionProps): SessionEntity {
    const id = randomID();

    const session = new SessionEntity({ id, props });

    return session;
  }

  validate(): void {}
}
