import { setContext, getContext } from 'svelte';
import { writable, type Writable, type Readable, derived } from 'svelte/store';
import type { User } from '$lib/types/auth.types';

const AUTH_SESSION_KEY = Symbol('auth-session');
const USER_PROFILE_KEY = Symbol('user-profile');
const AUTH_ACTIONS_KEY = Symbol('auth-actions');

export interface AuthSessionContext {
  token: Writable<string | null>;
  isAuthenticated: Readable<boolean>;
}

export interface UserProfileContext {
  user: Writable<User | null>;
  displayName: Readable<string>;
  email: Readable<string>;
  roleName: Readable<string>;
}

export interface AuthActionsContext {
  clearAll: () => void;
  setUser: (user: User) => void;
  setToken: (token: string) => void;
}

export function initAuthContext(initialUser?: User | null, initialToken?: string | null) {
  const token = writable<string | null>(initialToken ?? null);
  const user = writable<User | null>(initialUser ?? null);

  const isAuthenticated = derived(token, $token => $token !== null);

  const displayName = derived(user, $user => {
    if (!$user) return '';
    const parts = [
      $user.user_profile?.first_name,
      $user.user_profile?.last_name
    ].filter(Boolean);
    return parts.join(' ') || $user.username || '';
  });

  const email = derived(user, $user => $user?.user_profile?.email ?? '');
  const roleName = derived(user, $user => $user?.role?.name ?? '');

  const sessionCtx: AuthSessionContext = { token, isAuthenticated };
  const profileCtx: UserProfileContext = { user, displayName, email, roleName };
  const actionsCtx: AuthActionsContext = {
    clearAll: () => {
      token.set(null);
      user.set(null);
    },
    setUser: (u: User) => user.set(u),
    setToken: (t: string) => token.set(t)
  };

  setContext(AUTH_SESSION_KEY, sessionCtx);
  setContext(USER_PROFILE_KEY, profileCtx);
  setContext(AUTH_ACTIONS_KEY, actionsCtx);

  return {
    session: sessionCtx,
    profile: profileCtx,
    actions: actionsCtx,

    cleanup: () => {
      token.set(null);
      user.set(null);
    }
  };
}

export function getAuthSession(): AuthSessionContext {
  return getContext<AuthSessionContext>(AUTH_SESSION_KEY);
}

export function getUserProfile(): UserProfileContext {
  return getContext<UserProfileContext>(USER_PROFILE_KEY);
}

export function getAuthActions(): AuthActionsContext {
  return getContext<AuthActionsContext>(AUTH_ACTIONS_KEY);
}
