import { Button, Popconfirm, Space, Table, Tag, Tooltip } from 'antd'
import { TablePaginationConfig } from 'antd/lib/table'
import Column from 'antd/lib/table/Column'
import { FilterValue } from 'antd/lib/table/interface'
import { OrderStatusTag } from 'components/order-status-tag/order-status-tag'
import { OrderStatus } from 'constants/order'
import { ItemsPerPage } from 'constants/pagination'
import { AuthContext } from 'context/auth.context'
import dayjs from 'dayjs'
import { Order } from 'dto/order.dto'
import { buildPaginationRequest, PaginationResponse } from 'dto/pagination.dto'
import type { GetServerSideProps, NextPage } from 'next'
import { useTranslation } from 'next-i18next'
import { serverSideTranslations } from 'next-i18next/serverSideTranslations'
import Head from 'next/head'
import Link from 'next/link'
import { useRouter } from 'next/router'
import { useCallback, useContext, useEffect, useState } from 'react'
import auth from 'services/auth'
import { http } from 'services/http'
import { notification } from 'services/notification'
import { order } from 'services/order'
import { Maybe } from 'types/maybe'
import { handleApiError } from 'utils/error'

const UserOrderManagePage: NextPage = () => {
  const authCtx = useContext(AuthContext)
  const router = useRouter()
  const { t } = useTranslation('common')
  const [products, setProducts] = useState<Order[]>([])
  const [filteredStatus, setFilteredStatus] = useState<Maybe<OrderStatus>>(null)
  const [{ current, pageSize }, _] = useState({
    current: 1,
    pageSize: ItemsPerPage,
  })
  const [total, setTotal] = useState(-1)

  const fetchData = useCallback(
    (pageIndex: number, status: Maybe<OrderStatus> = null) => {
      setFilteredStatus(status)
      return http
        .get<PaginationResponse<Order>>('/orders', {
          params: {
            ...buildPaginationRequest(
              pageIndex,
              ItemsPerPage,
              status ? { status } : {}
            ),
          },
        })
        .then(({ data }) => {
          setTotal(data.total)

          setProducts(data.items)
        })
    },
    []
  )

  useEffect(() => {
    fetchData(current - 1)
  }, [fetchData, current])

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
        fetchData(current - 1)
      })
      .catch((error) => {
        handleApiError(t, error)
      })
  }

  const handleTableChange = (
    newPagination: TablePaginationConfig,
    filters: Record<string, FilterValue>
  ) => {
    fetchData(
      (newPagination.current || 1) - 1,
      filters.status ? (filters.status[0] as OrderStatus) : null
    )
  }

  const exportCSV = () => {
    order
      .exportCsv(false, filteredStatus)
      .catch((error) => handleApiError(t, error))
  }

  return (
    <div className="p-8">
      <Head>
        <title>Order Management System | My Orders</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <h1 className="text-5xl mb-4 text-center">{t('title_orders')}</h1>
      <main className="flex flex-col justify-start items-center min-h-screen p-16">
        <div className="mt-4 text-base w-full">
          <Tooltip title={t('export_selected_csv_tooltip')}>
            <Button onClick={exportCSV}>{t('export_selected_csv')}</Button>
          </Tooltip>
          <div className="mt-4">
            <Table
              className="min-w-full"
              dataSource={products}
              pagination={{ pageSize, total }}
              onChange={handleTableChange as any}
            >
              <Column
                title={t('order_id_label') as string}
                dataIndex="id"
                key="id"
                width="5%"
              />
              <Column
                title={t('order_vendor_label') as string}
                dataIndex="vendorName"
                key="vendorName"
                width="25%"
              />
              <Column
                title={t('order_created_label') as string}
                dataIndex="createdAt"
                key="createdAt"
                width="20%"
                render={(_: any, record: Order) => {
                  return dayjs(record.createdAt).format('MMM, DD YYYY hh:mm')
                }}
              />
              <Column
                title={t('order_updated_label') as string}
                dataIndex="statusChangeTime"
                key="statusChangeTime"
                width="20%"
                render={(_: any, record: Order) => {
                  return dayjs(record.statusChangeTime).format(
                    'MMM, DD YYYY hh:mm'
                  )
                }}
              />
              <Column
                title={t('order_price_label') as string}
                dataIndex="totalPrice"
                key="totalPrice"
                width="20%"
              />
              <Column
                title={t('order_status_label') as string}
                dataIndex="status"
                key="status"
                width="10%"
                filterMultiple={false}
                filters={[
                  {
                    text: t('order_cancelled_text'),
                    value: OrderStatus.Cancelled,
                  },
                  {
                    text: t('order_paid_text'),
                    value: OrderStatus.Paid,
                  },
                  {
                    text: t('order_shipping_text'),
                    value: OrderStatus.Shipping,
                  },
                  {
                    text: t('order_shipped_text'),
                    value: OrderStatus.Shipped,
                  },
                ]}
                render={(_: any, record: Order) => (
                  <OrderStatusTag status={record.status} />
                )}
              />
              <Column
                title="Action"
                key="action"
                width="10%"
                render={(_: any, record: Order) => {
                  return (
                    <Space size="middle">
                      <Link href={`/orders/${record.id}`}>
                        <Button type="primary">{t('order_view_text')}</Button>
                      </Link>
                      {order.isCancellableState(record) && (
                        <Popconfirm
                          title={t('order_cancel_confirm')}
                          onConfirm={() => cancelOrder(record.id)}
                          okText={t('confirm_ok_text')}
                          cancelText={t('confirm_cancel_text')}
                        >
                          <Button type="primary" danger>
                            {t('order_cancel_text')}
                          </Button>
                        </Popconfirm>
                      )}
                    </Space>
                  )
                }}
              />
            </Table>
          </div>
        </div>
      </main>
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

export default UserOrderManagePage
