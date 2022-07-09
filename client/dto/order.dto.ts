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
