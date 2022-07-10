import { yupResolver } from '@hookform/resolvers/yup'
import { Button } from 'antd'
import { CreateProduct } from 'dto/product.dto'
import { useTranslation } from 'next-i18next'
import React, { useMemo } from 'react'
import { useForm } from 'react-hook-form'
import { notification } from 'services/notification'
import { product } from 'services/product'
import { handleApiError } from 'utils/error'
import * as yup from 'yup'

interface ProductCreateFormProps {
  onCreated?: () => any
}
export const ProductCreateForm: React.FC<ProductCreateFormProps> = (props) => {
  const { t } = useTranslation('common')

  const schema = useMemo(() => {
    const schema = yup
      .object({
        name: yup.string().required(t('required_input')),
        description: yup.string().required(t('required_input')),
      })
      .required()
    return schema
  }, [t])

  const {
    register,
    handleSubmit,
    setError,
    reset,
    formState: { errors },
  } = useForm({
    resolver: yupResolver(schema),
  })

  const onSubmit = (data: any) => {
    const payload: CreateProduct = {
      name: data.name,
      descriptiton: data.description,
    }

    product
      .createProduct(payload)
      .then(() => {
        notification.info(t('action_success'), t('create_product_success'))
        reset()
        if (props.onCreated) {
          props.onCreated()
        }
      })
      .catch((error) => {
        handleApiError(t, error)
      })
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      {errors.apiError && <div>{errors.apiError!.message as any}</div>}

      <div className="mb-3">
        <label className="form-label inline-block mb-2 text-gray-700 text-xl">
          {t('product_name_label')}
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
      <div className="mb-3">
        <label className="form-label inline-block mb-2 text-gray-700 text-xl">
          {t('product_description_label')}
        </label>
        <textarea
          className={`form-control block w-full px-4 py-2 text-xl font-normal text-gray-700 bg-white bg-clip-padding border border-solid border-gray-300 rounded transition ease-in-out m-0 focus:text-gray-700 focus:bg-white focus:border-blue-600 focus:outline-none ${
            errors.description && 'border-red-400'
          }`}
          {...register('description')}
        />
        {errors.description && (
          <p className="text-sm text-red-400 mt-1">
            {errors.description?.message as any}
          </p>
        )}
      </div>
      <div className="mt-6">
        <Button type="primary" htmlType="submit">
          Create
        </Button>
      </div>
    </form>
  )
}
