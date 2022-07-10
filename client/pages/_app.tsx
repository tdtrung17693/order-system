import { notification } from 'services/notification'
import { appWithTranslation, useTranslation } from 'next-i18next'
import type { AppProps } from 'next/app'
import { useCallback, useEffect, useState } from 'react'
import { cart } from 'services/cart'
import { Layout } from '../components/layout/layout'
import { AuthContext } from '../context/auth.context'
import { CartContext } from '../context/cart.context'
import { AddCartItem, CartItem } from '../dto/cart.dto'
import { User } from '../dto/user.dto'
import auth from '../services/auth'
import '../styles/globals.css'
import { Maybe } from '../types/maybe'
import { handleApiError } from 'utils/error'
import Router from 'next/router'

function MyApp({ Component, pageProps }: AppProps) {
  const { t } = useTranslation('common')
  const [authenticated, setAuthenticated] = useState(false)
  const [user, setUser] = useState<Maybe<User>>(null)
  const [cartItems, setCartItems] = useState<CartItem[]>([])
  const [checkoutItems, setCheckoutItems] = useState<CartItem[]>([])

  const getCart = useCallback(async () => {
    try {
      const userCart = await cart.getUserCart()
      const newCheckoutItems = []

      for (let i of checkoutItems) {
        const checkoutItem = userCart.items.find(
          (cartItem) => i.productId === cartItem.productId
        )
        if (checkoutItem) newCheckoutItems.push(checkoutItem)
      }
      setCartItems(userCart.items)
      setCheckoutItems(newCheckoutItems)
    } catch (error) {
      handleApiError(t, error)
    }
  }, [checkoutItems, t])

  async function addCartItem(c: AddCartItem) {
    try {
      await cart.addCartItem(c)
      await getCart()
      notification.info(t('action_success'), t('add_cart_item_success'))
    } catch (error) {
      handleApiError(t, error)
    }
  }

  async function toggleCheckoutItem(c: CartItem) {
    const idx = checkoutItems.indexOf(c)
    const newCheckoutItems = [...checkoutItems]
    if (idx >= 0) {
      newCheckoutItems.splice(idx, 1)
    } else {
      newCheckoutItems.push(c)
    }

    setCheckoutItems(newCheckoutItems)
  }

  function logout() {
    auth.logout()
    setAuthenticated(false)
    setUser(null)
  }

  useEffect(() => {
    auth.onLogin(() => {
      setUser(auth.user)
      setAuthenticated(true)
      getCart()
      Router.push('/')
    })
    auth.init().then((user) => {
      if (user) {
        setUser(user)
        setAuthenticated(true)
        getCart()
      }
    })
  }, [getCart])

  return (
    <AuthContext.Provider value={{ authenticated, user, logout }}>
      <CartContext.Provider
        value={{
          cartItems,
          addCartItem,
          checkoutItems,
          toggleCheckoutItem,
          refreshCart: getCart,
        }}
      >
        <Layout>
          <Component {...pageProps} />
        </Layout>
      </CartContext.Provider>
    </AuthContext.Provider>
  )
}

export default appWithTranslation(MyApp)
