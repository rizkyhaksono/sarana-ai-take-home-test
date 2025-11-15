import Navbar from "./navbar"
import Footer from "./footer"
import { QueryProvider } from "@/providers/query-provider"

export default function BaseLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <div className="flex min-h-screen flex-col">
      <QueryProvider>
        <Navbar />
        {children}
        <Footer />
      </QueryProvider>
    </div>
  )
}