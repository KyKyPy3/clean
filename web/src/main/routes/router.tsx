import { RouterProvider, createBrowserRouter } from 'react-router-dom'

import Error from '@components/Error'
import { ProtectedRoute } from './protected'

import { SignUp } from '@registration/infrastructure/controller/http/v1/pages/signup'
import SignIn from '@session/infrastructure/controller/http/v1/pages/signin'
import { Profile, UserList } from '@user/infrastructure/controller/http/v1/pages'
import { GamesList } from '@game/infrastructure/controller/http/v1/pages'

export const AppRoutes = () => {
  const router = createBrowserRouter([
    {
      path: "/signup",
      element: <SignUp />,
      errorElement: <Error />,
    },
    {
      path: "/signin",
      element: <SignIn />,
      errorElement: <Error />,
    },
    {
      element: <ProtectedRoute />,
      children: [
        {
          path: "/",
          element: <GamesList />,
          errorElement: <Error />,
        },
        {
          path: "/game",
          element: <GamesList />,
          errorElement: <Error />,
        },
        {
          path: "/profile",
          element: <Profile />,
          errorElement: <Error />,
        },
        {
          path: "/user",
          element: <UserList />,
          errorElement: <Error />,
        }
      ]
    },
  ])

  return <RouterProvider router={router} />
}
