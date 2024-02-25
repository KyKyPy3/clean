import { useMutation } from "@tanstack/react-query";

import { Email } from "@core"
import { useLogger } from "@/src/main/logger";

import type { SessionRepository } from "./session.port"

export type LoginParams = {
  email: string;
  password: string;
}

export const useSessionCreate = (repository: SessionRepository) => {
  const logger = useLogger();

  const createSession = useMutation({
    mutationFn: (params: LoginParams) => {
      // Create email value object
      const email = new Email({ value: params.email });

      return repository.login(email, params.password);
    },
    onError: (error) => {
      logger.error(error);
    },
  })

  return createSession;
}
