import { Button, Popconfirm, Tag } from 'antd'
import { OrderStatusTag } from 'components/order-status-tag/order-status-tag'
import { OrderStatus } from 'constants/order'
import { AuthContext } from 'context/auth.context'
import Decimal from 'decimal.js'
import { Order } from 'dto/order.dto'
import type { GetServerSideProps, NextPage } from 'next'
import { useTranslation } from 'next-i18next'
import { serverSideTranslations } from 'next-i18next/serverSideTranslations'
import Head from 'next/head'
import Image from 'next/image'
import { useRouter } from 'next/router'
import { useCallback, useContext, useEffect, useMemo, useState } from 'react'
import auth from 'services/auth'
import { notification } from 'services/notification'
import { order } from 'services/order'
import { Maybe } from 'types/maybe'
import { handleApiError } from 'utils/error'

const OrderPage: NextPage = () => {
  const authCtx = useContext(AuthContext)
  const router = useRouter()
  const { t } = useTranslation('common')
  const [currentOrder, setCurrentOrder] = useState<Maybe<Order>>(null)
  const orderId = useMemo(() => {
    return parseInt((router.query.id as string) || '0', 10)
  }, [router.query.id])

  const fetchOrder = useCallback(
    (orderId: number) => {
      order
        .getOrder(orderId)
        .then((data) => setCurrentOrder(data))
        .catch((error) => handleApiError(t, error))
    },
    [t]
  )

  useEffect(() => {
    fetchOrder(orderId)
  }, [fetchOrder, orderId])

  useEffect(() => {
    if (auth.initialized && !authCtx.user) {
      router.push('/auth/signin')
    }
  }, [authCtx.user, router])

  function cancelOrder(orderId: number) {
    order
      .cancelOrder(orderId)
      .then(() => {
        notification.info(t('action_success'), t('cancel_order_success'))
        fetchOrder(orderId)
      })
      .catch((error) => {
        handleApiError(t, error)
      })
  }

  return (
    <div className="p-8">
      <Head>
        <title>Order Management System | My Orders</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      {currentOrder != null && (
        <>
          <div className=" text-center">
            <h2 className="text-5xl mb-4">
              {t('title_order')}: {currentOrder.id}
            </h2>
            <OrderStatusTag status={currentOrder.status} />
            {currentOrder.status !== OrderStatus.Cancelled &&
              currentOrder.status !== OrderStatus.Shipping &&
              currentOrder.status !== OrderStatus.Shipped && (
                <Popconfirm
                  title={t('order_cancel_confirm')}
                  onConfirm={() => cancelOrder(currentOrder.id)}
                  okText={t('confirm_ok_text')}
                  cancelText={t('confirm_cancel_text')}
                >
                  <Button type="primary" danger>
                    {t('order_cancel_text')}
                  </Button>
                </Popconfirm>
              )}
          </div>
          <main className="flex flex-col justify-start items-center min-h-screen p-16">
            <div className="mt-4 text-base w-full">
              <div className="mt-4">
                <h3 className="text-lg py-2 px-4 rounded bg-gray-100">
                  Payment info
                </h3>
                <div className="ml-4">
                  <span className="font-bold mr-4">
                    {t('payment_fullname')}:
                  </span>
                  <span className="">{currentOrder.recipientName}</span>
                </div>
                <div className="ml-4">
                  <span className="font-bold mr-4">
                    {t('payment_phone_number')}:
                  </span>
                  <span className="">{currentOrder.recipientPhone}</span>
                </div>
                <div className="ml-4">
                  <span className="font-bold mr-4">{t('payment_method')}:</span>
                  <span className="">{t(currentOrder.paymentMethodName)}</span>
                </div>
                <div className="ml-4">
                  <span className="font-bold mr-4">
                    {t('payment_address')}:
                  </span>
                  <span className="">{currentOrder.shippingAddress}</span>
                </div>
              </div>
              <div className="mt-4">
                <h3 className="text-lg py-2 px-4 rounded bg-gray-100">
                  {t('order_items_text')}
                </h3>
                <table className="w-full text-sm lg:text-base" cellSpacing="0">
                  <thead>
                    <tr className="h-12 uppercase">
                      <th className="table-cell w-40"></th>
                      <th className="text-left">{t('product_name_label')}</th>
                      <th className="text-center">
                        <span>{t('order_item_quantity_label')}</span>
                      </th>
                      <th className="text-right table-cell">
                        {t('product_price_label')}
                      </th>
                      <th className="text-right">
                        {t('order_item_total_label')}
                      </th>
                    </tr>
                  </thead>
                  <tbody>
                    {(currentOrder?.items || []).map((item) => {
                      return (
                        <>
                          <tr key={`${item.productId}-${item}`}>
                            <td className=" pb-4 table-cell">
                              <Image
                                src="/images/default-product-image.png"
                                width={160}
                                height={160}
                                alt={item.productName}
                              />
                            </td>
                            <td>{item.productName}</td>
                            <td className="text-center mt-6">
                              {item.quantity}
                            </td>
                            <td className="hidden text-right md:table-cell">
                              <span className="text-sm lg:text-base font-medium">
                                {item.unitPrice}
                              </span>
                            </td>
                            <td className="text-right">
                              <span className="text-sm lg:text-base font-medium">
                                {new Decimal(item.unitPrice)
                                  .times(item.quantity)
                                  .toFixed(2)}
                              </span>
                            </td>
                          </tr>
                        </>
                      )
                    })}
                  </tbody>
                </table>
                <div className="my-4 mt-6 -mx-2 lg:flex">
                  <div className="w-full">
                    <div className="p-4">
                      <div className="flex justify-between border-b">
                        <div className="lg:px-4 lg:py-2 m-2 text-lg lg:text-xl font-bold text-center text-gray-800">
                          {t('cart_order_total_text')}
                        </div>
                        <div className="lg:px-4 lg:py-2 m-2 lg:text-lg font-bold text-center text-gray-900">
                          {(currentOrder.items || [])
                            .reduce(
                              (s, i) =>
                                new Decimal(i.quantity)
                                  .times(new Decimal(i.unitPrice))
                                  .plus(s),
                              new Decimal(0)
                            )
                            .toFixed(2)}
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </main>
        </>
      )}
    </div>
  )
}
export const getServerSideProps: GetServerSideProps<any> = async ({
  locale,
}) => {
  return {
    props: {
      ...(await serverSideTranslations(locale || 'en', ['common'])),
    },
  }
}

export default OrderPage
