import { OrderStatus } from 'constants/order'
import Decimal from 'decimal.js'

export interface Order {
  createdAt: string
  updatedAt: string
  id: number
  status: OrderStatus
  statusChangeTime: string
  total: string
  paymentMethodId: string
  paymentMethodName: string
  shippingAddress: string
  recipientName: string
  recipientPhone: string
  vendorId: number
  vendorName: string
  userName: string
  userId: number
  items?: {
    productName: string
    productId: string
    quantity: number
    unitPrice: number
  }[]
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
