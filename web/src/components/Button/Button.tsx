import { TailSpin } from 'react-loader-spinner'

interface ButtonProps {
  id: string
  type: 'submit' | 'button'
  title: string | React.ReactElement
  isLoading?: boolean
  onClick?: React.MouseEventHandler<HTMLButtonElement>
}

const Button: React.FC<ButtonProps> = props => {
  return(
    <button
      className='flex justify-center w-full p-1 font-bold text-white bg-slate-950 rounded-lg transition duration-500 ease-in-out enabled:hover:bg-white enabled:hover:text-slate-950 border-2 border-slate-950'
      id={props.id}
      type={props.type}
      disabled={props.isLoading}
      onClick={props.onClick}
    >
      {props.isLoading ? (
        <TailSpin
          height='25'
          width='25'
          color='#FFFFFF'
        />
      ) : (
        props.title
      )}
    </button>
  )
}

export default Button