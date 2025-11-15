'use client'

import { useGetNotes } from '@/services/notes/notes.service'
import { CreateNoteForm } from '@/app/(protected)/dashboard/_components/create-note-form'
import { NoteCard } from '@/app/(protected)/dashboard/_components/note-card'
import type { Note } from '@/types'

export default function DashboardPage() {
  const { data: notesResponse, isLoading, error } = useGetNotes({ page: 1, limit: 50 })

  return (
    <div className="min-h-screen bg-background">
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

          {notesResponse?.data && notesResponse.data.total === 0 && !isLoading && (
            <div className="md:col-span-2 lg:col-span-3 text-center py-8">
              <p className="text-muted-foreground">
                No notes yet. Create your first note above!
              </p>
            </div>
          )}

          {notesResponse?.data?.notes?.map((note: Note) => (
            <NoteCard key={note.id} note={note} />
          ))}
        </div>
      </main>
    </div>
  )
}
