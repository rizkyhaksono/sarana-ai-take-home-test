'use client'

import { useState } from 'react'
import { Trash2, Upload, Image as ImageIcon } from 'lucide-react'
import { Note } from '@/hooks/use-notes'
import { useDeleteNote, useUploadNoteImage } from '@/hooks/use-notes'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'

interface NoteCardProps {
  note: Note
}

export function NoteCard({ note }: NoteCardProps) {
  const [imageFile, setImageFile] = useState<File | null>(null)
  const { mutate: deleteNote, isPending: isDeleting } = useDeleteNote()
  const { mutate: uploadImage, isPending: isUploading } = useUploadNoteImage()

  const handleDelete = () => {
    if (confirm('Are you sure you want to delete this note?')) {
      deleteNote(note.id)
    }
  }

  const handleImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      setImageFile(e.target.files[0])
    }
  }

  const handleUpload = () => {
    if (imageFile) {
      uploadImage({ id: note.id, file: imageFile }, {
        onSuccess: () => {
          setImageFile(null)
        },
      })
    }
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center justify-between">
          <span className="truncate">{note.title}</span>
          <Button
            variant="ghost"
            size="icon"
            onClick={handleDelete}
            disabled={isDeleting}
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
            <img
              src={`http://localhost:8080/${note.image_path}`}
              alt={note.title}
              className="object-cover w-full h-full"
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
            disabled={isUploading}
            className="hidden"
            id={`file-${note.id}`}
          />
          <label
            htmlFor={`file-${note.id}`}
            className="flex-1 cursor-pointer"
          >
            <Button
              type="button"
              variant="outline"
              className="w-full"
              disabled={isUploading}
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
              disabled={isUploading}
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
