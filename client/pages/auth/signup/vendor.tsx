import type {
  GetServerSideProps,
  GetStaticPaths,
  GetStaticProps,
  NextPage,
} from 'next'
import { useTranslation } from 'next-i18next'
import { serverSideTranslations } from 'next-i18next/serverSideTranslations'
import Head from 'next/head'
import Image from 'next/image'
import { SignUpForm } from '../../../components/signup-form/signup-form'
import { UserRole } from '../../../constants/user-role'

const Signup: NextPage = () => {
  const { t } = useTranslation('common')
  return (
    <div className="p-4">
      <Head>
        <title>Order Management System | Sign Up</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <main className="flex flex-col justify-center items-center min-h-screen p-16">
        <h1 className="text-5xl mb-4">{t('signup_vendor')}</h1>
        <div className="flex justify-center">
          <SignUpForm role={UserRole.Vendor} />
        </div>
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
      // Will be passed to the page component as props
    },
  }
}

export default Signup
