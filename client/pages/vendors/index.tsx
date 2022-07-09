import type { GetServerSideProps, NextPage } from 'next'
import Error from 'next/error'
import { useTranslation } from 'next-i18next'
import { serverSideTranslations } from 'next-i18next/serverSideTranslations'
import { useContext, useEffect } from 'react'
import { AuthContext } from 'context/auth.context'
import { UserRole } from 'constants/user-role'
import { useRouter } from 'next/router'
import auth from 'services/auth'
import { LayoutDashboard } from 'components/layout/layout-dashboard'

const VendorDashboard: NextPage = () => {
  const authCtx = useContext(AuthContext)
  const router = useRouter()
  const { t } = useTranslation('common')

  useEffect(() => {
    if (auth.initialized && !authCtx.user) {
      router.push('/auth/signin')
    }
  }, [authCtx.user, router])

  if (authCtx.user && authCtx.user.role != UserRole.Vendor) {
    return <Error statusCode={403} />
  }

  return (
    <div className="px-20">
      <main className="flex flex-col justify-start items-center min-h-screen p-16">
        <LayoutDashboard>
          <div className="text-center mt-4 text-3xl">
            Hello, {authCtx.user?.name}
          </div>
        </LayoutDashboard>
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

export default VendorDashboard
