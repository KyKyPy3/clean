import { makeApiUrl, makeHttpClient } from "@/src/main/http";
import { useUserList } from "@user/application/list";
import { useTranslation } from "react-i18next";
import { UserRepositoryImpl } from "@user/infrastructure/gateway/backendRepository";

export function UserList() {
  const { t } = useTranslation()
  const { users, isFetchUsersLoading, isFetchUsersSuccess } = useUserList(
    new UserRepositoryImpl(makeApiUrl('/user'), makeHttpClient())
  );

  return (
    <div className='w-screen min-h-screen p-2 flex flex-col justify-center items-center default-bg'>
      {isFetchUsersLoading && <div>Loading...</div>}
      {isFetchUsersSuccess && (
        <div>
          {users?.map(user => (
            <div key={user.id}>
              <div>{user?.id}</div>
              <div>{user?.getProps().email.value}</div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
