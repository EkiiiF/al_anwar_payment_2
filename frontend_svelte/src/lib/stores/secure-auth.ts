import { writable, derived, type Readable } from 'svelte/store';
import type { User } from '$lib/types/auth.types';

export interface SecureAuthState {
  user: User | null;
  accessToken: string | null;
  isAuthenticated: boolean;
  loading: boolean;
}

const INITIAL_STATE: SecureAuthState = {
  user: null,
  accessToken: null,
  isAuthenticated: false,
  loading: true
};

function createSecureAuthStore() {
  const { subscribe, set, update } = writable<SecureAuthState>({ ...INITIAL_STATE });
  function antiPersistenceAudit(): void {
    if (typeof window === 'undefined') return;

    const suspiciousKeys = ['access_token', 'refresh_token', 'token', 'auth', 'session', 'jwt'];
    for (const key of suspiciousKeys) {
      if (localStorage.getItem(key) !== null) {
        console.warn(
          `[SECURITY] ⚠️ Terdeteksi key "${key}" di localStorage! ` +
          `Data sensitif TIDAK BOLEH tersimpan di localStorage. ` +
          `Menghapus key tersebut...`
        );
        localStorage.removeItem(key);
      }
      if (sessionStorage.getItem(key) !== null) {
        console.warn(
          `[SECURITY] ⚠️ Terdeteksi key "${key}" di sessionStorage! ` +
          `Data sensitif TIDAK BOLEH tersimpan di sessionStorage. ` +
          `Menghapus key tersebut...`
        );
        sessionStorage.removeItem(key);
      }
    }
  }

  return {
    subscribe,
    login(user: User, accessToken: string): void {
      set({
        user,
        accessToken,
        isAuthenticated: true,
        loading: false
      });

      if (import.meta.env.DEV) {
        antiPersistenceAudit();
      }
    },
    logout(): void {
      set({ ...INITIAL_STATE, loading: false });
      if (typeof window !== 'undefined') {
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
        sessionStorage.removeItem('access_token');
        sessionStorage.removeItem('refresh_token');
      }
    },
    setUser(user: User): void {
      update(state => ({
        ...state,
        user,
        isAuthenticated: true,
        loading: false
      }));
    },
    setLoading(loading: boolean): void {
      update(state => ({ ...state, loading }));
    },
    reset(): void {
      set({ ...INITIAL_STATE, loading: false });
    },
    getToken(): string | null {
      let token: string | null = null;
      const unsubscribe = subscribe(state => {
        token = state.accessToken;
      });
      unsubscribe();
      return token;
    },
    audit: antiPersistenceAudit
  };
}

export const secureAuth = createSecureAuthStore();

export const isAuthenticated: Readable<boolean> = derived(
  secureAuth,
  $auth => $auth.isAuthenticated
);

export const currentUser: Readable<User | null> = derived(
  secureAuth,
  $auth => $auth.user
);

export const isLoading: Readable<boolean> = derived(
  secureAuth,
  $auth => $auth.loading
);
