'use client'

import { useState } from 'react'
import { useCreateNote } from '@/services/notes/notes.service'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { ImageIcon } from 'lucide-react'

interface CreateNoteFormProps {
  onSuccess?: () => void
}

export function CreateNoteForm({ onSuccess }: CreateNoteFormProps) {
  const [title, setTitle] = useState('')
  const [content, setContent] = useState('')
  const [imageFile, setImageFile] = useState<File | null>(null)

  const createNoteMutation = useCreateNote()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await createNoteMutation.mutateAsync({ title, content, image: imageFile || undefined })
      setTitle('')
      setContent('')
      setImageFile(null)
      onSuccess?.()
    } catch (err) {
      console.error('Failed to create note:', err)
    }
  }

  const handleImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      setImageFile(e.target.files[0])
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
            <Textarea
              id="content"
              placeholder="Enter note content"
              value={content}
              onChange={(e) => setContent(e.target.value)}
              required
              disabled={createNoteMutation.isPending}
              className="min-h-[120px]"
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="image">Image (optional)</Label>
            <input
              type="file"
              accept="image/*"
              onChange={handleImageChange}
              disabled={createNoteMutation.isPending}
              className="hidden"
              id="image-upload"
            />
            <label htmlFor="image-upload">
              <Button
                type="button"
                variant="neutral"
                className="w-full"
                disabled={createNoteMutation.isPending}
                asChild
              >
                <span>
                  <ImageIcon className="mr-2 h-4 w-4" />
                  {imageFile ? imageFile.name : 'Choose Image'}
                </span>
              </Button>
            </label>
          </div>
          <Button type="submit" disabled={createNoteMutation.isPending} className="w-full">
            {createNoteMutation.isPending ? 'Creating...' : 'Create Note'}
          </Button>
        </form>
      </CardContent>
    </Card>
  )
}
