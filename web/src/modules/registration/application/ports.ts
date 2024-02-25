import type { RegistrationEntity } from "@registration/domain/entity/registration"

export interface RegistrationRepository {
  register(registration: RegistrationEntity): Promise<void>;
}

