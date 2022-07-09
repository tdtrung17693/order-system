import axios from 'axios'
import { AddCartItem, Cart, CartItem } from 'dto/cart.dto'
import { ErrorResponse } from 'dto/common'
import { http } from './http'

export const cart = {
  async getUserCart(): Promise<Cart> {
    const response = await http.get<Cart>('/carts')

    return response.data
  },
  async addCartItem(c: AddCartItem): Promise<boolean | ErrorResponse> {
    try {
      await http.post('/carts', c)
      return true
    } catch (err) {
      if (!axios.isAxiosError(err)) return false
      if (!err?.response) {
        return false
      }
      throw err.response.data as ErrorResponse
    }
  },
}
