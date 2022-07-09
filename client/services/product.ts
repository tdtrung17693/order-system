import { SetProductPrice, UpdateProductStock } from 'dto/product.dto'
import { http } from './http'

export const product = {
  updateProductStock(data: UpdateProductStock) {
    return http.post(`/vendors/products/${data.productId}/stocks`, data)
  },
  setProductPrice(data: SetProductPrice) {
    return http.post(`/vendors/products/${data.productId}/prices`, data)
  },
}
