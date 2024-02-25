import type { Email } from "@core";
import { AggregateRoot } from "@base";
import { randomID } from "@core";

import type { FullName } from "@user/domain/vo/full-name";

export interface UserProps {
  email: Email;
  fullname: FullName;
}

export interface CreateUserProps {
  email: Email;
  fullname: FullName;
}

export class UserEntity extends AggregateRoot<UserProps> {
  static create(props: CreateUserProps): UserEntity {
    const id = randomID();

    const user = new UserEntity({ id, props });

    return user;
  }

  validate(): void {}
}
