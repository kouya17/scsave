import Head from 'next/head'
import Link from 'next/link'
import Navigation from './navigation'

export const siteTitle = 'scsave'

export default function Layout({
  children,
  home
}: {
  children: React.ReactNode
  home?: boolean
}) {
  return (
    <div>
      <Head>
        <link rel="icon" href="/favicon.ico" />
        <meta
          name="description"
          content="SWR App template"
        />
        <meta name="og:title" content={siteTitle} />
        <meta name="twitter:card" content="summary" />
      </Head>
      <Navigation
        title={siteTitle}
        menu={[{name: 'ジョブ', path: '/scrape-properties-job' + (process.env.URL_POSTFIX ? process.env.URL_POSTFIX : "")}, {name: '物件', path: '/property' + (process.env.URL_POSTFIX ? process.env.URL_POSTFIX : "")}]}
      />
      <div className='container mx-auto p-4'>
        <main>{children}</main>
        {!home && (
          <div className='py-3'>
            <Link href="/">← Back to home</Link>
          </div>
        )}
      </div>
    </div>
  )
}