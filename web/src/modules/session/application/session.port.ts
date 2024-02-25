import type { Email } from "@core";

export interface SessionRepository {
  login(session: Email, password: string): Promise<{ token: string }>;
}