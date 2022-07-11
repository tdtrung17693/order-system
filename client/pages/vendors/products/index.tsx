import { Button, Modal, Popconfirm, Space, Table } from 'antd'
import { TablePaginationConfig } from 'antd/lib/table'
import Column from 'antd/lib/table/Column'
import { FilterValue, SorterResult } from 'antd/lib/table/interface'
import { LayoutDashboard } from 'components/layout/layout-dashboard'
import { ProductCreateForm } from 'components/product-create-form/product-create-form'
import { ProductSetPriceForm } from 'components/product-set-price-form/product-set-price-form'
import { ProductUpdateForm } from 'components/product-update-form/product-update-form'
import { ItemsPerPage } from 'constants/pagination'
import { ProductUpdateType } from 'constants/product'
import { UserRole } from 'constants/user-role'
import { AuthContext } from 'context/auth.context'
import { buildPaginationRequest, PaginationResponse } from 'dto/pagination.dto'
import { Product } from 'dto/product.dto'
import type { GetServerSideProps, NextPage } from 'next'
import { useTranslation } from 'next-i18next'
import { serverSideTranslations } from 'next-i18next/serverSideTranslations'
import Error from 'next/error'
import { useRouter } from 'next/router'
import { useCallback, useContext, useEffect, useState } from 'react'
import auth from 'services/auth'
import { http } from 'services/http'
import { product } from 'services/product'
import { Maybe } from 'types/maybe'

enum ProductModalType {
  ModalUpdateStock,
  ModalSetPrice,
  ModalNewProduct,
}

const VendorDashboardProducts: NextPage = () => {
  const authCtx = useContext(AuthContext)
  const router = useRouter()
  const { t } = useTranslation('common')
  const [products, setProducts] = useState<Product[]>([])
  const [isModalVisible, setIsModalVisible] = useState(false)
  const [currentModalProduct, setCurrentModalProduct] =
    useState<Maybe<Product>>(null)
  const [modalType, setModalType] = useState(ProductModalType.ModalSetPrice)
  const [productUpdateType, setProductUpdateType] = useState(
    ProductUpdateType.Export
  )
  const [{ current, pageSize }, setPagination] = useState({
    current: 1,
    pageSize: ItemsPerPage,
  })
  const [total, setTotal] = useState(-1)

  const fetchData = useCallback((pageIndex: number) => {
    return http
      .get<PaginationResponse<Product>>('/vendors/products', {
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

  const handleTableChange = (newPagination: TablePaginationConfig) => {
    fetchData((newPagination.current || 1) - 1)
  }

  if (authCtx.user && authCtx.user.role != UserRole.Vendor) {
    return <Error statusCode={403} />
  }

  const handleCancel = () => {
    setIsModalVisible(false)
  }

  const showProductUpdateModal = (
    product: Product,
    updateType: ProductUpdateType
  ) => {
    setIsModalVisible(true)
    setCurrentModalProduct(product)
    setProductUpdateType(updateType)
    setModalType(ProductModalType.ModalUpdateStock)
  }

  const showSetPriceModal = (product: Product) => {
    setIsModalVisible(true)
    setCurrentModalProduct(product)
    setModalType(ProductModalType.ModalSetPrice)
  }

  const showNewProductModal = () => {
    setIsModalVisible(true)
    setModalType(ProductModalType.ModalNewProduct)
  }

  return (
    <>
      <div className="px-20">
        <main className="flex flex-col justify-start items-center min-h-screen p-16">
          <LayoutDashboard>
            <div className=" mt-4 text-base">
              <Button onClick={showNewProductModal}>New product</Button>
              <div className="mt-4">
                <Table
                  className="min-w-full"
                  dataSource={products}
                  pagination={{ pageSize, total }}
                  onChange={handleTableChange as any}
                >
                  <Column
                    title={t('product_id_label') as string}
                    dataIndex="id"
                    key="id"
                    width="5%"
                  />
                  <Column
                    title={t('product_name_label') as string}
                    dataIndex="name"
                    key="name"
                    width="60%"
                  />
                  <Column
                    title={t('product_price_label') as string}
                    dataIndex="productPrice"
                    key="productPrice"
                    width="10%"
                  />
                  <Column
                    title={t('product_stock_quantity_label') as string}
                    dataIndex="stockQuantity"
                    key="stockQuantity"
                    width="10%"
                  />
                  <Column
                    title="Action"
                    key="action"
                    render={(_: any, record: Product) => (
                      <Space size="middle">
                        <Button
                          type="default"
                          onClick={() =>
                            showProductUpdateModal(
                              record,
                              ProductUpdateType.Import
                            )
                          }
                        >
                          Import stock
                        </Button>
                        <Button
                          type="default"
                          onClick={() =>
                            showProductUpdateModal(
                              record,
                              ProductUpdateType.Export
                            )
                          }
                        >
                          Export stock
                        </Button>
                        <Button
                          type="default"
                          onClick={() => showSetPriceModal(record)}
                        >
                          Set price
                        </Button>
                      </Space>
                    )}
                  />
                </Table>
              </div>
            </div>
          </LayoutDashboard>
        </main>
      </div>
      <Modal
        title={
          modalType == ProductModalType.ModalUpdateStock
            ? productUpdateType == ProductUpdateType.Import
              ? t('product_import_modal_title')
              : t('product_export_modal_title')
            : modalType == ProductModalType.ModalNewProduct
            ? t('product_create_modal_title')
            : t('product_set_price_modal_title')
        }
        visible={isModalVisible}
        footer={null}
        onCancel={handleCancel}
      >
        {modalType === ProductModalType.ModalUpdateStock ? (
          <ProductUpdateForm
            productId={currentModalProduct?.id || -1}
            updateType={productUpdateType}
            onUpdated={() => {
              setIsModalVisible(false)
              fetchData((current || 1) - 1)
            }}
          />
        ) : modalType == ProductModalType.ModalNewProduct ? (
          <ProductCreateForm
            onCreated={() => {
              setIsModalVisible(false)
              fetchData((current || 1) - 1)
            }}
          />
        ) : (
          <ProductSetPriceForm
            productId={currentModalProduct?.id || -1}
            currentPrice={currentModalProduct?.productPrice || 0}
            onUpdated={() => {
              setIsModalVisible(false)
              fetchData((current || 1) - 1)
            }}
          />
        )}
      </Modal>
    </>
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
