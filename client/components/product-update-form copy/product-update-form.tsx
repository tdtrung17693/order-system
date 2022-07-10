import { yupResolver } from '@hookform/resolvers/yup'
import { Button } from 'antd'
import { ProductUpdateType } from 'constants/product'
import { UpdateProductStock } from 'dto/product.dto'
import { useTranslation } from 'next-i18next'
import React, { useMemo } from 'react'
import { useForm } from 'react-hook-form'
import { notification } from 'services/notification'
import { product } from 'services/product'
import { handleApiError } from 'utils/error'
import * as yup from 'yup'

interface ProductUpdateFormProps {
  updateType: ProductUpdateType
  productId: number
  onUpdated?: () => any
}
export const ProductUpdateForm: React.FC<ProductUpdateFormProps> = (props) => {
  const { t } = useTranslation('common')

  const schema = useMemo(() => {
    const schema = yup
      .object({
        quantity: yup.number().required(t('required_input')),
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
    const payload: UpdateProductStock = {
      productId: props.productId,
      description: data.description,
      type: props.updateType,
      quantity: data.quantity,
    }

    product
      .updateProductStock(payload)
      .then(() => {
        notification.info(
          t('action_success'),
          t('update_product_stock_success')
        )
        reset()
        if (props.onUpdated) {
          props.onUpdated()
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
          {t('product_update_quantity_label')}
        </label>
        <input
          type="number"
          inputMode="numeric"
          className={`form-control block w-full px-4 py-2 text-xl font-normal text-gray-700 bg-white bg-clip-padding border border-solid border-gray-300 rounded transition ease-in-out m-0 focus:text-gray-700 focus:bg-white focus:border-blue-600 focus:outline-none ${
            errors.quantity && 'border-red-400'
          }`}
          {...register('quantity')}
        />
        {errors.quantity && (
          <p className="text-sm text-red-400 mt-1">
            {errors.quantity?.message as any}
          </p>
        )}
      </div>
      <div className="mb-3">
        <label className="form-label inline-block mb-2 text-gray-700 text-xl">
          {t('product_update_description_label')}
        </label>
        <textarea
          className={`form-control block w-full px-4 py-2 text-xl font-normal text-gray-700 bg-white bg-clip-padding border border-solid border-gray-300 rounded transition ease-in-out m-0 focus:text-gray-700 focus:bg-white focus:border-blue-600 focus:outline-none ${
            errors.descriptioon && 'border-red-400'
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
          Update
        </Button>
      </div>
    </form>
  )
}
