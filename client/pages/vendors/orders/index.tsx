import type { GetServerSideProps, NextPage } from 'next'
import Error from 'next/error'
import { useTranslation } from 'next-i18next'
import { serverSideTranslations } from 'next-i18next/serverSideTranslations'
import { useCallback, useContext, useEffect, useState } from 'react'
import { AuthContext } from 'context/auth.context'
import { UserRole } from 'constants/user-role'
import { useRouter } from 'next/router'
import auth from 'services/auth'
import { LayoutDashboard } from 'components/layout/layout-dashboard'
import { ItemsPerPage } from 'constants/pagination'
import {
  PaginationResponse,
  buildPaginationRequest,
  PaginationQuery,
} from 'dto/pagination.dto'
import { Product } from 'dto/product.dto'
import { http } from 'services/http'
import { Table, Tag, Space } from 'antd'
import table, { TablePaginationConfig } from 'antd/lib/table'
import Column from 'antd/lib/table/Column'
import ColumnGroup from 'antd/lib/table/ColumnGroup'
import { FilterValue, SorterResult } from 'antd/lib/table/interface'
import { Order } from 'dto/order.dto'
import dayjs from 'dayjs'
import { OrderStatus } from 'constants/order'

const VendorDashboardProducts: NextPage = () => {
  const authCtx = useContext(AuthContext)
  const router = useRouter()
  const { t } = useTranslation('common')
  const [products, setProducts] = useState<Order[]>([])
  const [{ current, pageSize }, setPagination] = useState({
    current: 1,
    pageSize: ItemsPerPage,
  })
  const [total, setTotal] = useState(-1)

  const fetchData = useCallback((pageIndex: number) => {
    return http
      .get<PaginationResponse<Order>>('/vendors/orders', {
        params: {
          ...buildPaginationRequest(pageIndex, ItemsPerPage),
        },
      })
      .then(({ data }) => {
        setTotal(data.total)

        setProducts(data.items)
      })
  }, [])

  useEffect(() => {
    fetchData(current - 1)
  }, [fetchData, current, total])

  useEffect(() => {
    if (auth.initialized && !authCtx.user) {
      router.push('/auth/signin')
    }
  }, [authCtx.user, router])

  const handleTableChange = (
    newPagination: TablePaginationConfig,
    filters: Record<string, FilterValue>,
    sorter: SorterResult<Product>
  ) => {
    fetchData((newPagination.current || 1) - 1)
  }

  if (authCtx.user && authCtx.user.role != UserRole.Vendor) {
    return <Error statusCode={403} />
  }

  return (
    <div className="px-20">
      <main className="flex flex-col justify-start items-center min-h-screen p-16">
        <LayoutDashboard>
          <div className="text-center mt-4 text-base">
            <div className="p-2 ">
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
                  title={t('order_created_label') as string}
                  dataIndex="createdAt"
                  key="createdAt"
                  width="20%"
                  render={(_: any, record: Order) => {
                    return dayjs(record.createdAt).format('MMM, DD YYYY')
                  }}
                />
                <Column
                  title={t('order_updated_label') as string}
                  dataIndex="statusChangeTime"
                  key="statusChangeTime"
                  width="20%"
                  render={(_: any, record: Order) => {
                    return dayjs(record.statusChangeTime).format('MMM, DD YYYY')
                  }}
                />
                <Column
                  title={t('order_status_label') as string}
                  dataIndex="status"
                  key="status"
                  width="10%"
                  render={(_: any, record: Order) => {
                    if (record.status === OrderStatus.Cancelled) {
                      return (
                        <span className="text-red-500">{record.status}</span>
                      )
                    }
                    return (
                      <span className="text-gray-600">{record.status}</span>
                    )
                  }}
                />
                <Column
                  title="Action"
                  key="action"
                  width="10%"
                  render={(_: any, record: Order) => {
                    if (record.status === OrderStatus.Cancelled) return
                    return (
                      <Space size="middle">
                        <a onClick={() => cancelOrder()}>Cancel</a>
                      </Space>
                    )
                  }}
                />
              </Table>
            </div>
          </div>
        </LayoutDashboard>
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

export default VendorDashboardProducts
