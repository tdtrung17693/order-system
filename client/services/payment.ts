import { PaymentMethod } from 'dto/payment.dto'
import { http } from './http'

export const payment = {
  getSupportedMethods() {
    return http
      .get<PaymentMethod[]>('/payment-methods')
      .then((response) => response.data)
  },
}
