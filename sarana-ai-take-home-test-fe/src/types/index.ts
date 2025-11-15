// Base Types
export interface User {
  id: string;
  email: string;
  created_at: string;
}

export interface Note {
  id: string;
  user_id: string;
  title: string;
  content: string;
  image_path?: string;
  created_at: string;
  updated_at: string;
}

export interface Log {
  id: number;
  datetime: string;
  method: string;
  endpoint: string;
  status_code: number;
  headers: string;
  request_body: string;
  response_body: string;
  created_at: string;
}

// API Request Types
export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
}

export interface CreateNoteRequest {
  title: string;
  content: string;
  image?: File;
}

// API Response Types
export interface BaseResponse<T = unknown> {
  success: boolean;
  message: string;
  data?: T;
  error?: ErrorData;
}

export interface ErrorData {
  code: string;
  message: string;
  details?: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}

export interface PaginationParams {
  search?: string;
  sort_by?: string;
  order?: 'ASC' | 'DESC';
  page?: number;
  limit?: number;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}

export interface NotesData {
  notes: Note[];
  page: number;
  total: number;
  per_page: number;
  total_pages: number;
  total_items: number;
}

export interface NotesResponse {
  success: boolean;
  message: string;
  data: NotesData;
}

export interface LogsData {
  logs: Log[];
  page: number;
  per_page: number;
  total_pages: number;
  total: number;
}

export interface LogsResponse {
  success: boolean;
  message: string;
  data: LogsData;
}

export interface UploadImageResponse {
  message: string;
  image_path: string;
}
