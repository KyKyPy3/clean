import { useNavigate, Link } from 'react-router-dom'
import { useForm } from 'react-hook-form'
import { useTranslation } from 'react-i18next'
import type { SubmitHandler } from 'react-hook-form'

import { version } from '@/package.json';
import Button from '@components/Button';
import Input from '@components/Input';

import { makeApiUrl, makeHttpClient } from "@/src/main/http"
import { useSessionCreate } from '@session/application/session'
import { SessionRepositoryImpl } from "@session/infrastructure/gateway/backendRepository"
import { useSessionStore } from "@session/infrastructure/controller/http/v1/store"

interface IFormInput {
  email: string
  password: string
}

const SignIn: React.FC = () => {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const { register, formState: { errors }, handleSubmit, reset } = useForm<IFormInput>()

  const createSession = useSessionCreate(
    new SessionRepositoryImpl(makeApiUrl('/auth/login'), makeHttpClient())
  )
  const { setSession, token } = useSessionStore()

  const onSubmit: SubmitHandler<IFormInput> = (data) => {
    createSession.mutate(data, {
      onSuccess: data => {
        setSession(data);

        navigate("/");
      },
      onError: () => {
        reset();
      },
    });
  }

  return(
    <div className='w-screen min-h-screen p-2 flex flex-col justify-center items-center default-bg'>
      <div className='w-3/4 md:w-1/4 h-fit bg-neutral-50 border-none rounded-xl p-10'>
        <h1 className='text-slate-950 text-5xl font-bold'>Login</h1>

        <form
          onSubmit={handleSubmit(onSubmit)}
          className='flex flex-col mt-5'
        >
          <div>
            <Input
              id='login-email'
              type='text'
              placeholder='Email'
              label='Enter your email'
              hook={{...register('email')}}
            />
          </div>

          <div className='mt-5'>
            <Input
              id='login-password'
              type='password'
              placeholder='Password'
              label='Enter your password'
              hook={{...register('password')}}
            />
          </div>

          <div className='mt-5'>
            <Button
              id='login-submit'
              type='submit'
              isLoading={createSession.isPending}
              title={<p className='text-lg'>Sign in</p>}
            />
          </div>
          <div className='flex justify-center mt-5'>
            <Link to="/signup">Create Account</Link>
          </div>
        </form>
      </div>
      <h2 className='mt-5 text-white text-2m'>{t('version')} {version}</h2>
    </div>
  )
}

export default SignIn
