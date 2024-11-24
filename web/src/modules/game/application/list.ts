import { useQuery } from "@tanstack/react-query";

import { useLogger } from "@/src/main/logger";

import type { GameEntity } from "../domain/entity/game";

import type { GameRepository } from "./ports";

export const useGameList = (repository: GameRepository) => {
  const logger = useLogger();

  const { data, isLoading, isError, isSuccess } = useQuery<GameEntity[]>({
    queryKey: ["games"],
    queryFn: () => repository.list(),
  })

  return {
    games: data,
    isFetchGamesLoading: isLoading,
    isFetchGamesSuccess: isSuccess,
    isFetchGamesError: isError,
  };
}