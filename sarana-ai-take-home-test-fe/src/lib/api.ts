const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

export interface User {
  id: string
  email: string
  created_at: string
}

export interface Note {
  id: string
  user_id: string
  title: string
  content: string
  image_path?: string
  created_at: string
  updated_at: string
}

export interface AuthResponse {
  token: string
  user: User
}

export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  email: string
  password: string
}

export interface CreateNoteRequest {
  title: string
  content: string
}

class ApiClient {
  private baseURL: string

  constructor(baseURL: string) {
    this.baseURL = baseURL
  }

  private getAuthHeader(): HeadersInit {
    const token = localStorage.getItem('token')
    return token ? { Authorization: `Bearer ${token}` } : {}
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseURL}${endpoint}`
    const headers = {
      'Content-Type': 'application/json',
      ...this.getAuthHeader(),
      ...options.headers,
    }

    const response = await fetch(url, {
      ...options,
      headers,
    })

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'An error occurred')
    }

    return response.json()
  }

  // Auth endpoints
  async register(data: RegisterRequest): Promise<AuthResponse> {
    return this.request<AuthResponse>('/register', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async login(data: LoginRequest): Promise<AuthResponse> {
    return this.request<AuthResponse>('/login', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  // Notes endpoints
  async getNotes(): Promise<Note[]> {
    return this.request<Note[]>('/notes', {
      method: 'GET',
    })
  }

  async getNote(id: string): Promise<Note> {
    return this.request<Note>(`/notes/${id}`, {
      method: 'GET',
    })
  }

  async createNote(data: CreateNoteRequest): Promise<Note> {
    return this.request<Note>('/notes', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async deleteNote(id: string): Promise<{ message: string }> {
    return this.request<{ message: string }>(`/notes/${id}`, {
      method: 'DELETE',
    })
  }

  async uploadNoteImage(
    id: string,
    file: File
  ): Promise<{ message: string; image_path: string }> {
    const formData = new FormData()
    formData.append('image', file)

    const token = localStorage.getItem('token')
    const headers: HeadersInit = token
      ? { Authorization: `Bearer ${token}` }
      : {}

    const response = await fetch(`${this.baseURL}/notes/${id}/image`, {
      method: 'POST',
      headers,
      body: formData,
    })

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Failed to upload image')
    }

    return response.json()
  }
}

export const apiClient = new ApiClient(API_URL)
