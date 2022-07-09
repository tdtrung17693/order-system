import { notification } from 'services/notification'
import { appWithTranslation, useTranslation } from 'next-i18next'
import type { AppProps } from 'next/app'
import { useEffect, useState } from 'react'
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

function MyApp({ Component, pageProps }: AppProps) {
  const { t } = useTranslation('common')
  const [authenticated, setAuthenticated] = useState(false)
  const [user, setUser] = useState<Maybe<User>>(null)
  const [cartItems, setCartItems] = useState<CartItem[]>([])

  useEffect(() => {
    auth.init().then((user) => {
      if (user) {
        setUser(user)
        setAuthenticated(true)
        getCart()
      }
    })
  })

  async function getCart() {
    try {
      const userCart = cart.getUserCart()
      setCartItems((await userCart).items)
    } catch (error) {
      handleApiError(t, error)
    }
  }

  async function addCartItem(c: AddCartItem) {
    try {
      await cart.addCartItem(c)
      await getCart()
      notification.info(t('action_success'), t('add_cart_item_success'))
    } catch (error) {
      handleApiError(t, error)
    }
  }

  function logout() {
    auth.logout()
    setAuthenticated(false)
    setUser(null)
  }

  return (
    <AuthContext.Provider value={{ authenticated, user, logout }}>
      <CartContext.Provider value={{ cartItems, addCartItem }}>
        <Layout>
          <Component {...pageProps} />
        </Layout>
      </CartContext.Provider>
    </AuthContext.Provider>
  )
}

export default appWithTranslation(MyApp)
