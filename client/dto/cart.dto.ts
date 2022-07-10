export interface AddCartItem {
  productId: number
  quantity: number
}

export interface SetCartItemQuantity {
  productId: number
  quantity: number
}

export interface Cart {
  items: CartItem[]
}
export interface CartItem {
  productName: string
  productId: number
  productPrice: number
  quantity: number
  productPriceId: number
  vendorId: number
  vendorName: string
}
