interface MenuButtonProps {
  id: string
  title: string
  icon: React.ReactElement
  onClick?: React.MouseEventHandler<HTMLButtonElement>
}

const MenuButton: React.FC<MenuButtonProps> = props => {
  return(
    <button
      className='w-full p-1 font-bold text-slate-950 bg-white rounded-lg'
      id={props.id}
      type='button'
      onClick={props.onClick}
    >
      <div className='flex flex-col items-center'>
        {props.icon}
        <p>{props.title}</p>
      </div>
    </button>
  )
}

export default MenuButton
