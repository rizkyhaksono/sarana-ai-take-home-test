"use client"

import { ThemeToggle } from "./theme-toggle";
import { LogOut } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { useRouter } from 'next/navigation'
import { useLogout, useAuth } from '@/services/auth/auth.service'

export default function Navbar() {
  const router = useRouter()
  const { mutateAsync: logout } = useLogout()
  const { data: user } = useAuth()

  const handleLogout = async () => {
    await logout()
    router.push('/login')
  }

  return (
    <header className="border-b">
      <div>
        {!user ? (
          <div className="container mx-auto flex items-center justify-between py-4 px-4">
            <h1 className="text-base font-bold">Sarana Notes App</h1>
            <ThemeToggle />
          </div>
        ) : (
          <div className="container mx-auto flex items-center justify-between py-4 px-4">
            <div>
              <h1 className="text-2xl font-bold">My Notes</h1>
              <p className="text-sm text-muted-foreground">{user?.email}</p>
            </div>
            <div className="flex items-center space-x-2">
              <ThemeToggle />
              <Button onClick={handleLogout} variant="default">
                <LogOut className="mr-2 h-4 w-4" />
                Logout
              </Button>
            </div>
          </div>
        )}
      </div>
    </header>
  )
}