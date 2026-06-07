import { secureAuth } from './secure-auth';
import { tokenStore } from '$lib/api/client';
import type { User } from '$lib/types/auth.types';
import { derived } from 'svelte/store';

function createAuthStore() {
  const { subscribe } = derived(secureAuth, $state => ({
    user: $state.user,
    isAuthenticated: $state.isAuthenticated,
    loading: $state.loading
  }));

  return {
    subscribe,

    login: (user: User, token: string) => {
      tokenStore.set(token);
      secureAuth.login(user, token);
    },

    logout: () => {
      tokenStore.clear();
      secureAuth.logout();
      if (typeof window !== 'undefined') {
        fetch('/api/auth/session', { method: 'DELETE' }).catch(() => {});
      }
    },

    setUser: (user: User) => {
      secureAuth.setUser(user);
    },

    setLoading: (loading: boolean) => {
      secureAuth.setLoading(loading);
    }
  };
}

export const auth = createAuthStore();
