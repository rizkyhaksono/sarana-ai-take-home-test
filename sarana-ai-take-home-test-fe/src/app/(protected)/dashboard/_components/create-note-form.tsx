'use client'

import { useState } from 'react'
import { useCreateNote } from '@/services/notes/notes.service'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'

interface CreateNoteFormProps {
  onSuccess?: () => void
}

export function CreateNoteForm({ onSuccess }: CreateNoteFormProps) {
  const [title, setTitle] = useState('')
  const [content, setContent] = useState('')

  const createNoteMutation = useCreateNote()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await createNoteMutation.mutateAsync({ title, content })
      setTitle('')
      setContent('')
      onSuccess?.()
    } catch (err) {
      console.error('Failed to create note:', err)
    }
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Create New Note</CardTitle>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit} className="space-y-4">
          {createNoteMutation.error && (
            <div className="rounded-lg bg-destructive/15 p-3 text-sm text-destructive">
              {createNoteMutation.error.message}
            </div>
          )}
          <div className="space-y-2">
            <Label htmlFor="title">Title</Label>
            <Input
              id="title"
              placeholder="Enter note title"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              required
              disabled={createNoteMutation.isPending}
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="content">Content</Label>
            <textarea
              id="content"
              placeholder="Enter note content"
              value={content}
              onChange={(e) => setContent(e.target.value)}
              required
              disabled={createNoteMutation.isPending}
              className="flex min-h-[120px] w-full rounded-md border border-input bg-background px-3 py-2 text-base ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 md:text-sm"
            />
          </div>
          <Button type="submit" disabled={createNoteMutation.isPending} className="w-full">
            {createNoteMutation.isPending ? 'Creating...' : 'Create Note'}
          </Button>
        </form>
      </CardContent>
    </Card>
  )
}
