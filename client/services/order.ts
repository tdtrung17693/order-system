import { OrderStatus } from 'constants/order'
import { OrdersCreate } from 'dto/order.dto'
import { Maybe } from 'types/maybe'
import { http } from './http'

export const order = {
  cancelOrder(orderId: number) {
    return http.post(`/orders/${orderId}/cancel`)
  },
  createOrders(orders: OrdersCreate) {
    return http.post('/orders', orders)
  },
  orderNextStatus(orderId: number) {
    return http.put(`/vendors/orders/${orderId}`)
  },
  async exportCsv(status?: Maybe<OrderStatus>) {
    const response = await http.get('/orders/export-csv', {
      responseType: 'blob',
      params: {
        status,
      },
    })

    const url_2 = window.URL.createObjectURL(new Blob([response.data]))
    const link = document.createElement('a')
    link.href = url_2
    link.setAttribute('download', `orders${status ? `-${status}` : ''}.csv`) //or any other extension
    document.body.appendChild(link)
    link.click()
    setTimeout(function () {
      document.body.removeChild(link)
      window.URL.revokeObjectURL(url_2)
    }, 200)
  },
}
