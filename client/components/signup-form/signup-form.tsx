import React, { useMemo } from 'react'
import { useTranslation } from 'next-i18next'
import { useForm } from 'react-hook-form'
import * as yup from 'yup'
import { yupResolver } from '@hookform/resolvers/yup'
import { parseRole, UserRole } from '../../constants/user-role'
import auth from '../../services/auth'
import { users } from '../../services/users'
import { UserSignUp } from '../../dto/auth.dto'

interface SignUpFormProps {
  role: UserRole
}

export const SignUpForm: React.FC<SignUpFormProps> = (props) => {
  const { t } = useTranslation('common')

  const schema = useMemo(() => {
    const schema = yup
      .object({
        name: yup.string().required(t('required_input')),
        email: yup
          .string()
          .email(t('email_invalid'))
          .required(t('required_input')),
        password: yup
          .string()
          .min(6, t('password_length_min'))
          .required(t('required_input')),
        confirmPassword: yup
          .string()
          .oneOf([yup.ref('password'), null], t('password_mismatched'))
          .required(t('required_input')),
      })
      .required()
    return schema
  }, [t])
  const {
    register,
    handleSubmit,
    watch,
    setError,
    formState: { errors },
  } = useForm({
    resolver: yupResolver(schema),
  })
  const onSubmit = (data: any) => {
    const userData: UserSignUp = {
      ...data,
      role: parseRole(data.role),
    }
    users
      .register(userData)
      .then((response) => {
        auth.setAccessToken(response.accessToken)
        window.location.href = '/'
      })
      .catch((err) => {
        setError('apiError', { message: err.message })
      })
  }
  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      {errors.apiError && <div>{errors.apiError!.message as any}</div>}
      <input type="hidden" value={props.role} {...register('role')} />
      <div className="mb-3 xl:w-96">
        <label className="form-label inline-block mb-2 text-gray-700 text-xl">
          {t('name_label')}
        </label>
        <input
          type="text"
          className={`form-control block w-full px-4 py-2 text-xl font-normal text-gray-700 bg-white bg-clip-padding border border-solid border-gray-300 rounded transition ease-in-out m-0 focus:text-gray-700 focus:bg-white focus:border-blue-600 focus:outline-none ${
            errors.name && 'border-red-400'
          }`}
          {...register('name')}
        />
        {errors.name && (
          <p className="text-sm text-red-400 mt-1">
            {errors.name?.message as any}
          </p>
        )}
      </div>
      <div className="mb-3 xl:w-96">
        <label className="form-label inline-block mb-2 text-gray-700 text-xl">
          {t('email_label')}
        </label>
        <input
          type="email"
          inputMode="email"
          className={`form-control block w-full px-4 py-2 text-xl font-normal text-gray-700 bg-white bg-clip-padding border border-solid border-gray-300 rounded transition ease-in-out m-0 focus:text-gray-700 focus:bg-white focus:border-blue-600 focus:outline-none ${
            errors.email && 'border-red-400'
          }`}
          {...register('email')}
        />
        {errors.email && (
          <p className="text-sm text-red-400 mt-1">
            {errors.email?.message as any}
          </p>
        )}
      </div>
      <div className="mb-3 xl:w-96">
        <label className="form-label inline-block mb-2 text-gray-700 text-xl">
          {t('password_label')}
        </label>
        <input
          type="password"
          className={`form-control block w-full px-4 py-2 text-xl font-normal text-gray-700 bg-white bg-clip-padding border border-solid border-gray-300 rounded transition ease-in-out m-0 focus:text-gray-700 focus:bg-white focus:border-blue-600 focus:outline-none ${
            errors.password && 'border-red-400'
          }`}
          {...register('password')}
        />
        {errors.password && (
          <p className="text-sm text-red-400 mt-1">
            {errors.email?.message as any}
          </p>
        )}
      </div>
      <div className="mb-3 xl:w-96">
        <label className="form-label inline-block mb-2 text-gray-700 text-xl">
          {t('confirm_password_label')}
        </label>
        <input
          type="password"
          className={`form-control block w-full px-4 py-2 text-xl font-normal text-gray-700 bg-white bg-clip-padding border border-solid border-gray-300 rounded transition ease-in-out m-0 focus:text-gray-700 focus:bg-white focus:border-blue-600 focus:outline-none ${
            errors.confirmPassword && 'border-red-400'
          }`}
          {...register('confirmPassword')}
        />
        {errors.confirmPassword && (
          <p className="text-sm text-red-400 mt-1">
            {errors.confirmPassword?.message as any}
          </p>
        )}
      </div>
      <div className="mt-6">
        <button
          type="submit"
          className="block min-w-full px-6 py-2.5 bg-blue-600 text-white font-medium text-xs leading-tight uppercase rounded shadow-md hover:bg-blue-700 hover:shadow-lg focus:bg-blue-700 focus:shadow-lg focus:outline-none focus:ring-0 active:bg-blue-800 active:shadow-lg transition duration-150 ease-in-out"
        >
          Sign Up
        </button>
      </div>
    </form>
  )
}
