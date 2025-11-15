'use client'

import { useState } from 'react'
import { Trash2, Image as ImageIcon, Edit } from 'lucide-react'
import { Note } from '@/types'
import { useDeleteNote, useUpdateNote, useGetNoteById } from '@/services/notes/notes.service'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
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
import Image from 'next/image'

interface NoteCardProps {
  note: Note
}

export function NoteCard({ note }: NoteCardProps) {
  const [isEditing, setIsEditing] = useState(false)
  const [title, setTitle] = useState(note.title)
  const [content, setContent] = useState(note.content)
  const [imageFile, setImageFile] = useState<File | null>(null)

  const deleteNoteMutation = useDeleteNote()
  const updateNoteMutation = useUpdateNote()
  const { data: noteData } = useGetNoteById(note.id, isEditing)

  const isLoading = deleteNoteMutation.isPending || updateNoteMutation.isPending

  if (noteData && !isEditing) {
    if (noteData.title !== title) setTitle(noteData.title)
    if (noteData.content !== content) setContent(noteData.content)
  }

  const handleDelete = async () => {
    try {
      await deleteNoteMutation.mutateAsync(note.id)
    } catch (err) {
      console.error('Failed to delete note:', err)
    }
  }

  const handleImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      setImageFile(e.target.files[0])
    }
  }

  const handleUpdate = async () => {
    try {
      await updateNoteMutation.mutateAsync({
        id: note.id,
        data: { title, content, image: imageFile || undefined }
      })
      setIsEditing(false)
      setImageFile(null)
    } catch (err) {
      console.error('Failed to update note:', err)
    }
  }

  const handleCancelEdit = () => {
    setTitle(note.title)
    setContent(note.content)
    setImageFile(null)
    setIsEditing(false)
  }

  const imageUrl = note.image_path
    ? `${process.env.NEXT_PUBLIC_API_URL}/${note.image_path.replace(/\\/g, '/')}`
    : null

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center justify-between">
          <span className="truncate">{isEditing ? 'Edit Note' : note.title}</span>
          <div className="flex gap-2">
            {!isEditing && (
              <>
                <Button
                  variant="default"
                  size="icon"
                  onClick={() => setIsEditing(true)}
                  disabled={isLoading}
                >
                  <Edit className="h-4 w-4" />
                </Button>
                <AlertDialog>
                  <AlertDialogTrigger asChild>
                    <Button
                      variant="default"
                      size="icon"
                      disabled={isLoading}
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </AlertDialogTrigger>
                  <AlertDialogContent>
                    <AlertDialogHeader>
                      <AlertDialogTitle>Are you sure?</AlertDialogTitle>
                      <AlertDialogDescription>
                        This action cannot be undone. This will permanently delete your note.
                      </AlertDialogDescription>
                    </AlertDialogHeader>
                    <AlertDialogFooter>
                      <AlertDialogCancel>Cancel</AlertDialogCancel>
                      <AlertDialogAction onClick={handleDelete}>
                        Delete
                      </AlertDialogAction>
                    </AlertDialogFooter>
                  </AlertDialogContent>
                </AlertDialog>
              </>
            )}
          </div>
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        {isEditing ? (
          <>
            <div className="space-y-2">
              <Label htmlFor={`title-${note.id}`}>Title</Label>
              <Input
                id={`title-${note.id}`}
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                disabled={isLoading}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor={`content-${note.id}`}>Content</Label>
              <Textarea
                id={`content-${note.id}`}
                value={content}
                onChange={(e) => setContent(e.target.value)}
                disabled={isLoading}
                className="min-h-[120px]"
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor={`image-${note.id}`}>Image (optional)</Label>
              <input
                type="file"
                accept="image/*"
                onChange={handleImageChange}
                disabled={isLoading}
                className="hidden"
                id={`image-${note.id}`}
              />
              <label htmlFor={`image-${note.id}`}>
                <Button
                  type="button"
                  variant="neutral"
                  className="w-full"
                  disabled={isLoading}
                  asChild
                >
                  <span>
                    <ImageIcon className="mr-2 h-4 w-4" />
                    {imageFile ? imageFile.name : 'Choose New Image'}
                  </span>
                </Button>
              </label>
            </div>
          </>
        ) : (
          <>
            <p className="text-sm text-muted-foreground whitespace-pre-wrap">
              {note.content}
            </p>
            {note.image_path && (
              <div className="relative aspect-video w-full overflow-hidden rounded-lg border">
                <Image
                  src={imageUrl!}
                  alt={note.title}
                  fill
                  className="object-cover"
                  loading='eager'
                  unoptimized
                />
              </div>
            )}
          </>
        )}
      </CardContent>
      <CardFooter className="flex flex-col space-y-2">
        {isEditing ? (
          <div className="flex w-full gap-2">
            <Button
              variant="neutral"
              onClick={handleCancelEdit}
              disabled={isLoading}
              className="flex-1"
            >
              Cancel
            </Button>
            <Button
              onClick={handleUpdate}
              disabled={isLoading}
              className="flex-1"
            >
              {isLoading ? 'Saving...' : 'Save'}
            </Button>
          </div>
        ) : (
          <p className="text-xs text-muted-foreground">
            {new Date(note.created_at).toLocaleDateString()}
          </p>
        )}
      </CardFooter>
    </Card>
  )
}
