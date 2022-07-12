import {
  Book24Filled,
  Box24Filled,
  Cart24Regular,
  DoorArrowLeft24Regular,
  Home24Regular,
  PersonAccounts24Filled,
} from '@fluentui/react-icons'
import { Tooltip } from 'antd'
import { UserRole } from 'constants/user-role'
import { useTranslation } from 'next-i18next'
import Link from 'next/link'
import React, { useContext } from 'react'
import { AuthContext } from '../../context/auth.context'
import { CartContext } from '../../context/cart.context'

type Props = {
  children?: React.ReactNode
}

export const Layout: React.FC<Props> = (props) => {
  const { t } = useTranslation('common')
  const authCtx = useContext(AuthContext)
  const cartCtx = useContext(CartContext)
  return (
    <div className="p-4">
      <nav className="flex justify-center items-center space-x-4">
        <Link className="" href="/">
          <a>
            <Home24Regular />
          </a>
        </Link>
        {authCtx.authenticated && (
          <>
            {authCtx.user?.role == UserRole.Vendor && (
              <Link href="/vendors">
                <Tooltip title={t('dashboard_title')}>
                  <a>
                    <Book24Filled />
                  </a>
                </Tooltip>
              </Link>
            )}

            <Link className="" href="/products">
              <Tooltip title={t('browse_products_tooltip')}>
                <a>
                  <Box24Filled />
                </a>
              </Tooltip>
            </Link>
            <Link href="/orders">
              <Tooltip title={t('manage_my_order_title')}>
                <a>
                  <PersonAccounts24Filled />
                </a>
              </Tooltip>
            </Link>
            <Tooltip title={t('logout_title')}>
              <a
                title="logout"
                onClick={() => authCtx.logout()}
                className="flex px-4 py-2 justify-center items-center rounded-md bg-gray-100"
              >
                <span className="flex mr-4">
                  <DoorArrowLeft24Regular />{' '}
                </span>
                <span>
                  {authCtx.user?.email} ({authCtx.user?.name})
                </span>
              </a>
            </Tooltip>
            <Link className="" href="/cart">
              <Tooltip title={t('cart_title')}>
                <a className="relative inline-flex justify-between items-center py-2 px-4 rounded-md text-white bg-red-400 hover:bg-red-500 hover:text-white transition ">
                  <Cart24Regular />
                  <span className="text-sm ml-4">
                    {cartCtx.cartItems.reduce((sum, i) => sum + i.quantity, 0)}
                  </span>
                </a>
              </Tooltip>
            </Link>
          </>
        )}
      </nav>
      {props.children}
    </div>
  )
}
