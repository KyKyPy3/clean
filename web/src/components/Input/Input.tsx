import { useState } from 'react'
import type { HTMLInputTypeAttribute } from 'react'
import { Eye, EyeOff } from 'lucide-react'

interface InputProps {
  id: string
  name?: string
  autocomplete?: boolean
  placeholder?: string
  label?: string | React.ReactElement
  type: HTMLInputTypeAttribute
  value?: any
  error?: boolean
  helperText?: string
  onChange?: (event: React.ChangeEvent<HTMLInputElement>) => void
  hook?: object
}

const Input: React.FC<InputProps> = props => {
  const [isPasswordVisible, setIsPasswordVisible] = useState<boolean>(false)
  const [type, setType] = useState<HTMLInputTypeAttribute>(props.type)

  const handlePasswordVisibility = () => {
    if(type === 'text') {
      setIsPasswordVisible(false)
      return setType('password')
    }

    setIsPasswordVisible(true)
    setType('text')
  }

  return(
    <div className='flex flex-col'>
      {props.label && (
        <p className={`mb-1 text-md ${props.error ? 'text-red-500' : 'text-gray-600'}`}>
          {props.label}
        </p>
      )}
      <div className='w-full flex flex-row items-center'>
        <input
          className={`w-full rounded-lg outline-none p-1 border-2 ${props.error ? 'border-red-500' : 'border-gray-500'}`}
          id={props.id}
          name={props.name}
          placeholder={props.placeholder}
          type={type}
          value={props.value}
          autoComplete="{{autocomplete ? 'on' : 'off'}"
          onChange={props.onChange}
          {...props.hook}
        />
        {(props.type === 'password' && !isPasswordVisible) && (
          <button
            type='button'
            className='bg-white rounded-lg border-2 border-gray-500 p-1 ml-3'
            onClick={handlePasswordVisibility}
          >
            <Eye className='text-slate-700' />
          </button>
        )}
        {isPasswordVisible && (
          <button
            type='button'
            className='bg-white rounded-lg border-2 border-gray-500 p-1 ml-3'
            onClick={handlePasswordVisibility}
          >
            <EyeOff className='text-slate-700' />
          </button>
        )}
      </div>
      {props.error && (
        <p className='text-red-500 mt-1'>
          {props.helperText}
        </p>
      )}
    </div>
  )
}

export default Input