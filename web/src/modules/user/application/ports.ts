import type { UserEntity } from "@user/domain/entity/user"

export interface UserRepository {
  me(): Promise<UserEntity>;
  list(): Promise<UserEntity[]>
}