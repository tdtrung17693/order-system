import type { GetServerSideProps, NextPage } from 'next'
import Link from 'next/link'
import { useTranslation } from 'next-i18next'
import { serverSideTranslations } from 'next-i18next/serverSideTranslations'
import Head from 'next/head'
import styles from '../styles/Home.module.css'
import { useContext } from 'react'
import { AuthContext } from '../context/auth.context'
import { UserRole } from '../constants/user-role'

const Home: NextPage = () => {
  const authCtx = useContext(AuthContext)
  const { t } = useTranslation('common')
  return (
    <div className={styles.container}>
      <Head>
        <title>Order Management System</title>
        <meta name="description" content="Generated by create next app" />
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <main className="flex flex-col justify-center items-center min-h-screen p-16">
        <h1 className="text-5xl mb-4">{t('title')}</h1>

        {!authCtx.authenticated ? (
          <p className="text-center space-x-4">
            <Link href="/auth/signup/vendor">
              <a
                type="button"
                className="inline-block px-6 py-2.5 bg-blue-600 text-white font-medium text-xs leading-tight uppercase rounded shadow-md hover:bg-blue-700 hover:text-white hover:shadow-lg focus:bg-blue-700 focus:shadow-lg focus:outline-none focus:ring-0 active:bg-blue-800 active:shadow-lg transition duration-150 ease-in-out"
              >
                Vendor Signup
              </a>
            </Link>
            <Link href="/auth/signup/user">
              <a
                type="button"
                className="inline-block px-6 py-2.5 bg-blue-600 text-white font-medium text-xs leading-tight uppercase rounded shadow-md hover:bg-blue-700 hover:text-white hover:shadow-lg focus:bg-blue-700 focus:shadow-lg focus:outline-none focus:ring-0 active:bg-blue-800 active:shadow-lg transition duration-150 ease-in-out"
              >
                User Signup
              </a>
            </Link>
            <Link href="/auth/signin">
              <a
                type="button"
                className="inline-block px-6 py-2.5 bg-gray-200 text-gray-800 font-medium text-xs leading-tight uppercase rounded shadow-md hover:bg-gray-300 hover:text-gray-800 hover:shadow-lg focus:bg-gray-300 focus:shadow-lg focus:outline-none focus:ring-0 active:bg-gray-400 active:shadow-lg transition duration-150 ease-in-out"
              >
                Sign in
              </a>
            </Link>
          </p>
        ) : (
          <p className="text-center space-x-4">
            {authCtx.user!.role === UserRole.Vendor && (
              <Link href="/vendors">
                <a
                  type="button"
                  className="inline-block px-6 py-2.5 bg-blue-600 text-white font-medium text-xs leading-tight uppercase rounded shadow-md hover:bg-blue-700 hover:shadow-lg focus:bg-blue-700 focus:shadow-lg focus:outline-none focus:ring-0 active:bg-blue-800 active:shadow-lg transition duration-150 ease-in-out"
                >
                  Vendor Dashboard
                </a>
              </Link>
            )}
            <Link href="/products">
              <a
                type="button"
                className="inline-block px-6 py-2.5 bg-blue-600 text-white font-medium text-xs leading-tight uppercase rounded shadow-md hover:bg-blue-700 hover:shadow-lg focus:bg-blue-700 focus:shadow-lg focus:outline-none focus:ring-0 active:bg-blue-800 active:shadow-lg transition duration-150 ease-in-out"
              >
                Browse Products
              </a>
            </Link>
          </p>
        )}
      </main>
    </div>
  )
}

export const getServerSideProps: GetServerSideProps<any> = async ({
  locale,
}) => {
  return {
    props: {
      ...(await serverSideTranslations(locale || 'en', ['common'])),
    },
  }
}

export default Home
