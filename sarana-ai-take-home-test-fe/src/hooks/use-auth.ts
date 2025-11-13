import { useMutation } from '@tanstack/react-query'
import { apiClient } from '@/lib/api'
import { useAuth } from '@/contexts/AuthContext'

export interface LoginInput {
  email: string
  password: string
}

export interface RegisterInput {
  email: string
  password: string
}

// Login mutation
export function useLogin() {
  const { login } = useAuth()

  return useMutation({
    mutationFn: (data: LoginInput) => apiClient.login(data),
    onSuccess: (response) => {
      login(response.token, response.user)
    },
  })
}

// Register mutation
export function useRegister() {
  const { login } = useAuth()

  return useMutation({
    mutationFn: (data: RegisterInput) => apiClient.register(data),
    onSuccess: (response) => {
      login(response.token, response.user)
    },
  })
}
