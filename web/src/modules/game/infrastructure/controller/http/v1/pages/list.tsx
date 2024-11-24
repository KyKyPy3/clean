import { useTranslation } from "react-i18next";

import { useGameList } from "@game/application/list";
import { makeApiUrl, makeHttpClient } from "@/src/main/http";
import { GameRepositoryImpl } from "@game/infrastructure/gateway/backendRepository";

export function GamesList() {
  const { t } = useTranslation()
  const { games, isFetchGamesLoading, isFetchGamesSuccess } = useGameList(
    new GameRepositoryImpl(makeApiUrl('/game'), makeHttpClient())
  );

  return (
    <div className='w-screen min-h-screen p-2 flex flex-col justify-center items-center default-bg'>
      {isFetchGamesLoading && <div>Loading...</div>}
      {isFetchGamesSuccess && (
        <div>
          {games?.map(game => (
            <div key={game.id}>
              <div>{game?.id}</div>
              <div>{game?.getProps().name}</div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
