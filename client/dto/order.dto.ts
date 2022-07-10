import { OrderStatus } from 'constants/order'

export interface Order {
  createdAt: string
  updatedAt: string
  id: number
  status: OrderStatus
  statusChangeTime: string
  total: string
  paymentMethodId: string
  shippingAddress: string
  recipientName: string
  recipientPhone: string
  vendorId: number
  vendorName: string
}
export interface OrdersCreate {
  orders: OrderCreate[]
  paymentMethodId: string
  recipientAddress: string
  recipientName: string
  recipientPhone: string
}

export interface OrderCreate {
  items: OrderItem[]
}

export interface OrderItem {
  productId: number
  quantity: number
}
