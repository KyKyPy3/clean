import { Email } from "@/src/core";
import { RegistrationEntity } from "@registration/domain/entity/registration";

export type RegistrationDTO = {
  id: string;
  email: string;
  password: string;
}

export function mapToDomain(reg: RegistrationDTO): RegistrationEntity {
  const email = new Email({ value: reg.email })

  const entity = new RegistrationEntity({
    id: reg.id,
    props: {
      email: email,
      password: reg.password,
    }
  });

  return entity;
}

export function mapToRequest(reg: RegistrationEntity): RegistrationDTO {
  const copy = reg.getProps();

  return {
    id: copy.id,
    email: copy.email.value,
    password: copy.password,
  }
}
