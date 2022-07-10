import { yupResolver } from '@hookform/resolvers/yup'
import { Button, Select } from 'antd'
import { PaymentInfo, PaymentMethod } from 'dto/payment.dto'
import { CreateProduct } from 'dto/product.dto'
import { useTranslation } from 'next-i18next'
import React, { useMemo } from 'react'
import { useForm } from 'react-hook-form'
import { notification } from 'services/notification'
import { product } from 'services/product'
import { handleApiError } from 'utils/error'
import * as yup from 'yup'
interface PaymentFormProps {
  onSubmit: (data: PaymentInfo) => any
  paymentMethods: PaymentMethod[]
}
export const PaymentForm: React.FC<PaymentFormProps> = (props) => {
  const { t } = useTranslation('common')

  const schema = useMemo(() => {
    const schema = yup
      .object({
        paymentMethodId: yup.string().required(t('required_input')),
        recipientPhone: yup.string().required(t('required_input')),
        recipientName: yup.string().required(t('required_input')),
        recipientAddress: yup.string().required(t('required_input')),
      })
      .required()
    return schema
  }, [t])

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm({
    resolver: yupResolver(schema),
  })

  const onSubmit = (data: any) => {
    const payload: PaymentInfo = {
      paymentMethodId: data.paymentMethodId,
      recipientAddress: data.recipientAddress,
      recipientPhone: data.recipientPhone,
      recipientName: data.recipientName,
    }

    props.onSubmit(payload)
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      <div className="mb-3">
        <label className="form-label inline-block mb-2 text-gray-700 text-xl">
          {t('payment_fullname')}
        </label>
        <input
          type="text"
          className={`form-control block w-full px-4 py-2 text-xl font-normal text-gray-700 bg-white bg-clip-padding border border-solid border-gray-300 rounded transition ease-in-out m-0 focus:text-gray-700 focus:bg-white focus:border-blue-600 focus:outline-none ${
            errors.recipientName && 'border-red-400'
          }`}
          {...register('recipientName')}
        />
        {errors.recipientName && (
          <p className="text-sm text-red-400 mt-1">
            {errors.recipientName?.message as any}
          </p>
        )}
      </div>
      <div className="mb-3">
        <label className="form-label inline-block mb-2 text-gray-700 text-xl">
          {t('payment_phone_number')}
        </label>
        <input
          type="text"
          className={`form-control block w-full px-4 py-2 text-xl font-normal text-gray-700 bg-white bg-clip-padding border border-solid border-gray-300 rounded transition ease-in-out m-0 focus:text-gray-700 focus:bg-white focus:border-blue-600 focus:outline-none ${
            errors.recipientPhone && 'border-red-400'
          }`}
          {...register('recipientPhone')}
        />
        {errors.recipientPhone && (
          <p className="text-sm text-red-400 mt-1">
            {errors.recipientPhone?.message as any}
          </p>
        )}
      </div>
      <div className="mb-3">
        <label className="form-label inline-block mb-2 text-gray-700 text-xl">
          {t('payment_address')}
        </label>
        <input
          type="text"
          className={`form-control block w-full px-4 py-2 text-xl font-normal text-gray-700 bg-white bg-clip-padding border border-solid border-gray-300 rounded transition ease-in-out m-0 focus:text-gray-700 focus:bg-white focus:border-blue-600 focus:outline-none ${
            errors.recipientAddress && 'border-red-400'
          }`}
          {...register('recipientAddress')}
        />
        {errors.recipientAddress && (
          <p className="text-sm text-red-400 mt-1">
            {errors.recipientAddress?.message as any}
          </p>
        )}
      </div>
      <div className="mb-3">
        <label className="form-label inline-block mb-2 text-gray-700 text-xl">
          {t('payment_method')}
        </label>
        <select
          className={`form-control block w-full px-4 py-2 text-xl font-normal text-gray-700 bg-white bg-clip-padding border border-solid border-gray-300 rounded transition ease-in-out m-0 focus:text-gray-700 focus:bg-white focus:border-blue-600 focus:outline-none ${
            errors.paymentMethodId && 'border-red-400'
          }`}
          {...register('paymentMethodId')}
        >
          {props.paymentMethods.map((method, idx) => (
            <option
              key={method.id}
              value={method.id}
              defaultChecked={idx === 0}
            >
              {t(method.name)}
            </option>
          ))}
        </select>
        {errors.paymentMethodId && (
          <p className="text-sm text-red-400 mt-1">
            {errors.paymentMethodId?.message as any}
          </p>
        )}
      </div>
      <div className="mt-6">
        <Button
          type="primary"
          htmlType="submit"
          size="large"
          className="flex justify-center items-center w-full px-10 py-3 mt-6 font-medium uppercase rounded-full shadow item-center hover:bg-gray-700 focus:shadow-outline focus:outline-none"
        >
          <span className="ml-2 mt-5px flex">{t('submit_order')}</span>
        </Button>
      </div>
    </form>
  )
}
