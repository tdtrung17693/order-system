import { CartItem } from 'dto/cart.dto'

type VendorId = number
export function groupCartItemByVendor(
  cartItems: CartItem[]
): Record<
  VendorId,
  { vendorId: number; vendorName: string; items: CartItem[] }
> {
  const vendorItems: Record<
    VendorId,
    { vendorId: number; vendorName: string; items: CartItem[] }
  > = {}
  cartItems.forEach((item) => {
    if (!vendorItems[item.vendorId]) {
      vendorItems[item.vendorId] = {
        vendorId: item.vendorId,
        vendorName: item.vendorName,
        items: [],
      }
    }

    vendorItems[item.vendorId].items.push(item)
  })
  return vendorItems
}
