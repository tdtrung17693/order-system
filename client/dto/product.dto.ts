import { ProductUpdateType } from 'constants/product'

export interface Product {
  id: number
  name: string
  description: string
  vendorId: number
  unit: string
  stockQuantity: number
  productPrice: number
}

export interface CreateProduct {
  name: string
  descriptiton: string
}

export interface UpdateProductStock {
  productId: number
  description: string
  type: ProductUpdateType
  quantity: number
}

export interface SetProductPrice {
  productId: number
  price: number
}
