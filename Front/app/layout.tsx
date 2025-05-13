import type React from "react"
import { Inter } from "next/font/google"
import "./globals.css"

const inter = Inter({ subsets: ["latin"] })

export const metadata = {
  title: "Bob's Corn Shop",
  description: "The best place to buy corn from Farmer Bob",
    generator: 'v0.dev'
}

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode
}>) {
  return (
    <html lang="en">
      <body className={inter.className}>
        <div className="min-h-screen bg-amber-50 flex flex-col">
          <header className="bg-green-800 text-white p-4 shadow-md">
            <div className="container mx-auto">
              <h1 className="text-xl font-bold">🌽 Bob's Farm</h1>
            </div>
          </header>

          <div className="flex-grow">{children}</div>

          <footer className="bg-green-800 text-white p-4 mt-auto">
            <div className="container mx-auto text-center">
              <p>© {new Date().getFullYear()} Bob's Farm - Fresh produce since 1985</p>
            </div>
          </footer>
        </div>
      </body>
    </html>
  )
}
