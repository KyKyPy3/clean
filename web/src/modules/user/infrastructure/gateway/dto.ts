import { Email } from "@/src/core";

import { UserEntity } from "../../domain/entity/user";
import { FullName } from "../../domain/vo/full-name";

export type UserDTO = {
  id: string;
  email: string;
  name: string;
  surname: string;
  middlename: string;
}

export function mapToDomain(user: UserDTO): UserEntity {
  const email = new Email({ value: user.email })
  const fullname = new FullName({ name: user.name, surname: user.surname, middlename: user.middlename })

  const entity = new UserEntity({
    id: user.id,
    props: {
      email: email,
      fullname: fullname,
    }
  });

  return entity;
}

export function mapToRequest(user: UserEntity): UserDTO {
  const copy = user.getProps();

  return {
    id: copy.id,
    email: copy.email.value,
    name: copy.fullname.name,
    surname: copy.fullname.surname,
    middlename: copy.fullname.middlename,
  }
}