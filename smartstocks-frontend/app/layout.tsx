import type { Metadata } from "next"
import { Inter } from "next/font/google"
import "./globals.css"
import { Providers } from "@/lib/providers"

const inter = Inter({
  subsets: ["latin"],
  variable: "--font-inter",
})

export const metadata: Metadata = {
  title: "SmartStocks - Educacion Financiera para Jovenes",
  description: "Plataforma educativa gamificada de educacion financiera para jovenes argentinos. Aprende a invertir de forma divertida.",
  keywords: ["educacion financiera", "jovenes", "invertir", "finanzas", "argentina", "gamificacion"],
}

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode
}>) {
  return (
    <html lang="es">
      <body className={`${inter.variable} font-sans antialiased`}>
        <Providers>{children}</Providers>
      </body>
    </html>
  )
}
