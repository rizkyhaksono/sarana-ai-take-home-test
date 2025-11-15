'use client'

import { useState } from 'react'
import { Trash2, Upload, Image as ImageIcon } from 'lucide-react'
import { Note } from '@/types'
import { useDeleteNote, useUploadNoteImage } from '@/services/notes/notes.service'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import Image from 'next/image'

interface NoteCardProps {
  note: Note
}

export function NoteCard({ note }: NoteCardProps) {
  const [imageFile, setImageFile] = useState<File | null>(null)
  const deleteNoteMutation = useDeleteNote()
  const uploadImageMutation = useUploadNoteImage()

  const isLoading = deleteNoteMutation.isPending || uploadImageMutation.isPending

  const handleDelete = async () => {
    if (confirm('Are you sure you want to delete this note?')) {
      try {
        await deleteNoteMutation.mutateAsync(note.id)
      } catch (err) {
        console.error('Failed to delete note:', err)
      }
    }
  }

  const handleImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      setImageFile(e.target.files[0])
    }
  }

  const handleUpload = async () => {
    if (imageFile) {
      try {
        await uploadImageMutation.mutateAsync({ noteId: note.id, file: imageFile })
        setImageFile(null)
      } catch (err) {
        console.error('Failed to upload image:', err)
      }
    }
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center justify-between">
          <span className="truncate">{note.title}</span>
          <Button
            variant={'default'}
            size="icon"
            onClick={handleDelete}
            disabled={isLoading}
          >
            <Trash2 className="h-4 w-4" />
          </Button>
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <p className="text-sm text-muted-foreground whitespace-pre-wrap">
          {note.content}
        </p>
        {note.image_path && (
          <div className="relative aspect-video w-full overflow-hidden rounded-lg border">
            <Image
              src={`http://localhost:8080/${note.image_path}`}
              alt={note.title}
              fill
              className="object-cover"
            />
          </div>
        )}
      </CardContent>
      <CardFooter className="flex flex-col space-y-2">
        <div className="flex w-full items-center space-x-2">
          <input
            type="file"
            accept="image/*"
            onChange={handleImageChange}
            disabled={isLoading}
            className="hidden"
            id={`file-${note.id}`}
          />
          <label
            htmlFor={`file-${note.id}`}
            className="flex-1 cursor-pointer"
          >
            <Button
              type="button"
              variant={"default"}
              className="w-full"
              disabled={isLoading}
              asChild
            >
              <span>
                <ImageIcon className="mr-2 h-4 w-4" />
                {imageFile ? imageFile.name : 'Choose Image'}
              </span>
            </Button>
          </label>
          {imageFile && (
            <Button
              onClick={handleUpload}
              disabled={isLoading}
              size="icon"
            >
              <Upload className="h-4 w-4" />
            </Button>
          )}
        </div>
        <p className="text-xs text-muted-foreground">
          {new Date(note.created_at).toLocaleDateString()}
        </p>
      </CardFooter>
    </Card>
  )
}
