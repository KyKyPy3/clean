import { useForm } from 'react-hook-form'
import { useTranslation } from 'react-i18next'
import { Link, useNavigate } from 'react-router-dom'
import type { SubmitHandler } from 'react-hook-form'

import { version } from '@/package.json'
import Button from '@components/Button'
import Input from '@components/Input'
import { makeApiUrl, makeHttpClient } from '@/src/main/http'

import { useRegistrationCreate } from '@/src/modules/registration/application/create'
import { RegistrationRepositoryImpl } from '@registration/infrastructure/gateway/backendRepository'

interface IFormInput {
  email: string
  password: string
}

export function SignUp() {
  const { t } = useTranslation()
  const { register, formState: { errors }, handleSubmit, reset } = useForm<IFormInput>()
  const createRegistration = useRegistrationCreate(
    new RegistrationRepositoryImpl(makeApiUrl('/registration'), makeHttpClient())
  )

  const onSubmit: SubmitHandler<IFormInput> = (data) => {
    createRegistration.mutate(data, {
      onSuccess: () => {
        reset();
      },
    });
  }

  return (
    <div className='w-screen min-h-screen p-2 flex flex-col justify-center items-center default-bg'>
      <div className='w-3/4 md:w-1/4 h-fit bg-neutral-50 border-none rounded-xl p-10'>
        <h1 className='text-slate-950 text-5xl font-bold'>{t('signup.create_account')}</h1>
        <form
          onSubmit={handleSubmit(onSubmit)}
          className='flex flex-col mt-5'
        >
          <div className='mt-5'>
            <Input
              id='signup-email'
              type='text'
              autocomplete={true}
              error={errors.email ? true : false}
              placeholder='Email'
              label='Enter your email'
              hook={{...register('email', { required: "Email address is required" })}}
            />
          </div>
          <div>{errors.email?.message}</div>

          <div className='mt-5'>
            <Input
              id='signup-password'
              autocomplete={true}
              error={errors.password ? true : false}
              type='password'
              placeholder='Password'
              label='Enter your password'
              hook={{...register('password', { required: "Password address is required" })}}
            />
          </div>
          <div>{errors.password?.message}</div>

          <div className='mt-5'>
            <Button
              id='signup-submit'
              type='submit'
              isLoading={createRegistration.isPending}
              title={<p className='text-lg'>{t('signup.signup')}</p>}
            />
          </div>

          <div className='flex justify-center mt-5'>
            <Link to="/signin">{t('signin.signin')}</Link>
          </div>
        </form>
      </div>
      <h2 className='mt-5 text-white text-2m'>{t('version')} {version}</h2>
    </div>
  )
}
