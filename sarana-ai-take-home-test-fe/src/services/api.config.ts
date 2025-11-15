import { getAuthCookie } from "@/lib/cookie";

export const BASE_API = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export const baseHeaders = async () => {
  const token = getAuthCookie()?.token || null;
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  };
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }
  return headers;
};
