import { PUBLIC_API_URL } from '$env/static/public';
import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import type { ApiResponse } from '$lib/types/api.types';
import { secureAuth } from '$lib/stores/secure-auth';

const API_BASE_URL = ''; 

let _memoryToken: string | null = null;

export const tokenStore = {

  get: (): string | null => _memoryToken,

  set: (token: string): void => {
    _memoryToken = token;
  },

  clear: (): void => {
    _memoryToken = null;

    if (typeof localStorage !== 'undefined') {
      localStorage.removeItem('access_token');
      localStorage.removeItem('refresh_token');
    }
    if (typeof sessionStorage !== 'undefined') {
      sessionStorage.removeItem('access_token');
      sessionStorage.removeItem('refresh_token');
    }
  }
};

// ─── HTTP Request Helper ────────────────────────────────────────
async function request<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<ApiResponse<T>> {
  const token = tokenStore.get();

  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    'Cache-Control': 'no-cache, no-store, must-revalidate',
    'Pragma': 'no-cache',
    'Expires': '0',
    ...(options.headers as Record<string, string>)
  };

  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  let url = `${API_BASE_URL}${endpoint}`;
  if (options.method === 'GET') {
    const separator = url.includes('?') ? '&' : '?';
    url += `${separator}_t=${Date.now()}`;
  }

  const res = await fetch(url, {
    ...options,
    headers
  });

  if (res.status === 401) {
    // Token expired/invalid — bersihkan dari memori
    tokenStore.clear();
    secureAuth.logout();
    if (typeof window !== 'undefined') {
      fetch('/api/auth/session', { method: 'DELETE' }).catch(() => {});
    }

    if (typeof window !== 'undefined' && window.location.pathname !== '/login') {
      goto('/login');
    }
    throw new Error('Sesi telah berakhir. Silakan login kembali.');
  }

  let json: ApiResponse<T>;
  try {
    json = await res.json();
  } catch {
    throw new Error(`Server error (${res.status}): Respons tidak valid.`);
  }

  if (!res.ok) {
    throw new Error(json.data as string || json.message || `Terjadi kesalahan (${res.status})`);
  }

  return json;
}

// ─── Exported API Client ────────────────────────────────────────
export const apiClient = {
  get: <T>(endpoint: string) => request<T>(endpoint, { method: 'GET' }),
  post: <T>(endpoint: string, body?: unknown) =>
    request<T>(endpoint, {
      method: 'POST',
      body: body ? JSON.stringify(body) : undefined
    }),
  put: <T>(endpoint: string, body?: unknown) =>
    request<T>(endpoint, {
      method: 'PUT',
      body: body ? JSON.stringify(body) : undefined
    }),
  patch: <T>(endpoint: string, body?: unknown) =>
    request<T>(endpoint, {
      method: 'PATCH',
      body: body ? JSON.stringify(body) : undefined
    }),
  delete: <T>(endpoint: string) => request<T>(endpoint, { method: 'DELETE' }),

  download: async (endpoint: string, defaultFileName: string): Promise<void> => {
    const token = tokenStore.get();
    const headers: Record<string, string> = {};
    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    const res = await fetch(`${API_BASE_URL}${endpoint}`, {
      method: 'GET',
      headers
    });

    if (res.status === 401) {
      tokenStore.clear();
      if (typeof window !== 'undefined') {
        window.location.href = '/login';
      }
      throw new Error('Sesi telah berakhir. Silakan login kembali.');
    }

    if (!res.ok) {
      throw new Error(`Gagal mengunduh berkas (${res.status})`);
    }

    const disposition = res.headers.get('Content-Disposition');
    let fileName = defaultFileName;
    if (disposition && disposition.includes('attachment')) {
      const filenameRegex = /filename[^;=\n]*=((['"]).*?\2|[^;\n]*)/;
      const matches = filenameRegex.exec(disposition);
      if (matches != null && matches[1]) {
        fileName = matches[1].replace(/['"]/g, '');
      }
    }

    const blob = await res.blob();
    if (typeof window !== 'undefined') {
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = fileName;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      window.URL.revokeObjectURL(url);
    }
  }
};
