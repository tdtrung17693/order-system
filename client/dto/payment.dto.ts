export interface PaymentMethod {
  id: string
  name: string
}

export interface PaymentInfo {
  paymentMethodId: string
  recipientPhone: string
  recipientName: string
  recipientAddress: string
}
