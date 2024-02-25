import { useTranslation } from "react-i18next";
import { UserRepositoryImpl } from "../../../../gateway/backendRepository";
import { makeApiUrl, makeHttpClient } from "@/src/main/http";
import { useUserProfile } from "@/src/modules/user/application/profile";

export function Profile() {
  const { t } = useTranslation()
  const { user, isFetchProfileLoading, isFetchProfileSuccess } = useUserProfile(
    new UserRepositoryImpl(makeApiUrl('/user'), makeHttpClient())
  );

  return (
    <div className='w-screen min-h-screen p-2 flex flex-col justify-center items-center default-bg'>
      {isFetchProfileLoading && <div>Loading...</div>}
      {isFetchProfileSuccess && (
        <div>
          <div>{user?.getProps().id}</div>
          <div>{user?.getProps().email.value}</div>
        </div>
      )}
    </div>
  )
}
