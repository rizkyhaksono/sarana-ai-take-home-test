"use client"

import { ThemeToggle } from "./theme-toggle";
import { LogOut, Book } from 'lucide-react'
import { Button } from '@/components/ui/button'
import Link from 'next/link'
import { useRouter, usePathname } from 'next/navigation'
import { useLogout, useAuth } from '@/services/auth/auth.service'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '@/components/ui/alert-dialog'

export default function Navbar() {
  const router = useRouter()
  const pathname = usePathname()
  const { mutateAsync: logout } = useLogout()
  const { data: user } = useAuth()

  const handleLogout = async () => {
    await logout()
    router.push('/login')
  }

  const isAuthPage = pathname?.startsWith('/login') || pathname?.startsWith('/register')

  return (
    <header className="sticky top-0 z-50 w-full border-b-2 border-border bg-background backdrop-blur">
      <div className="container mx-auto">
        {!user || isAuthPage ? (
          <div className="flex items-center justify-between py-4 px-4">
            <Link href="/" className="flex items-center space-x-2">
              <Book className="h-6 w-6" />
              <h1 className="text-xl font-heading">Sarana Notes</h1>
            </Link>
            <ThemeToggle />
          </div>
        ) : (
          <div className="flex items-center justify-between py-4 px-4">
            <Link href="/dashboard" className="flex items-center space-x-3">
              <Book className="h-6 w-6" />
              <div>
                <h1 className="text-2xl font-heading">My Notes</h1>
                <p className="text-sm text-muted-foreground">{user?.user?.email}</p>
              </div>
            </Link>
            <div className="flex items-center gap-3">
              <nav className="hidden md:flex items-center gap-2">
                <Link href="/dashboard">
                  <Button variant={pathname === '/dashboard' ? 'default' : 'neutral'} size="sm">
                    Dashboard
                  </Button>
                </Link>
                <Link href="/logs">
                  <Button variant={pathname === '/logs' ? 'default' : 'neutral'} size="sm">
                    Logs
                  </Button>
                </Link>
              </nav>
              <ThemeToggle />
              <AlertDialog>
                <AlertDialogTrigger asChild>
                  <Button variant="default">
                    <LogOut className="mr-2 h-4 w-4" />
                    Logout
                  </Button>
                </AlertDialogTrigger>
                <AlertDialogContent>
                  <AlertDialogHeader>
                    <AlertDialogTitle>Are you sure you want to logout?</AlertDialogTitle>
                    <AlertDialogDescription>
                      You will be redirected to the login page and need to sign in again to access your notes.
                    </AlertDialogDescription>
                  </AlertDialogHeader>
                  <AlertDialogFooter>
                    <AlertDialogCancel>Cancel</AlertDialogCancel>
                    <AlertDialogAction onClick={handleLogout}>
                      Logout
                    </AlertDialogAction>
                  </AlertDialogFooter>
                </AlertDialogContent>
              </AlertDialog>
            </div>
          </div>
        )}
      </div>
    </header>
  )
}