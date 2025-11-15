import { LoginRequest, RegisterRequest, User } from "@/types";
import { useMutation, useQuery } from "@tanstack/react-query";
import { BASE_API } from "../api.config";
import { getAuthCookie, setAuthCookie } from "@/lib/cookie";

interface AuthResponse {
  success: boolean;
  message: string;
  data: {
    token: string;
    user: User;
  };
}

export const authKeys = {
  user: ['auth', 'user'] as const,
  token: ['auth', 'token'] as const,
};

export const useLogin = () => {
  return useMutation({
    mutationFn: async (data: LoginRequest) => {
      const res = await fetch(`${BASE_API}/login`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
      });

      if (!res.ok) {
        const error = await res.json();
        throw new Error(error.message || "Failed to login");
      }

      const response = await res.json() as AuthResponse;
      return response;
    },
    onSuccess: (data) => {
      setAuthCookie(data.data.token);
    }
  });
};

export const useRegister = () => {

  return useMutation({
    mutationFn: async (data: RegisterRequest) => {
      const res = await fetch(`${BASE_API}/register`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
      });

      if (!res.ok) {
        const error = await res.json();
        throw new Error(error.message || "Failed to register");
      }

      const response = await res.json() as AuthResponse;

      return response;
    },
    onSuccess: (data) => {
      setAuthCookie(data.data.token);
    }
  });
};

export const useLogout = () => {
  return useMutation({
    mutationFn: async () => {
      document.cookie = 'sarana-note-app-auth=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT; SameSite=Lax';
      document.cookie = 'user_email=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT; SameSite=Lax';
    },
  });
};


export const useAuth = () => {
  return useQuery({
    queryKey: authKeys.user,
    queryFn: async () => {
      const token = getAuthCookie();

      if (!token) throw new Error('No token found');

      const res = await fetch(`${BASE_API}/notes`, {
        method: "GET",
        headers: {
          "Authorization": `Bearer ${token}`,
        },
      });

      if (!res.ok) {
        document.cookie = 'sarana-note-app-auth=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT; SameSite=Lax';
        document.cookie = 'user_email=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT; SameSite=Lax';
        throw new Error('Invalid token');
      }
    },
    retry: false,
    staleTime: Infinity,
  });
};