import { useNavigate } from "react-router-dom"
import { User, LogOut, Users } from 'lucide-react'

import MenuButton from "../MenuButton/MenuButton"
import { useSessionStore } from "@session/infrastructure/controller/http/v1/store"

interface LayoutProps {
  children: React.ReactElement
}

const Layout: React.FC<LayoutProps> = props => {
  const navigate = useNavigate()
  const { clearSession } = useSessionStore()

  return(
    <div className='w-screen min-h-screen flex flex-row default-bg'>
      <div className='md:w-1/12 w-3/12'>
        Logo put here
        <div className='flex flex-col gap-2 h-screen p-2'>
          <MenuButton
            id='menu-button-user'
            title='User'
            icon={<Users />}
            onClick={() => navigate("user")}
          />
          <div className="justify-self-end">
            <MenuButton
              id='menu-button-user'
              title='Profile'
              icon={<User />}
              onClick={() => navigate("/profile")}
            />
          </div>
          <div className="justify-self-end">
            <MenuButton
              id='menu-button-user'
              title='Profile'
              icon={<LogOut />}
              onClick={() => clearSession()}
            />
          </div>
        </div>
      </div>
      <div className='md:w-11/12 w-9/12'>
        {props.children}
      </div>
    </div>
  )
}

export default Layout
