'use client'

import { useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { LogOut } from 'lucide-react'
import { useAuth } from '@/contexts/AuthContext'
import { useNotes } from '@/hooks/use-notes'
import { Button } from '@/components/ui/button'
import { CreateNoteForm } from '@/components/create-note-form'
import { NoteCard } from '@/components/note-card'
import { ThemeToggle } from '@/components/theme-toggle'

export default function DashboardPage() {
  const router = useRouter()
  const { isAuthenticated, user, logout } = useAuth()
  const { data: notes, isLoading, error } = useNotes()

  useEffect(() => {
    if (!isAuthenticated) {
      router.push('/login')
    }
  }, [isAuthenticated, router])

  const handleLogout = () => {
    logout()
    router.push('/login')
  }

  if (!isAuthenticated) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <p>Redirecting...</p>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-background">
      <header className="border-b">
        <div className="container mx-auto flex items-center justify-between py-4 px-4">
          <div>
            <h1 className="text-2xl font-bold">My Notes</h1>
            <p className="text-sm text-muted-foreground">{user?.email}</p>
          </div>
          <div className="flex items-center space-x-2">
            <ThemeToggle />
            <Button onClick={handleLogout} variant="outline">
              <LogOut className="mr-2 h-4 w-4" />
              Logout
            </Button>
          </div>
        </div>
      </header>

      <main className="container mx-auto py-8 px-4">
        <div className="grid gap-8 md:grid-cols-2 lg:grid-cols-3">
          <div className="md:col-span-2 lg:col-span-3">
            <CreateNoteForm />
          </div>

          {isLoading && (
            <div className="md:col-span-2 lg:col-span-3 text-center py-8">
              <p className="text-muted-foreground">Loading notes...</p>
            </div>
          )}

          {error && (
            <div className="md:col-span-2 lg:col-span-3">
              <div className="rounded-lg bg-destructive/15 p-4 text-destructive">
                Error loading notes: {error.message}
              </div>
            </div>
          )}

          {notes && notes.length === 0 && !isLoading && (
            <div className="md:col-span-2 lg:col-span-3 text-center py-8">
              <p className="text-muted-foreground">
                No notes yet. Create your first note above!
              </p>
            </div>
          )}

          {notes?.map((note: any) => (
            <NoteCard key={note.id} note={note} />
          ))}
        </div>
      </main>
    </div>
  )
}
