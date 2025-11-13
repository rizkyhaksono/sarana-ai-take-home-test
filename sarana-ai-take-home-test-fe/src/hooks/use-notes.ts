import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { apiClient } from '@/lib/api'

export interface Note {
  id: string
  user_id: string
  title: string
  content: string
  image_path?: string
  created_at: string
  updated_at: string
}

export interface CreateNoteInput {
  title: string
  content: string
}

// Query keys
export const noteKeys = {
  all: ['notes'] as const,
  lists: () => [...noteKeys.all, 'list'] as const,
  list: () => [...noteKeys.lists()] as const,
  details: () => [...noteKeys.all, 'detail'] as const,
  detail: (id: string) => [...noteKeys.details(), id] as const,
}

// Get all notes
export function useNotes() {
  return useQuery({
    queryKey: noteKeys.list(),
    queryFn: () => apiClient.getNotes(),
  })
}

// Get single note
export function useNote(id: string) {
  return useQuery({
    queryKey: noteKeys.detail(id),
    queryFn: () => apiClient.getNote(id),
    enabled: !!id,
  })
}

// Create note
export function useCreateNote() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (data: CreateNoteInput) => apiClient.createNote(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: noteKeys.lists() })
    },
  })
}

// Delete note
export function useDeleteNote() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: string) => apiClient.deleteNote(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: noteKeys.lists() })
    },
  })
}

// Upload image
export function useUploadNoteImage() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, file }: { id: string; file: File }) =>
      apiClient.uploadNoteImage(id, file),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: noteKeys.detail(variables.id) })
      queryClient.invalidateQueries({ queryKey: noteKeys.lists() })
    },
  })
}
