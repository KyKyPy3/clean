import { useMutation } from "@tanstack/react-query";
import { Email } from "@core";
import { useLogger } from "@/src/main/logger";

import { RegistrationEntity } from "../domain/entity/registration"

import type { RegistrationRepository } from "./ports"

export type RegistrationParams = {
  email: string;
  password: string;
}

export const useRegistrationCreate = (repository: RegistrationRepository) => {
  const logger = useLogger();

  const createRegistration = useMutation({
    mutationFn: (params: RegistrationParams) => {
      // Create email value object
      const email = new Email({ value: params.email })

      // Create registration entity
      const registration = RegistrationEntity.create({
        email: email,
        password: params.password,
      })

      return repository.register(registration);
    },
    onError: (error) => {
      logger.error(error);
    },
  })

  return createRegistration;
}
