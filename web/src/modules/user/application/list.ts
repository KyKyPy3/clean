import { useQuery } from "@tanstack/react-query";

import { useLogger } from "@/src/main/logger";

import type { UserEntity } from "../domain/entity/user";

import type { UserRepository } from "./ports";

export const useUserList = (repository: UserRepository) => {
  const logger = useLogger();

  const { data, isLoading, isError, isSuccess } = useQuery<UserEntity[]>({
    queryKey: ["users"],
    queryFn: () => repository.list(),
  })

  return {
    users: data,
    isFetchUsersLoading: isLoading,
    isFetchUsersSuccess: isSuccess,
    isFetchUsersError: isError,
  };
}