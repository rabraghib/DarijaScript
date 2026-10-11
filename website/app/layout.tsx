import type { Metadata } from 'next';
import { Inter } from 'next/font/google';
import './globals.css';
import Link from 'next/link';

const inter = Inter({ subsets: ['latin'] });

export const metadata: Metadata = {
  title: 'DarijaScript Programming Language',
  description: 'The greatest language that ever was or will be!',
  openGraph: {
    type: 'website',
    siteName: 'DarijaScript',
    url: 'https://ds.rabraghib.com/',
    title: 'DarijaScript Programming Language',
    description: 'The greatest language that ever was or will be!',
    images: [
      {
        url: 'https://ds.rabraghib.com/meme.webp',
        width: 1280,
        height: 720,
        alt: 'DarijaScript Programming Language',
      },
    ],
  },
  twitter: {
    site: '@rabraghib',
    card: 'summary',
    title: 'DarijaScript Programming Language',
    description: 'The greatest language that ever was or will be!',
    images: [
      {
        url: 'https://ds.rabraghib.com/meme.webp',
        width: 1280,
        height: 720,
        alt: 'DarijaScript Programming Language',
      },
    ],
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={inter.className}>
        <div className="text-white bg-black">
          <Header />
          <div className="min-h-screen pt-16">{children}</div>
          <Footer />
        </div>
      </body>
    </html>
  );
}

function Header() {
  return (
    <nav className="fixed h-16 overflow-hidden z-20 w-full bg-slate-900 bg-opacity-20 border-b border-slate-700 backdrop-blur-2xl">
      <div className="px-6 m-auto max-w-6xl">
        <div className="flex flex-wrap items-center justify-between py-3">
          <div className="items-center flex justify-between w-auto">
            <Link href="/">
              {/* eslint-disable-next-line @next/next/no-img-element */}
              <img
                src="/logo.svg"
                alt="DarijaScript Logo"
                style={{ height: '2.5rem' }}
              />
            </Link>
          </div>
          <div className="w-fit flex-nowrap justify-end items-center flex h-fit">
            <div className="tracking-wide text-sm text-gray-300 pr-2">
              <a
                target="_blank"
                href="https://github.com/rabraghib/darijascript"
                className="block px-4 transition hover:text-indigo-400"
              >
                <span>GitHub</span>
              </a>
            </div>

            <Link
              href="/playground"
              className="text-sm font-medium inline-flex items-center justify-center rounded px-4 py-2 shadow-lg bg-blue-600 hover:bg-blue-700 text-white"
            >
              Playground!
            </Link>
          </div>
        </div>
      </div>
    </nav>
  );
}

function Footer() {
  return (
    <div className="bg-slate-900 bg-opacity-30 text-lg text-center p-4">
      Made with <span className="text-red-600">&hearts;</span> by{' '}
      <a
        target="_blank"
        href="https://rabraghib.com/"
        className="underline whitespace-nowrap"
      >
        Rabyâ Raghib
      </a>
    </div>
  );
}
