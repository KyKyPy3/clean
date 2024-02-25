import type { Email } from "@core";
import type { AggregateID } from "@base";
import { AggregateRoot } from "@base";
import { randomID } from "@core";

export interface RegistrationProps {
  email: Email;
  password: string;
}

export interface CreateRegistrationProps {
  email: Email;
  password: string;
}

export class RegistrationEntity extends AggregateRoot<RegistrationProps> {
  protected readonly _id!: AggregateID;

  static create(props: CreateRegistrationProps): RegistrationEntity {
    const id = randomID();

    const registration = new RegistrationEntity({ id, props });

    return registration;
  }

  validate(): void {}
}