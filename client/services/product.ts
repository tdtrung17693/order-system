import {
  buildPaginationRequest,
  PaginationQuery,
  PaginationResponse,
} from 'dto/pagination.dto'
import {
  CreateProduct,
  Product,
  SetProductPrice,
  UpdateProductStock,
} from 'dto/product.dto'
import { http } from './http'

export const product = {
  updateProductStock(data: UpdateProductStock) {
    return http.post(`/vendors/products/${data.productId}/stocks`, data)
  },
  setProductPrice(data: SetProductPrice) {
    return http.post(`/vendors/products/${data.productId}/prices`, data)
  },
  createProduct(data: CreateProduct) {
    return http.post(`/vendors/products`, data)
  },
  getProducts(data: PaginationQuery): Promise<PaginationResponse<Product>> {
    return http
      .get<PaginationResponse<Product>>('/products', {
        params: data,
      })
      .then(({ data }) => {
        return data
      })
  },
}
