import React from 'react'
import { AddCartItem, CartItem } from '../dto/cart.dto'

export const CartContext = React.createContext<{
  cartItems: CartItem[]
  addCartItem: (c: AddCartItem) => void
}>({
  cartItems: [],
  addCartItem: (c) => {},
})
