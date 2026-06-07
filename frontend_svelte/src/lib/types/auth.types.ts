export type UserRole = 'super_user' | 'guardian' | 'pengasuh';

export interface Role {
  id: string;
  name: UserRole;
  description: string;
}

export interface User {
  id: string;
  username: string;
  role_id: string;
  role: Role;
  first_name: string;
  middle_name?: string;
  last_name?: string;
  email: string;
  phone_number: string;
  gender: 'L' | 'P';
  birth_date?: string;
  address?: string;
  is_active: boolean;
  last_login_at?: string;
  created_at: string;
}

export interface AuthSession {
  access_token: string;
  token_type: string;
  user_id: string;
  username: string;
  role: string;
}
