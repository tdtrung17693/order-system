import { OrderStatus } from 'constants/order'
import { Order, OrdersCreate } from 'dto/order.dto'
import { Maybe } from 'types/maybe'
import { http } from './http'

export const order = {
  isCancellableState(order: Order) {
    return (
      order.status !== OrderStatus.Cancelled &&
      order.status !== OrderStatus.Shipping &&
      order.status !== OrderStatus.Shipped
    )
  },
  isFinalState(order: Order) {
    return (
      order.status !== OrderStatus.Cancelled &&
      order.status !== OrderStatus.Shipped
    )
  },
  cancelOrder(orderId: number) {
    return http.post(`/orders/${orderId}/cancel`)
  },
  createOrders(orders: OrdersCreate) {
    return http.post('/orders', orders)
  },
  orderNextStatus(orderId: number) {
    return http.put(`/vendors/orders/${orderId}`)
  },
  getOrder(orderId: number): Promise<Order> {
    return http.get<Order>(`/orders/${orderId}`).then((res) => res.data)
  },
  async exportCsv(isVendor: boolean, status?: Maybe<OrderStatus>) {
    let endpoint = '/orders/export-csv'

    if (isVendor) {
      endpoint = `/vendors${endpoint}`
    }
    const response = await http.get(endpoint, {
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
