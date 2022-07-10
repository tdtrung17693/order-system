import React from 'react'
import { AddCartItem, CartItem } from '../dto/cart.dto'

export const CartContext = React.createContext<{
  cartItems: CartItem[]
  checkoutItems: CartItem[]
  addCartItem: (c: AddCartItem) => void
  toggleCheckoutItem: (c: CartItem) => void
  refreshCart: () => void
}>({
  cartItems: [],
  checkoutItems: [],
  addCartItem: (c) => {},
  toggleCheckoutItem: (c) => {},
  refreshCart: () => {},
})
