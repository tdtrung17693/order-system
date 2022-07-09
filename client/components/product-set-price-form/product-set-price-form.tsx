import React, { useEffect, useMemo } from 'react'
import { useTranslation } from 'next-i18next'
import { useForm } from 'react-hook-form'
import * as yup from 'yup'
import { yupResolver } from '@hookform/resolvers/yup'
import auth from '../../services/auth'
import { users } from '../../services/users'
import { UserSignIn } from '../../dto/auth.dto'
import { ProductUpdateType } from 'constants/product'
import { SetProductPrice, UpdateProductStock } from 'dto/product.dto'
import { product } from 'services/product'
import { notification } from 'services/notification'
import { handleApiError } from 'utils/error'
import { Button } from 'antd'

interface ProductSetPriceFormProps {
  productId: number
  currentPrice: number
  onUpdated?: () => any
}
export const ProductSetPriceForm: React.FC<ProductSetPriceFormProps> = (
  props
) => {
  const { t } = useTranslation('common')

  const schema = useMemo(() => {
    const schema = yup
      .object({
        price: yup
          .number()
          .moreThan(0, t('price_zero_error'))
          .required(t('required_input')),
      })
      .required()
    return schema
  }, [])

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm({
    resolver: yupResolver(schema),
  })

  useEffect(() => {
    reset({
      price: props.currentPrice,
    })
  }, [props.currentPrice, reset])

  const onSubmit = (data: any) => {
    const payload: SetProductPrice = {
      price: data.price,
      productId: props.productId,
    }

    product
      .setProductPrice(payload)
      .then(() => {
        notification.info(t('action_success'), t('set_product_price_success'))
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
          {t('product_price_label')}
        </label>
        <input
          type="text"
          inputMode="numeric"
          className={`form-control block w-full px-4 py-2 text-xl font-normal text-gray-700 bg-white bg-clip-padding border border-solid border-gray-300 rounded transition ease-in-out m-0 focus:text-gray-700 focus:bg-white focus:border-blue-600 focus:outline-none ${
            errors.quantity && 'border-red-400'
          }`}
          {...register('price', {})}
        />
        {errors.quantity && (
          <p className="text-sm text-red-400 mt-1">
            {errors.quantity?.message as any}
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
