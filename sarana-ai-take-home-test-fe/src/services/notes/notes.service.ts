import { CreateNoteRequest, Note, NotesResponse, PaginationParams, UploadImageResponse } from "@/types";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { BASE_API, baseHeaders } from "../api.config";
import { getAuthCookie } from "@/lib/cookie";

export const notesKeys = {
  all: ['notes'] as const,
  lists: () => [...notesKeys.all, 'list'] as const,
  list: (params?: PaginationParams) => [...notesKeys.lists(), params] as const,
  details: () => [...notesKeys.all, 'detail'] as const,
  detail: (id: string) => [...notesKeys.details(), id] as const,
};

export const useGetNotes = (params?: PaginationParams) => {
  return useQuery({
    queryKey: notesKeys.list(params),
    queryFn: async () => {
      const queryParams = new URLSearchParams();

      if (params?.search) queryParams.append('search', params.search);
      if (params?.sort_by) queryParams.append('sort_by', params.sort_by);
      if (params?.order) queryParams.append('order', params.order);
      if (params?.page) queryParams.append('page', params.page.toString());
      if (params?.limit) queryParams.append('limit', params.limit.toString());

      const query = queryParams.toString();
      const endpoint = query ? `/notes?${query}` : '/notes';

      const res = await fetch(`${BASE_API}${endpoint}`, {
        method: "GET",
        headers: await baseHeaders(),
      });

      if (!res.ok) {
        const error = await res.json();
        throw new Error(error.message || "Failed to fetch notes");
      }

      const response = await res.json();
      return response as NotesResponse;
    },
  });
};

export const useGetNoteById = (id: string, enabled = true) => {
  return useQuery({
    queryKey: notesKeys.detail(id),
    queryFn: async () => {
      const res = await fetch(`${BASE_API}/notes/${id}`, {
        method: "GET",
        headers: await baseHeaders(),
      });

      if (!res.ok) {
        const error = await res.json();
        throw new Error(error.message || "Failed to fetch note");
      }

      const data = await res.json();
      return data.data as Note;
    },
    enabled: !!id && enabled,
    staleTime: 30000,
  });
};

export const useCreateNote = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (data: CreateNoteRequest) => {
      const formData = new FormData();
      formData.append('title', data.title);
      formData.append('content', data.content);
      if (data.image) {
        formData.append('image', data.image);
      }

      const res = await fetch(`${BASE_API}/notes`, {
        method: "POST",
        headers: {
          'Authorization': `Bearer ${getAuthCookie()?.token}`,
        },
        body: formData,
      });

      if (!res.ok) {
        const error = await res.json();
        throw new Error(error.message || "Failed to create note");
      }

      const response = await res.json();
      return response.data as Note;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: notesKeys.lists() });
    },
  });
};

export const useUpdateNote = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, data }: { id: string; data: CreateNoteRequest }) => {
      const formData = new FormData();
      formData.append('title', data.title);
      formData.append('content', data.content);
      if (data.image) {
        formData.append('image', data.image);
      }

      const res = await fetch(`${BASE_API}/notes/${id}`, {
        method: "PUT",
        headers: {
          'Authorization': `Bearer ${getAuthCookie()?.token}`,
        },
        body: formData,
      });

      if (!res.ok) {
        const error = await res.json();
        throw new Error(error.message || "Failed to update note");
      }

      const response = await res.json();
      return response.data as Note;
    },
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: notesKeys.lists() });
      queryClient.invalidateQueries({ queryKey: notesKeys.detail(variables.id) });
    },
  });
};

export const useDeleteNote = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (id: string) => {
      const res = await fetch(`${BASE_API}/notes/${id}`, {
        method: "DELETE",
        headers: await baseHeaders(),
      });

      if (!res.ok) {
        const error = await res.json();
        throw new Error(error.message || "Failed to delete note");
      }

      return await res.json();
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: notesKeys.lists() });
    },
  });
};

export const useUploadNoteImage = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ noteId, file }: { noteId: string; file: File }) => {
      const token = getAuthCookie();
      const formData = new FormData();
      formData.append('image', file);

      const res = await fetch(`${BASE_API}/notes/${noteId}/image`, {
        method: "POST",
        headers: {
          'Authorization': `Bearer ${token}`,
        },
        body: formData,
      });

      if (!res.ok) {
        const error = await res.json();
        throw new Error(error.message || "Failed to upload image");
      }

      return await res.json() as UploadImageResponse;
    },
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: notesKeys.lists() });
      queryClient.invalidateQueries({ queryKey: notesKeys.detail(variables.noteId) });
    },
  });
};

export const useGetNoteImage = (noteId: string, enabled = true) => {
  return useQuery({
    queryKey: ['note-image', noteId],
    queryFn: async () => {
      const res = await fetch(`${BASE_API}/notes/${noteId}/image`, {
        method: "GET",
        headers: await baseHeaders(),
      });
      if (!res.ok) {
        const error = await res.json();
        throw new Error(error.message || "Failed to fetch note image");
      }
      const blob = await res.blob();
      return URL.createObjectURL(blob);
    },
    enabled: !!noteId && enabled,
    staleTime: 30000,
  });
}