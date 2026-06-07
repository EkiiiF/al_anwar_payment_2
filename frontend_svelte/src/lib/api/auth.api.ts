import { apiClient } from './client';
import type { User, AuthSession } from '$lib/types/auth.types';

export const authApi = {
  login: (username: string, password: string) =>
    apiClient.post<AuthSession>('/api/v1/auth/login', { username, password }),
  
  logout: () => apiClient.post('/api/v1/auth/logout'),
  
  getProfile: () => apiClient.get<User>('/api/v1/auth/profile'),
  
  updateProfile: (data: unknown) => apiClient.put<User>('/api/v1/auth/profile', data),
  
  changePassword: (data: unknown) => apiClient.put('/api/v1/auth/change-password', data)
};
