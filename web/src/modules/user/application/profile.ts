import { useQuery } from "@tanstack/react-query";

import { useLogger } from "@/src/main/logger";

import type { UserEntity } from "../domain/entity/user";

import type { UserRepository } from "./ports";

export const useUserProfile = (repository: UserRepository) => {
  const logger = useLogger();

  const { data, isLoading, isError, isSuccess } = useQuery<UserEntity>({
    queryKey: ["me"],
    queryFn: () => repository.me(),
  })

  return {
    user: data,
    isFetchProfileLoading: isLoading,
    isFetchProfileSuccess: isSuccess,
    isFetchProfileError: isError,
  };
}
